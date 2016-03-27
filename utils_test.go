package rpc

import (
	log "github.com/Sirupsen/logrus"
	"testing"
)

func TestSenderSign(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	s := Sender{
		Name: "demo",
		UUID: "uuid",
	}

	sign, err := s.Sign()
	data := []byte{52, 0, 0, 52, 0, 0, 100, 101, 109, 111, 117, 117, 105, 100}

	if err != nil {
		t.Error("Failed signing sender:", err)
	}
	if string(sign) != string(data) {
		t.Error("Failed signing sender: expected %x, got %x", data, sign)
	}
}

func TestProxySign(t *testing.T) {
	p := Receiver{
		Name:    "demo",
		UUID:    "uuid",
		Handler: "handler",
	}

	sign, err := p.Sign()
	data := []byte{52, 0, 0, 52, 0, 0, 55, 0, 0,
		100, 101, 109, 111,
		117, 117, 105, 100,
		104, 97, 110, 100, 108, 101, 114,
	}

	if err != nil {
		t.Error("Failed signing proxy:", err)
	}

	if string(sign) != string(data) {
		t.Error("Failed signing proxy: expected %x, got %x", data, sign)
	}
}

func TestReceiverSign(t *testing.T) {
	r := Destination{
		Name:    "demo",
		UUID:    "uuid",
		Handler: "handler",
	}

	sign, err := r.Sign()
	data := []byte{
		52, 0, 0, 52, 0, 0, 55, 0, 0,

		100, 101, 109, 111,
		117, 117, 105, 100,
		104, 97, 110, 100, 108, 101, 114,
	}

	if err != nil {
		t.Error("Failed signing dest:", err)
	}

	if string(sign) != string(data) {
		t.Error("Failed signing proxy: expected %x, got %x", data, sign)
	}
}

func TestEncode(t *testing.T) {

	r := RPC{}
	r.token = "token"

	s := Sender{
		Name: "demo",
		UUID: "uuid",
	}

	d := Destination{
		Name:    "demo",
		UUID:    "uuid",
		Handler: "handler",
	}

	p := Receiver{}

	body, err := r.encode(s, d, p, []byte{123, 125})

	data := []byte{

		53, 0,
		116, 111, 107, 101, 110,

		52, 0, 0, 52, 0, 0,
		100, 101, 109, 111,
		117, 117, 105, 100,

		52, 0, 0, 52, 0, 0, 55, 0, 0,
		100, 101, 109, 111,
		117, 117, 105, 100,
		104, 97, 110, 100, 108, 101, 114,
		48, 0, 0, 48, 0, 0, 48, 0, 0,
		123, 125,
	}

	t.Log(body)
	t.Log(data)

	if err != nil {
		t.Error("Failed encode:", err)
	}

	if string(body) != string(data) {
		t.Error("Failed signing proxy: expected %x, got %x", data, body)
	}
}
