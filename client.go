// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

const (
	HEARTBEAT_INTERVAL = 10 // time to wait in between heartbeats
	// [scheme]://ws.pusherapp.com:[port]/app/[key]
	WS_SCHEME         = "ws"
	WSS_SCHEME        = "wss"
	PUSHER_HOST       = "ws.pusherapp.com"
	APP               = "APP"
	WS_PORT           = "80"
	WSS_PORT          = "443"
	MAX_MESSAGE_BYTES = 10000 // max message size for pusher is 10kB http://www.quora.com/What-is-the-maximum-message-size-for-Pusher-1
	PROTOCOL_VERSION  = "7"
	ORIGIN            = "http://localhost/"
)

// concepts of different types of clients
type Pusher struct {
	ws            *websocket.Conn
	subscriptions map[string]*Subscription
}

func (c *Pusher) Heartbeat() {
	for {
		websocket.Message.Send(c.ws, PUSHER_PING_EVENT_PAYLOAD)
		time.Sleep(HEARTBEAT_INTERVAL * time.Second)
	}
}

func (c *Pusher) Listen() {
	for {
		var event Event
		err := websocket.JSON.Receive(c.ws, &event)
		if err != nil {
		} else {
			switch event.Event {
			case PUSHER_PING_EVENT:
				websocket.Message.Send(c.ws, PUSHER_PONG_EVENT_PAYLOAD)
			case PUSHER_PONG_EVENT:
			default:
				c.subscriptions[event.Channel].channel <- &event
			}
		}
	}
}

// TODO ERROR
func (c *Pusher) Subscribe(name string) *Subscription {
	err := websocket.Message.Send(c.ws, fmt.Sprintf(`{"event":"pusher:subscribe","data":{"channel":"%s"}}`, name))
	if err != nil {
		return nil
	}

	// looks like there was no error. lets add this subscription :D
	c.subscriptions[name] = NewSubscription(name)
	return c.subscriptions[name]
}

// TODO give this options pertaining to the type of client
func NewPusher(key string) (*Pusher, error) {
	url := WSS_SCHEME + "://" + PUSHER_HOST + "/app/" + key + "?protocol=" + PROTOCOL_VERSION

	// Attempt to Establish Connection wuth Pusher
	ws, err := websocket.Dial(url, "", ORIGIN)
	if err != nil {
		return nil, err
	}

	// Read Response from Pusher
	res := make([]byte, MAX_MESSAGE_BYTES)

	msgLength, err := ws.Read(res)
	if err != nil {
		return nil, err
	}

	var event Event
	err = json.Unmarshal(res[0:msgLength], &event)
	if err != nil {
		return nil, err
	}

	switch event.Event {
	case PUSHER_ERROR:
		// TODO
		fmt.Println("pusher error")
	case PUSHER_CONNECTION_ESTABLISHED:
		c := Pusher{ws: ws, subscriptions: make(map[string]*Subscription)}
		go c.Heartbeat()
		go c.Listen()
		return &c, nil
	}
	return nil, nil
}
