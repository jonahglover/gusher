// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

func main() {
	socket, err := NewPusher("de504dc5763aeef9ff52", Scheme("WSS"), ProtocolVersion("7"))
	subscription := socket.Subscribe("order_book")
	eventChannel := subscription.Bind("data")
	if err != nil {
		fmt.Println("issues")
	}
	for {
		ev := <-eventChannel
		fmt.Println(ev.Data)
	}

}
