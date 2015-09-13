// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

const (
	PING_EVENT                          = "pusher:ping"
	PONG_EVENT                          = "pusher:pong"
	PING_EVENT_PAYLOAD                  = `{"event":"pusher:ping","data":"{}"}`
	PONG_EVENT_PAYLOAD                  = `{"event":"pusher:pong","data":"{}"}`
	PUSHER_CONNECTION_ESTABLISHED_EVENT = "pusher:connection_established"
	PUSHER_SUBSCRIBE_EVENT              = "pusher:subscribe"
	PUSHER_ERROR_EVENT                  = "pusher:error"
)

type Event struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}
