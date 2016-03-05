package rpc

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestCallMessage(t *testing.T) {

	var (
		name  = "test-call"
		uuid  = "uuid"
		token = "token"
		a     = 0
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)
	if r.uri != uri {
		t.Error("Expected uri: %s got %s", uri, r.uri)
	}

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	m := struct{ Name string }{"name"}

	handler := func(s Sender, p []byte) error {
		log.Print("received", p)

		if a <= 4 {

			r.Call(d, m)
			a++
			return nil
		}

		i := struct {
			Name string
		}{}

		err = json.Unmarshal(p, &i)
		if err != nil {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if m.Name != i.Name {
			t.Error("Received message validation failed: %x got %x", m, i)
		}
		end <- true
		return nil
	}
	r.SetHandler("handler", handler)

	r.Listen()
	t.Log("RPC registered and setuped")

	timer := time.NewTimer(time.Second * 5)

	for {
		select {
		case <-r.connected:
			r.Call(d, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestCastMessage(t *testing.T) {

	var (
		name  = "test-case"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	m := struct {
		Name string
	}{
		"name",
	}

	handler := func(s Sender, p []byte) error {
		log.Print("received", p)

		i := struct {
			Name string
		}{}

		err = json.Unmarshal(p, &i)
		if err != nil {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if m.Name != i.Name {
			t.Error("Received message validation failed: %x got %x", m, i)
		}
		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second * 5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.Cast(d, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestCallBinaryMessage(t *testing.T) {

	var (
		name  = "test-call-b"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	m := []byte{123, 125}

	handler := func(s Sender, p []byte) error {
		log.Print("received", p)

		if string(m) != string(p) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second * 5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.CallBinary(d, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}

}

func TestCastBinaryMessage(t *testing.T) {

	var (
		name  = "test-cast-b"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)
	if r.uri != uri {
		t.Error("Expected uri: %s got %s", uri, r.uri)
	}

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	m := []byte{123, 125}

	handler := func(s Sender, p []byte) error {
		log.Print("received", p)

		if string(m) != string(p) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second * 5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.CastBinary(d, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestCallSignedMessage(t *testing.T) {

	var (
		name  = "test-call"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)
	if r.uri != uri {
		t.Error("Expected uri: %s got %s", uri, r.uri)
	}

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	m := struct{ Name string }{"name"}
	so := Sender{
		Name: "sender",
		UUID: "suid",
	}

	handler := func(s Sender, p []byte) error {
		log.Print("received", p)

		i := struct {
			Name string
		}{}

		err = json.Unmarshal(p, &i)
		if err != nil {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if m.Name != i.Name {
			t.Error("Received message validation failed: %x got %x", m, i)
		}

		if s.Name != so.Name {
			t.Error("Received message validation failed: %x got %x", so.Name, s.Name)
		}

		end <- true
		return nil
	}
	r.SetHandler("handler", handler)

	r.Listen()
	t.Log("RPC registered and setuped")

	timer := time.NewTimer(time.Second * 5)

	for {
		select {
		case <-r.connected:
			r.CallSigned(so, d, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestCastSignedMessage(t *testing.T) {

	var (
		name  = "test-case"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	m := struct{ Name string }{"name"}
	so := Sender{
		Name: "sender",
		UUID: "suid",
	}

	handler := func(s Sender, p []byte) error {
		log.Print("received", p)

		i := struct {
			Name string
		}{}

		err = json.Unmarshal(p, &i)
		if err != nil {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if m.Name != i.Name {
			t.Error("Received message validation failed: %x got %x", m, i)
		}
		if so.Name != s.Name {
			t.Error("Received message validation failed: %x got %x", so.Name, s.Name)
		}

		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second * 5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.CastSigned(so, d, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestCallSignedBinaryMessage(t *testing.T) {

	var (
		name  = "test-call-b"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.Shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	m := []byte{123, 125}
	so := Sender{
		Name: "sender",
		UUID: "suid",
	}

	handler := func(s Sender, p []byte) error {
		log.Print("received", p)

		if string(m) != string(p) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		if so.Name != s.Name {
			t.Error("Received message validation failed: %x got %x", so.Name, s.Name)
		}
		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second * 5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.CallSignedBinary(so, d, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}

}

func TestCastSignedBinaryMessage(t *testing.T) {

	var (
		name  = "test-cast-b"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.Shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)
	if r.uri != uri {
		t.Error("Expected uri: %s got %s", uri, r.uri)
	}

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	m := []byte{123, 125}
	so := Sender{
		Name: "sender",
		UUID: "suid",
	}

	handler := func(s Sender, p []byte) error {
		log.Print("received", p)

		if string(m) != string(p) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if so.Name != s.Name {
			t.Error("Received message validation failed: %x got %x", so.Name, s.Name)
		}

		end <- true
		return nil
	}

	r.SetHandler("handler", handler)
	r.Listen()

	timer := time.NewTimer(time.Second * 5)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.CastSignedBinary(so, d, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestProxyCall(t *testing.T) {

	var (
		name  = "test-proxy-cl"
		uuid  = "uuid"
		token = "G&vW8ag#2q8MU&h78@Gu^pjtg@R5CPcd"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		//r.cleanup()
		r.Shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}

	p := Receiver{
		Name:    name,
		UUID:    uuid,
		Handler: "proxy",
	}

	m := struct{ Name string }{"name"}

	proxy := func(s Sender, g Destination, b []byte) error {
		t.Log("received in proxy", string(b), s, g)
		log.Println(g)

		r.CallBinary(g, b)
		return nil
	}
	r.SetUpstream("proxy", proxy)

	h := func(s Sender, p []byte) error {
		t.Log("received", p)

		i := struct {
			Name string
		}{}

		err = json.Unmarshal(p, &i)
		if err != nil {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if m.Name != i.Name {
			t.Error("Received message validation failed: %x got %x", m, i)
		}

		end <- true
		return nil
	}
	r.SetHandler("handler", h)

	r.Listen()

	timer := time.NewTimer(time.Second * 10)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.ProxyCall(d, p, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestProxyCast(t *testing.T) {

	var (
		name  = "test-proxy-ct"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.Shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}
	p := Receiver{
		Name:    name,
		UUID:    uuid,
		Handler: "proxy",
	}

	m := struct{ Name string }{"name"}

	proxy := func(s Sender, d Destination, b []byte) error {
		t.Log("received in proxy", b)
		r.CallBinary(d, b)
		return nil
	}

	r.SetUpstream("proxy", proxy)

	handler := func(s Sender, p []byte) error {
		t.Log("received", p)
		i := struct {
			Name string
		}{}

		err = json.Unmarshal(p, &i)
		if err != nil {
			t.Error("Received message validation failed: %x got %x", m, p)
		}

		if m.Name != i.Name {
			t.Error("Received message validation failed: %x got %x", m, i)
		}
		end <- true
		return nil
	}
	r.SetHandler("handler", handler)

	r.Listen()

	timer := time.NewTimer(time.Second * 20)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.ProxyCast(d, p, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestProxyCallBinary(t *testing.T) {

	var (
		name  = "test-proxy-clb"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.Shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}
	p := Receiver{
		Name:    name,
		UUID:    uuid,
		Handler: "proxy",
	}

	m := []byte{123, 125}

	proxy := func(s Sender, d Destination, b []byte) error {
		r.CallBinary(d, b)
		return nil
	}

	r.SetUpstream("proxy", proxy)

	handler := func(s Sender, p []byte) error {
		t.Log("received", p)

		if string(m) != string(p) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		end <- true
		return nil
	}
	r.SetHandler("handler", handler)

	r.Listen()

	timer := time.NewTimer(time.Second * 20)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.ProxyCallBinary(d, p, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}

func TestProxyCastBinary(t *testing.T) {

	var (
		name  = "test-proxy-ctb"
		uuid  = "uuid"
		token = "token"
	)

	r, err := Register(name, uuid, token)
	if err != nil {
		t.Error("Register APP error", err)
	}

	defer func() {
		r.cleanup()
		r.Shutdown()
	}()

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)

	end := make(chan bool)
	d := Destination{
		Name:    name,
		UUID:    uuid,
		Handler: "handler",
	}
	p := Receiver{
		Name:    name,
		UUID:    uuid,
		Handler: "proxy",
	}

	m := []byte{123, 125}

	proxy := func(s Sender, d Destination, b []byte) error {
		r.CastBinary(d, b)
		return nil
	}

	r.SetUpstream("proxy", proxy)

	handler := func(s Sender, p []byte) error {
		t.Log("received", p)

		if string(m) != string(p) {
			t.Error("Received message validation failed: %x got %x", m, p)
		}
		end <- true
		return nil
	}
	r.SetHandler("handler", handler)

	r.Listen()

	timer := time.NewTimer(time.Second * 20)
	t.Log("RPC registered and setuped")

	for {
		select {
		case <-r.connected:
			r.ProxyCastBinary(d, p, m)
		case <-end:
			return
		case <-timer.C:
			t.Error("No message received: failed")
			timer.Stop()
			return
		}
	}
}
