package rpc

import "github.com/streadway/amqp"

type RPC struct {

	uri  string
	conn *amqp.Connection

	name  string
	uuid  string
	token string

	connect   chan bool
	reconnect chan bool
	error     chan error

	handlers  map[string]Handler
	upstreams map[string]Upstream
}

type Sender struct {
	name string
	uuid  string
}

type Receiver struct {
	name    string
	uuid    string
	handler string
	any   bool
	all   bool
}

type Proxy struct {
	name    string
	uuid    string
	handler string
}

type Handler  func (Sender, []byte) error

type Upstream func (Sender, Receiver, []byte) error