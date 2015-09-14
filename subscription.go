// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gusher

import (
	"errors"
	"time"
)

const (
	BIND_TIMEOUT = 10
)

type Subscription struct {
	Name      string
	channel   chan *Event
	bindings  map[string]chan *Event
	readyChan chan bool
	Ready     bool
}

func (s *Subscription) Bind(eventName string) chan *Event {
	if !s.Ready {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(BIND_TIMEOUT * time.Second)
			timeout <- true
		}()
		select {
		case <-s.readyChan:
			log.Info("Client now ready to bind")
		case <-timeout:
			return nil
		}
	}
	log.Info("Binding to event \"" + eventName + "\" on channel \"" + s.Name + "\"")
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
	s.Ready = true
	s.readyChan <- s.Ready
	log.Info("Listening for events on channel \"" + s.Name + "\"")
	for {
		e := <-s.channel
		if e.Event == INTERNAL_UNSUBSCRIBE_EVENT {
			log.Info("Stop listening for events on channel \"" + s.Name + "\"")
			break
		}
		if s.bindings[e.Event] != nil {
			s.bindings[e.Event] <- e
		}
	}
}

func (s *Subscription) Unsubscribe() {
	log.Info("Unsubscribing from \"" + s.Name + "\"")
	s.channel <- &Event{Event: INTERNAL_UNSUBSCRIBE_EVENT}
}

func NewSubscription(name string) *Subscription {
	s := &Subscription{Name: name, channel: make(chan *Event), bindings: make(map[string]chan *Event), Ready: false, readyChan: make(chan bool, 1)}
	go s.listen()
	return s
}
