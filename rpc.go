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
	r := rpc.RegisterAPP()
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

var (
	rpc RPC
)

// Register application in RPC
func Register (name string, uuid string, token string) (RPC, error) {

	rpc.name  = name
	rpc.uuid  = uuid
	rpc.token = token

	rpc.connect   = make (chan bool)
	rpc.reconnect = make (chan bool)
	rpc.error     = make (chan error)

	rpc.handlers  = make (map[string]Handler)
	rpc.upstreams = make (map[string]Upstream)

	go rpc.listen()
	return rpc, nil
}

// SetURI - set URI connection
func SetURI (uri string) {
	rpc.uri = uri
}

// Start listening for incoming messages
func Listen () {
	rpc.connect <- true
}