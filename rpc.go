// Copyright 2016 Last.Backend. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package lastbackend/prc implements an amqp queue based messaging system.
Like the standard AMQP application rpc.Router matches incoming messages against
a list of registered routes and upstream and calls a handler for the route or upstream
that matches the route name or upstream or other conditions. The main features are:
	* Messages routing can be matched based received router, upstream.
	* You can send many types of messages, with or without confirmation.


Each app automatically join in a group by its name, you can set uuid of registered app to
provide better delivery.

You should start with your app registering:
	func main() {
		r := rpc.Register()
		r.SetHandler("handler",   SomeHandler)
	}

Setup handler and upstream examples:
	r := rpc.Register()
	r.SetHandler("handler",   SomeHandler)
	r.SetUpstream("upstream", SomeUpstream)

Handlers and Upstreams examples:

	// SomeHandler definition
 	func SomeHandler(s rpc.Sender, message []byte) error {

 	}

 	// SomeUpstream definition
 	func SomeUpstream(s rpc.Sender, r rpc.Recipient, message []byte) error {

 	}
*/
package rpc

// Register application in RPC
func Register(name string, uuid string, token string) (*RPC, error) {

	var rpc RPC

	rpc.name = name
	rpc.uuid = uuid
	rpc.token = token

	rpc.connect = make(chan bool)
	rpc.reconnect = make(chan bool)
	rpc.connected = make(chan bool)

	rpc.limit = 1

	rpc.done = make(chan error)
	rpc.error = make(chan error)

	rpc.handlers = make(map[string]Handler)
	rpc.upstreams = make(map[string]Upstream)
	return &rpc, nil
}

// SetURI - set URI connection
func (r *RPC) SetURI(uri string) {
	r.uri = uri
}

func (r *RPC) SetLimit(limit int) {
	r.limit = limit
}

// Start listening for incoming messages
func (r *RPC) Listen() {
	go r.listen()
	r.online = true
	r.connect <- true
}

func (r *RPC) Connected() chan bool {
	return r.connected
}

func (r *RPC) Cleanup() {
	r.cleanup()
}

func (r *RPC) Shutdown() {
	r.online = false
	r.shutdown()
}
