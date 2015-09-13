// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

const (
	PUSHER_ERROR                  = "pusher:error"
	PUSHER_CONNECTION_ESTABLISHED = "pusher:connection_established"
	PUSHER_PING_EVENT             = "pusher:ping"
	PUSHER_PONG_EVENT             = "pusher:pong"
	PUSHER_PING_EVENT_PAYLOAD     = `{"event":"pusher:ping","data":"{}"}`
	PUSHER_PONG_EVENT_PAYLOAD     = `{"event":"pusher:pong","data":"{}"}`
)

type Event struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}
