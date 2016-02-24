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
	connected chan bool
	done      chan error
	error     chan error

	handlers  map[string]Handler
	upstreams map[string]Upstream

	channels  channels
	exchanges exchanges
	queues    queues

	online bool
}

type channels struct {
	direct *amqp.Channel
	topic  *amqp.Channel
}
type exchanges struct {
	direct string
	topic  string
}
type queues struct {
	direct string
	topic  string
}

type Sender struct {
	name string
	uuid  string
}

type Destination struct {
	name    string
	uuid    string
	handler string
}

type Receiver struct {
	name    string
	uuid    string
	handler string
}

type Handler  func (Sender, []byte) error

type Upstream func (Sender, Destination, []byte) error