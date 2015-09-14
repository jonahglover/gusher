// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
)

type Subscription struct {
	Name     string
	channel  chan *Event
	bindings map[string]chan *Event
}

func (s *Subscription) Bind(eventName string) chan *Event {
	s.bindings[eventName] = make(chan *Event)
	return s.bindings[eventName]
}

func (s *Subscription) Unbind(eventName string) error {
	if s.bindings[eventName] != nil {
		delete(s.bindings, eventName)
		log.Notice("Unbound from event " + eventName + ".")
		return nil
	} else {
		return errors.New("This event binding does not exist")
	}
}

func (s *Subscription) listen() {
	log.Info("Listening for events on channel \"" + s.Name + "\"")
	for {
		e := <-s.channel
		if s.bindings[e.Event] != nil {
			s.bindings[e.Event] <- e
		}
	}
}

func NewSubscription(name string) *Subscription {
	s := &Subscription{Name: name, channel: make(chan *Event), bindings: make(map[string]chan *Event)}
	go s.listen()
	return s
}
