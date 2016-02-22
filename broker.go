package rpc

import (
	"log"
	"time"
	"os"
	"fmt"
	"github.com/streadway/amqp"
	"strings"
)

func (r *RPC) listen() {
	for {
		select {
		case <-r.connect:
			log.Printf("RPC: %s connect", r.name)
			go r.dial()

		case <-r.reconnect:
			log.Printf("RPC: %s reconnect", r.name)
			timer := time.NewTimer(time.Second)
			<-timer.C
			go r.dial()
			timer.Stop()
		}
	}
}

func (r *RPC) dial() {
	var err error

	if r.uri == "" {
		AMQP_USER := os.Getenv("AMQP_USER")
		AMQP_PASS := os.Getenv("AMQP_PASS")
		AMQP_HOST := os.Getenv("AMQP_HOST")
		AMQP_PORT := os.Getenv("AMQP_PORT")

		r.uri = fmt.Sprintf("amqp://%s:%s@%s:%s/", AMQP_USER, AMQP_PASS, AMQP_HOST, AMQP_PORT)
	}

	r.conn, err = amqp.Dial(r.uri)

	if err != nil {
		log.Println("RPC: Dial error", err)
		r.reconnect <- true
		return
	}

	go func() {
		log.Printf("RPC: Closing: %s", <-r.conn.NotifyClose(make(chan *amqp.Error)))
		r.connect <- true
	}()

	log.Println("RPC: connect consumer and proxier")

}

func (r *RPC) call (s Sender, p Proxy, d Receiver, data []byte) error {
	return r.publish(true, s, p, d, data)
}

func (r *RPC) cast (s Sender, p Proxy, d Receiver, data []byte) error {
	return r.publish(false, s, p, d, data)
}

func (r *RPC) publish (call bool, s Sender, p Proxy, d Receiver, data []byte) error {

	var body []byte
	var hash [256]byte

	var token [32]byte
	copy(token[:], r.token)

	copy(hash[0:32], token[:])
	copy(hash[32:84], s.Sign()[:])
	copy(hash[84:152], p.Sign()[:])
	copy(hash[152:220], d.Sign()[:])


	body = append(body, hash[:]...)
	body = append(body, data[:]...)

	log.Printf("PRC: call %s:%s, send: %dB body (%s)", p.name, p.uuid, len(body), body)

	channel, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	exchange := fmt.Sprintf("%s:%s", p.name, "direct")
	if !call {
		exchange = fmt.Sprintf("%s:%s", p.name, "topic")
	}

	bind     := p.uuid
	if bind == "" {
		bind = p.name
	}

	if err := channel.Publish(exchange, bind, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

func (r *RPC) subscribe() error {
	var err error
	var done = make(chan error)

	var channels = struct {
		direct *amqp.Channel
		topic  *amqp.Channel
	}{}

	var exchanges = struct {
		direct string
		topic  string
	}{
		fmt.Sprintf("%s:%s", r.name, "direct"),
		fmt.Sprintf("%s:%s", r.name, "topic"),
	}


	var queues   = struct {
		direct string
		topic  string
	}{
		fmt.Sprintf("%s:%s", r.name, "direct"),
		fmt.Sprintf("%s:%s", r.name, "topic"),
	}


	// Get hostname for register current instance
	log.Printf("RPC: Create new consumer: %s", r.name)

	channels.direct, err = r.conn.Channel();
	if err != nil {
		log.Println("Channel:", err)
		return err
	}

	err = channels.direct.Qos(100, 0, true)
	if err != nil {
		log.Println("Channel:", err)
		return err
	}

	channels.topic, err = r.conn.Channel();
	if err != nil {
		log.Println("Channel:", err)
		return err
	}

	err = channels.topic.Qos(100, 0, true)
	if err != nil {
		log.Println("Channel:", err)
		return err
	}

	// create direct exchange for guarantee delivery messages
	if err = channels.direct.ExchangeDeclare(exchanges.direct, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// create topic exchange for non guarantee delivery messages
	if err = channels.topic.ExchangeDeclare(exchanges.topic, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// create direct queue for guarantee delivery messages
	if _, err := channels.direct.QueueDeclare(queues.direct, true, false, false, false, nil); err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	// create topic queue for non guarantee delivery messages
	if _, err := channels.topic.QueueDeclare(queues.topic, true, true, false, false, nil); err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	// create bindings for direct messages
	if err = channels.direct.QueueBind(queues.direct, r.uuid, exchanges.direct, false, nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if err = channels.direct.QueueBind(queues.direct, r.name, exchanges.direct, false, nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	// create bindings for topic messages
	if err = channels.direct.QueueBind(queues.topic, r.uuid, exchanges.direct, false, nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if err = channels.direct.QueueBind(queues.topic, r.name, exchanges.direct, false, nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}


	messages, err := channels.direct.Consume(queues.direct, r.name, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	streams, err := channels.topic.Consume(queues.topic, r.name, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	go r.handle (messages, done)
	go r.handle (streams, done)

	return nil
}

func (r *RPC) handle (msgs <-chan amqp.Delivery, done chan error) {
	for d := range msgs {

		go func(m amqp.Delivery) {

			// validate token
			if r.token != string(m.Body[0:32]) {
				log.Println("RPC: token verification failed")
				m.Ack(false)
			}

			// parse sender information
			s := Sender{}
			s.name = string(d.Body[32:48])
			s.uuid = string(d.Body[48:84])

			s.name = s.name[:strings.Index(string(s.name), "\x00")]

			// parse destination information
			p := Proxy{}
			p.name = string(d.Body[84:100])
			p.uuid = string(d.Body[100:136])
			p.handler = string(d.Body[136:152])

			p.name = p.name[:strings.Index(string(p.name), "\x00")]
			p.handler = p.handler[:strings.Index(string(p.handler), "\x00")]

			// parse receiver information
			d := Receiver{}
			d.name = string(m.Body[156:168])
			d.uuid = string(m.Body[168:204])
			d.handler = string(m.Body[204:220])

			d.name = d.name[:strings.Index(string(d.name), "\x00")]
			d.handler = d.handler[:strings.Index(string(d.handler), "\x00")]

			data := m.Body[255:len(m.Body)]

			log.Println("data:", string(data))

			if p.name != r.name {
				log.Println("PRC: send to upstream")
				_, ok := r.upstreams[p.handler]

				if !ok {
					log.Println("RPC: upstream not found", p.handler)
					m.Ack(false)
					return
				}

				err := r.upstreams[p.handler](s, d, data)
				if err != nil {
					log.Println("RPC: Proxy error:", err)
				}
				m.Ack(false)
			}

		}(d)
	}

	fmt.Println("handle: deliveries channel closed")
	done <- nil
}
