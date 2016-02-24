package rpc

import (
	"testing"
)

func TestSenderSign (t *testing.T) {
	s := Sender{
		name: "demo",
		uuid: "uuid",
	}

	sign := s.Sign()
	data := []byte{
		100, 101, 109, 111, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		117, 117, 105, 100, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0,
	}

	if string(sign) != string(data)  {
		t.Error("Failed signing sender: expected %x, got %x", data, sign)
	}
}

func TestProxySign (t *testing.T) {
	p := Receiver{
		name: "demo",
		uuid: "uuid",
		handler: "handler",
	}

	sign := p.Sign()
	data := []byte{
		100, 101, 109, 111, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,

		117, 117, 105, 100, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0,

		104, 97, 110, 100, 108, 101, 114, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	if string(sign) != string(data)  {
		t.Error("Failed signing proxy: expected %x, got %x", data, sign)
	}
}

func TestReceiverSign (t *testing.T) {
	r := Destination{
		name: "demo",
		uuid: "uuid",
		handler: "handler",
	}

	sign := r.Sign()
	data := []byte{
		100, 101, 109, 111, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,

		117, 117, 105, 100, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0,

		104, 97, 110, 100, 108, 101, 114, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	if string(sign) != string(data)  {
		t.Error("Failed signing proxy: expected %x, got %x", data, sign)
	}
}

func TestEncode (t *testing.T) {

	r := RPC{}
	r.token = "token"

	s := Sender{
		name: "demo",
		uuid: "uuid",
	}


	d := Destination{
		name: "demo",
		uuid: "uuid",
		handler: "handler",
	}

	p := Receiver{}



	body := r.encode(s, d, p, []byte{123,125})

	data := []byte{

		116, 111, 107, 101, 110, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,

		100, 101, 109, 111, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		117, 117, 105, 100, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0,


		100, 101, 109, 111, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,

		117, 117, 105, 100, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0,

		104, 97, 110, 100, 108, 101, 114, 0,
		0, 0, 0, 0, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,

		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0,

		123, 125,

	}

	t.Log(body)
	t.Log(data)

	if string(body) != string(data)  {
		t.Error("Failed signing proxy: expected %x, got %x", data, body)
	}
}