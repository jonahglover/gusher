// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"time"
)

type Pusher struct {
	ws              *websocket.Conn
	subscriptions   map[string]*Subscription
	ActivityTimeout float64
	SocketID        string
	scheme          string
	host            string
	appKey          string
	protocolVersion string
	pusherHost      string
	readBuffer      chan *Event
}

// TODO fix me to use write pump
func (c *Pusher) heartbeat() {
	log.Info("Starting Heartbeat")
	log.Info("Sending Ping every " + fmt.Sprintf("%.6f", c.ActivityTimeout) + " seconds.")
	for {
		websocket.Message.Send(c.ws, PING_EVENT_PAYLOAD)
		time.Sleep(time.Duration(c.ActivityTimeout) * time.Second)
	}
}

// TODO Build a 'read pump' intermediate channel
// that can do this processing rather than having it all in this loop.
// much cleaner

func (c *Pusher) listen() {
	log.Info("Listening for events from Pusher")
	for {
		var event Event
		err := websocket.JSON.Receive(c.ws, &event)
		if err != nil {
		} else {
			switch event.Event {
			case PING_EVENT:
				websocket.Message.Send(c.ws, PONG_EVENT_PAYLOAD)
			case PONG_EVENT:
			default:
				if err != nil {
					log.Error("Could not read Channel Event from websocket")
				}
				c.subscriptions[event.Channel].channel <- &event
			}
		}
	}
}

func (c *Pusher) sendEvent(e *Event) {
	c.readBuffer <- e
}

func (c *Pusher) serve() error {
	log.Info("Listening for events to send to Pusher")
	for {
		event := <-c.readBuffer
		err := websocket.Message.Send(c.ws, event.Encode())
		if err != nil {
			return err
		}
	}
}

func (c *Pusher) Subscribe(name string) *Subscription {
	subscribeEvent := &Event{Event: PUSHER_SUBSCRIBE_EVENT, Data: fmt.Sprintf(`{"channel":"%s"}`, name)}
	c.sendEvent(subscribeEvent)
	// looks like there was no error. lets add this subscription :D
	c.subscriptions[name] = NewSubscription(name)
	return c.subscriptions[name]
}

func (c *Pusher) Unsubscribe(name string) error {
	var subscription *Subscription
	if c.subscriptions[name] != nil {
		subscription = c.subscriptions[name]
	}

	unsubscribeEvent := &Event{Event: PUSHER_UNSUBSCRIBE_EVENT, Data: fmt.Sprintf(`{"channel":"%s"}`, name)}
	c.sendEvent(unsubscribeEvent)

	subscription.Unsubscribe()

	delete(c.subscriptions, name)
	return nil
}

func (p *Pusher) listenAndServe() error {
	go p.heartbeat()
	go p.listen()
	go p.serve()
	return nil
}

func NewPusher(key string, options ...func(*Pusher) error) (*Pusher, error) {

	p := Pusher{scheme: WSS_SCHEME, pusherHost: PUSHER_HOST, appKey: key, protocolVersion: PROTOCOL_VERSION}

	// process options

	for _, option := range options {
		option(&p)
	}

	// set URL

	url := p.scheme + "://" + p.pusherHost + "/app/" + p.appKey + "?protocol=" + p.protocolVersion

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
	case PUSHER_ERROR_EVENT:
		if err != nil {
			return nil, err
		}
		log.Error(event.Data)
	case PUSHER_CONNECTION_ESTABLISHED_EVENT:
		log.Info("Connection Established")
		var event Event
		err = json.Unmarshal(res[0:msgLength], &event)
		if err != nil {
			return nil, err
		}
		var data map[string]interface{}
		json.Unmarshal([]byte(event.Data), &data)
		c := Pusher{ws: ws, subscriptions: make(map[string]*Subscription), ActivityTimeout: data["activity_timeout"].(float64), SocketID: data["socket_id"].(string), readBuffer: make(chan *Event)}
		c.listenAndServe()
		return &c, nil
	}
	return nil, errors.New("Received unknown event from Pusher.")
}
