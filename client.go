package rpc

import (
	"encoding/json"
	"log"
)

// Call - send message with delivery guarantee
func (r *RPC) Call (d Destination, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.call(Sender{r.name, r.uuid}, d, Receiver{}, msg);
	if err != nil {
		return err
	}

	return nil
}

// Cast - send message without delivery guarantee
func (r *RPC) Cast (d Destination, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.cast(Sender{r.name, r.uuid}, d, Receiver{}, msg);
	if err != nil {
		return err
	}

	return nil
}

// CallBinary - send binary message with delivery guarantee
func (r *RPC) CallBinary (d Destination, message []byte) error {

	err := r.call(Sender{r.name, r.uuid}, d, Receiver{}, message);
	if err != nil {
		return err
	}

	return nil
}

// CastBinary - send binary message without delivery guarantee
func (r *RPC) CastBinary (d Destination, message []byte) error {
	err := r.cast(Sender{r.name, r.uuid}, d, Receiver{}, message);
	if err != nil {
		return err
	}

	return nil
}

// CallSigned - send signed message with delivery guarantee
func (r *RPC) CallSigned (s Sender, d Destination, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.call(s, d, Receiver{}, msg);
	if err != nil {
		return err
	}

	return nil
}

// CastSigned - send signed message without delivery guarantee
func (r *RPC) CastSigned (s Sender, d Destination, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.cast(s, d, Receiver{}, msg);
	if err != nil {
		return err
	}

	return nil
}

// CallSignedBinary - send signed binary message with delivery guarantee
func (r *RPC) CallSignedBinary (s Sender, d Destination, message []byte) error {

	err := r.call(s, d, Receiver{}, message);
	if err != nil {
		return err
	}

	return nil
}

// CastSignedBinary - send signed binary message without delivery guarantee
func (r *RPC) CastSignedBinary (s Sender, d Destination, message []byte) error {

	err := r.cast(s, d, Receiver{}, message);
	if err != nil {
		return err
	}

	return nil
}

// Proxy send message methods
// ProxyCall - send message throw another application with delivery guarantee
func (r *RPC) ProxyCall (d Destination, p Receiver, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.call(Sender{r.name, r.uuid}, d, p, msg);
	if err != nil {
		return err
	}

	return nil
}

// ProxyCast - send message throw another application without delivery guarantee
func (r *RPC) ProxyCast (d Destination, p Receiver, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.cast(Sender{r.name, r.uuid}, d, p, msg);
	if err != nil {
		return err
	}

	return nil
}

// ProxyCallBinary - send binary message throw another application with delivery guarantee
func (r *RPC) ProxyCallBinary (d Destination, p Receiver, message []byte) error {

	err := r.call(Sender{r.name, r.uuid}, d, p, message);
	if err != nil {
		return err
	}

	return nil
}

// ProxyCastBinary - send message throw another application without delivery guarantee
func (r *RPC) ProxyCastBinary (d Destination, p Receiver, message []byte) error {

	err := r.cast(Sender{r.name, r.uuid}, d, p, message);
	if err != nil {
		return err
	}

	return nil
}


// SetHeader - set handler routing
func (r *RPC) SetHandler(h string, f Handler) {
	r.handlers[h]=f
}

// SetUpstream - set upstream routing
func (r *RPC) SetUpstream(u string, f Upstream) {
	r.upstreams[u]=f
}