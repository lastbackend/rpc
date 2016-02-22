package rpc

import (
	"testing"
	"fmt"
	"log"
	"time"
	"encoding/json"
)


func TestCallMessage (t *testing.T) {
	var err error


	_, err = Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	rpc = SetURI(uri)
	if rpc.uri != uri {
		t.Error("Expected uri: %s got %s", uri, rpc.uri)
	}

	Listen()

	<- rpc.connected
	end := make (chan bool)
	r := Destination{
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

	rpc.SetHandler("handler", handler)

	timer := time.NewTimer(time.Second*5)
	rpc.Call(r, m)

	select {
	case <- end:
		return
	case <-timer.C:
		t.Error("No message received: failed")
		timer.Stop()
	}

}

func TestCastMessage (t *testing.T) {
	var err error


	_, err = Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	rpc = SetURI(uri)
	if rpc.uri != uri {
		t.Error("Expected uri: %s got %s", uri, rpc.uri)
	}

	Listen()

	<- rpc.connected
	end := make (chan bool)
	r := Destination{
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

	rpc.SetHandler("handler", handler)

	timer := time.NewTimer(time.Second*5)
	rpc.Cast(r, m)

	select {
	case <- end:
		return
	case <-timer.C:
		t.Error("No message received: failed")
		timer.Stop()
	}

}

func TestCallBinaryMessage (t *testing.T) {
	var err error


	_, err = Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	rpc = SetURI(uri)
	if rpc.uri != uri {
		t.Error("Expected uri: %s got %s", uri, rpc.uri)
	}

	Listen()

	<- rpc.connected
	end := make (chan bool)
	r := Destination{
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

	rpc.SetHandler("handler", handler)

	timer := time.NewTimer(time.Second*5)
	rpc.CallBinary(r, m)

	select {
	case <- end:
		return
	case <-timer.C:
		t.Error("No message received: failed")
		timer.Stop()
	}

}

func TestCastBinaryMessage (t *testing.T) {
	var err error


	_, err = Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	rpc = SetURI(uri)
	if rpc.uri != uri {
		t.Error("Expected uri: %s got %s", uri, rpc.uri)
	}

	Listen()

	<- rpc.connected
	end := make (chan bool)
	r := Destination{
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

	rpc.SetHandler("handler", handler)

	timer := time.NewTimer(time.Second*5)
	rpc.CastBinary(r, m)

	select {
	case <- end:
		return
	case <-timer.C:
		t.Error("No message received: failed")
		timer.Stop()
	}

}

func TestProxyCallBinary (t *testing.T) {
	var err error

	_, err = Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	rpc = SetURI(uri)
	if rpc.uri != uri {
		t.Error("Expected uri: %s got %s", uri, rpc.uri)
	}

	Listen()

	<- rpc.connected

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

	proxy := func (s Sender, d Destination, p []byte) error {
		log.Print("received in proxy", p)

		if (string(m) != string(p)) {
			t.Error("Received message validation in proxy failed: %x got %x", m, p)
		}

		rpc.CallBinary(d, p)

		return nil
	}
	rpc.SetUpstream("proxy", proxy)

	handler := func (s Sender, p []byte) error {
		log.Print("received", p)

		if (string(m) != string(p)) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		end <- true
		return nil
	}
	rpc.SetHandler("handler", handler)

	timer := time.NewTimer(time.Second*5)
	rpc.ProxyCallBinary(d, p, m)

	select {
	case <- end:
		return
	case <-timer.C:
		t.Error("No message received: failed")
		timer.Stop()
	}

}