// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// TODO PusherChannel interface
// public and private channel

type Subscription struct {
	Name     string
	channel  chan *Event
	bindings map[string]chan *Event
}

// TODO return an error or error code or something
// callback should take an Event and do some processing with it
// this is why it would be good for Event to be an interface

func (s *Subscription) Bind(eventName string) chan *Event {
	s.bindings[eventName] = make(chan *Event)
	return s.bindings[eventName]
}

func (s *Subscription) listen() {
	for {
		// XXX
		//maybe this should be a buffered channel?
		//should I be running the call backs in a goroutine? probably not
		e := <-s.channel
		if s.bindings[e.Event] != nil {
			s.bindings[e.Event] <- e
		}
	}
}

// should already be connected at this point
func NewSubscription(name string) *Subscription {
	s := &Subscription{Name: name, channel: make(chan *Event), bindings: make(map[string]chan *Event)}
	go s.listen()
	return s
}
