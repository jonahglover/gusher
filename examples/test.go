// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"gusher"
)

func main() {
	socket, err := gusher.NewPusher("de504dc5763aeef9ff52", gusher.Scheme("WSS"), gusher.ProtocolVersion("7"))
	subscription := socket.Subscribe("order_book")
	eventChannel := subscription.Bind("data")
	if err != nil {
		fmt.Println("issues")
	}
	for {
		ev := <-eventChannel
		fmt.Println(ev.Data)
		socket.Unsubscribe("order_book")
		break
	}

	newSubscription := socket.Subscribe("order_book")
	newEventChannel := newSubscription.Bind("data")
	for {
		ev := <-newEventChannel
		fmt.Println(ev.Data)
	}

}
