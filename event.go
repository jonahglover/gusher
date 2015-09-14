// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gusher

import (
	"fmt"
)

const (
	PING_EVENT                          = "pusher:ping"
	PONG_EVENT                          = "pusher:pong"
	PING_EVENT_PAYLOAD                  = `{"event":"pusher:ping","data":"{}"}`
	PONG_EVENT_PAYLOAD                  = `{"event":"pusher:pong","data":"{}"}`
	PUSHER_CONNECTION_ESTABLISHED_EVENT = "pusher:connection_established"
	PUSHER_SUBSCRIBE_EVENT              = "pusher:subscribe"
	PUSHER_UNSUBSCRIBE_EVENT            = "pusher:unsubscribe"
	PUSHER_ERROR_EVENT                  = "pusher:error"
	INTERNAL_UNSUBSCRIBE_EVENT          = "internal:unsubscribe"
)

type Event struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

func (e *Event) Encode() string {
	return fmt.Sprintf(`{"event":"%s","data":%s}`, e.Event, e.Data)
}
