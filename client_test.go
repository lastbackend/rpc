package rpc

import (
	"testing"
	"fmt"
	"log"
	"time"
	"encoding/json"
)


func TestCallMessage (t *testing.T) {

	var (
		name  = "test"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer r.shutdown()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)
	if r.uri != uri {
		t.Error("Expected uri: %s got %s", uri, r.uri)
	}

	end := make (chan bool)
	d := Destination{
		name: "test",
		uuid: "uuid",
		handler: "handler",
	}

	m := struct{ Name string }{"name",}

	handler := func (s Sender, p []byte) error {
		log.Print("received", p)

		i := struct{
			Name string
		}{

		}

		err = json.Unmarshal(p, &i)
		if err != nil {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if (m.Name != i.Name) {
			t.Error("Received message validation failed: %x got %x", m, i)
		}
		end <- true
		return nil
	}
	r.SetHandler("handler", handler)

	r.Listen()
	t.Log("RPC registered and setuped")

	timer := time.NewTimer(time.Second*5)

	for {
		select {
		case <- r.connected:
			r.Call(d, m)
		case <- end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestCastMessage (t *testing.T) {

	var (
		name  = "test"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}
	defer r.shutdown()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make (chan bool)
	d := Destination{
		name: "test",
		uuid: "uuid",
		handler: "handler",
	}

	m := struct{
		Name string
	}{
		"name",
	}

	handler := func (s Sender, p []byte) error {
		log.Print("received", p)

		i := struct{
			Name string
		}{

		}

		err = json.Unmarshal(p, &i)
		if err != nil {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if (m.Name != i.Name) {
			t.Error("Received message validation failed: %x got %x", m, i)
		}
		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second*5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <- r.connected:
			r.Cast(d, m)
		case <- end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestCallBinaryMessage (t *testing.T) {

	var (
		name  = "test"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}
	defer r.shutdown()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make (chan bool)
	d := Destination{
		name: "test",
		uuid: "uuid",
		handler: "handler",
	}

	m := []byte{123,125}

	handler := func (s Sender, p []byte) error {
		log.Print("received", p)

		if (string(m) != string(p)) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second*5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <- r.connected:
			r.CallBinary(d, m)
		case <- end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}

}

func TestCastBinaryMessage (t *testing.T) {

	var (
		name  = "test"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}
	defer r.shutdown()


	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)
	if r.uri != uri {
		t.Error("Expected uri: %s got %s", uri, r.uri)
	}

	end := make (chan bool)
	d := Destination{
		name: "test",
		uuid: "uuid",
		handler: "handler",
	}

	m := []byte{123,125}

	handler := func (s Sender, p []byte) error {
		log.Print("received", p)

		if (string(m) != string(p)) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second*5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <- r.connected:
			r.CastBinary(d, m)
		case <- end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestProxyCallBinary (t *testing.T) {

	var (
		name  = "test"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}
	defer r.shutdown()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make (chan bool)
	d := Destination{
		name: "test",
		uuid: "uuid",
		handler: "handler",
	}
	p := Receiver{
		name: "test",
		uuid: "uuid",
		handler: "proxy",
	}

	m := []byte{123,125}

	proxy := func (s Sender, d Destination, b []byte) error {
		log.Print("received in proxy", b)

		if (string(m) != string(b)) {
			t.Error("Received message validation in proxy failed: %x got %x", m, p)
		}

		log.Println(s, d, b)
		log.Println("Call Binary data:", d, b)
		go r.CallBinary(d, b)

		return nil
	}
	r.SetUpstream("proxy", proxy)

	handler := func (s Sender, p []byte) error {
		log.Print("received", p)

		if (string(m) != string(p)) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		end <- true
		return nil
	}
	r.SetHandler("handler", handler)

	r.Listen()

	timer := time.NewTimer(time.Second*5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <- r.connected:
			r.ProxyCallBinary(d, p, m)
		case <- end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}