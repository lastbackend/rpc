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
	var attempt int

	for {
		select {
		case <-r.connect:
			if attempt >= 60 {
				log.Println("attempt limit reached: 5")
				return
			}
			attempt++
			log.Printf("RPC: %s connect", r.name)
			go r.dial()

		case <-r.reconnect:
			log.Printf("RPC: %s reconnect", r.name)
			timer := time.NewTimer(time.Second)
			<-timer.C
			go r.dial()
			timer.Stop()
		case <- r.done :
			return
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

	log.Println("RPC: Dial to:", r.uri)
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

	r.connected <- true
	r.subscribe()

}

func (r *RPC) call (s Sender, d Destination, p Receiver, data []byte) error {
	return r.publish(true, s, d, p, data)
}

func (r *RPC) cast (s Sender, d Destination, p Receiver, data []byte) error {
	return r.publish(false, s, d, p, data)
}

func (r *RPC) publish (call bool, s Sender, d Destination, p Receiver, data []byte) error {

	body := r.encode(s, d, p, data)

	log.Printf("PRC: publish to %s:%s, dest: %s:%s, send: %dB body (%s)", d.name, d.uuid, p.name, p.uuid, len(body), body)

	channel, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	exchange := fmt.Sprintf("%s:%s", d.name, "direct")
	if !call {
		exchange = fmt.Sprintf("%s:%s", d.name, "topic")
	}

	bind     := d.uuid
	if bind == "" {
		bind = d.name
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

	r.exchanges.direct = fmt.Sprintf("%s:%s", r.name, "direct")
	r.exchanges.topic = fmt.Sprintf("%s:%s", r.name, "topic")


	r.queues.direct = fmt.Sprintf("%s:%s", r.name, "direct")
	r.queues.topic = fmt.Sprintf("%s:%s", r.name, "topic")

	// Get hostname for register current instance
	log.Printf("RPC: Create new consumer: %s", r.name)

	r.channels.direct, err = r.conn.Channel();
	if err != nil {
		log.Println("Channel:", err)
		return err
	}

	err = r.channels.direct.Qos(100, 0, true)
	if err != nil {
		log.Println("Channel:", err)
		return err
	}

	r.channels.topic, err = r.conn.Channel();
	if err != nil {
		log.Println("Channel:", err)
		return err
	}

	err = r.channels.topic.Qos(100, 0, true)
	if err != nil {
		log.Println("Channel:", err)
		return err
	}

	// create direct exchange for guarantee delivery messages
	if err = r.channels.direct.ExchangeDeclare(r.exchanges.direct, "direct", true, false, false, false, nil); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// create topic exchange for non guarantee delivery messages
	if err = r.channels.topic.ExchangeDeclare(r.exchanges.topic, "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// create direct queue for guarantee delivery messages
	if _, err := r.channels.direct.QueueDeclare(r.queues.direct, true, false, false, false, nil); err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	// create topic queue for non guarantee delivery messages
	if _, err := r.channels.topic.QueueDeclare(r.queues.topic, true, true, false, false, nil); err != nil {
		return fmt.Errorf("Queue Declare: %s", err)
	}

	// create bindings for direct messages
	if err = r.channels.direct.QueueBind(r.queues.direct, r.uuid, r.exchanges.direct, false, nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if err = r.channels.direct.QueueBind(r.queues.direct, r.name, r.exchanges.direct, false, nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	// create bindings for topic messages
	if err = r.channels.direct.QueueBind(r.queues.topic, r.uuid, r.exchanges.direct, false, nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}

	if err = r.channels.direct.QueueBind(r.queues.topic, r.name, r.exchanges.direct, false, nil); err != nil {
		return fmt.Errorf("Queue Bind: %s", err)
	}


	messages, err := r.channels.direct.Consume(r.queues.direct, r.queues.direct, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Queue Consume: %s", err)
	}

	streams, err := r.channels.topic.Consume(r.queues.topic, r.queues.topic, false, false, false, false, nil)
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

			log.Println("receiver message from:", m.DeliveryTag, m.ConsumerTag, string(m.Body))
			// validate token
			token := string(m.Body[0:32])
			token  = token[:strings.Index(string(token), "\x00")]

			log.Println(token, r.token)

			if r.token != token {
				log.Println("RPC: token verification failed")
				m.Ack(false)
			}

			// parse sender information
			s := Sender{}
			s.name = string(d.Body[32:48])
			s.uuid = string(d.Body[48:84])

			s.name = s.name[:strings.Index(string(s.name), "\x00")]

			// parse destination information
			e := Destination{}
			e.name = string(d.Body[84:100])
			e.uuid = string(d.Body[100:136])
			e.handler = string(d.Body[136:152])

			e.name = e.name[:strings.Index(string(e.name), "\x00")]
			e.handler = e.handler[:strings.Index(string(e.handler), "\x00")]

			// parse receiver information
			p := Receiver{}
			p.name = string(m.Body[152:168])
			p.uuid = string(m.Body[168:204])
			p.handler = string(m.Body[204:220])

			p.name = p.name[:strings.Index(string(p.name), "\x00")]
			p.handler = p.handler[:strings.Index(string(p.handler), "\x00")]

			data := m.Body[256:len(m.Body)]

			if p.name != "" {
				log.Println("PRC: need upstream", m.ConsumerTag)
				_, ok := r.upstreams[p.handler]

				if !ok {
					log.Println("RPC: upstream not found", p.handler)
					m.Ack(false)
				}

				err := r.upstreams[p.handler](s, e, data)
				if err != nil {
					log.Println("RPC: Proxy error:", err)
				}

				m.Ack(false)
				return
			}

				log.Println("PRC: send to handler", m.ConsumerTag)

				_, ok := r.handlers[e.handler]

				if !ok {
					log.Println("RPC: handler not found", e.handler)
					m.Ack(false)
					return
				}

				err := r.handlers[e.handler](s, data)
				if err != nil {
					log.Println("RPC: Proxy error:", err)
				}

				m.Ack(false)
		}(d)
	}

	fmt.Println("handle: deliveries channel closed")
	r.done <- nil
	return
}

func (r *RPC) shutdown () error {
	// will close() the deliveries channel
	log.Println("RPC: Shutdown broker")
	// close direct channels
	if err := r.channels.direct.Cancel(r.queues.direct, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := r.channels.topic.Cancel(r.queues.topic, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer fmt.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-r.done
}
