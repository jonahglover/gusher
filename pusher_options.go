// Copyright 2015 Jonah Glover. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gusher

const (
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

func Scheme(scheme string) func(*Pusher) error {
	log.Info("Setting Scheme: " + scheme)
	return func(p *Pusher) error {
		return p.setScheme(scheme)
	}
}

func (p *Pusher) setScheme(scheme string) error {
	p.scheme = scheme
	return nil
}

func PusherHost(host string) func(*Pusher) error {
	log.Info("Setting Pusher Host: " + host)
	return func(p *Pusher) error {
		return p.setHost(host)
	}
}

func (p *Pusher) setHost(host string) error {
	p.host = host
	return nil
}

func ProtocolVersion(protocolVersion string) func(*Pusher) error {
	log.Info("Setting Protocol Version: " + protocolVersion)
	return func(p *Pusher) error {
		return p.setProtocolVersion(protocolVersion)
	}
}

func (p *Pusher) setProtocolVersion(protocolVersion string) error {
	p.protocolVersion = protocolVersion
	return nil
}
