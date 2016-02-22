package rpc

import (
	"encoding/json"
	"log"
)

// Call - send message with delivery guarantee
func (r *RPC) Call (d Receiver, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.call(Sender{r.name, r.uuid}, Proxy{d.name, d.uuid, d.handler}, d, msg);
	if err != nil {
		return err
	}

	return nil
}

// Cast - send message without delivery guarantee
func (r *RPC) Cast (d Receiver, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.cast(Sender{r.name, r.uuid}, Proxy{d.name, d.uuid, d.handler}, d, msg);
	if err != nil {
		return err
	}

	return nil
}

// CallBinary - send binary message with delivery guarantee
func (r *RPC) CallBinary (d Receiver, message []byte) error {
	err := r.call(Sender{r.name, r.uuid}, Proxy{d.name, d.uuid, d.handler}, d, message);
	if err != nil {
		return err
	}

	return nil
}

// CastBinary - send binary message without delivery guarantee
func (r *RPC) CastBinary (d Receiver, message []byte) error {
	err := r.cast(Sender{r.name, r.uuid}, Proxy{d.name, d.uuid, d.handler}, d, message);
	if err != nil {
		return err
	}

	return nil
}

// CallSigned - send signed message with delivery guarantee
func (r *RPC) CallSigned (s Sender, d Receiver, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.call(s, Proxy{d.name, d.uuid, d.handler,}, d, msg);
	if err != nil {
		return err
	}

	return nil
}

// CastSigned - send signed message without delivery guarantee
func (r *RPC) CastSigned (s Sender, d Receiver, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.cast(s, Proxy{d.name, d.uuid, d.handler,}, d, msg);
	if err != nil {
		return err
	}

	return nil
}

// CallSignedBinary - send signed binary message with delivery guarantee
func (r *RPC) CallSignedBinary (s Sender, d Receiver, message []byte) error {

	err := r.call(s, Proxy{d.name, d.uuid, d.handler,}, d, message);
	if err != nil {
		return err
	}

	return nil
}

// CastSignedBinary - send signed binary message without delivery guarantee
func (r *RPC) CastSignedBinary (s Sender, d Receiver, message []byte) error {

	err := r.cast(s, Proxy{d.name, d.uuid, d.handler,}, d, message);
	if err != nil {
		return err
	}

	return nil
}

// Proxy send message methods
// ProxyCall - send message throw another application with delivery guarantee
func (r *RPC) ProxyCall (p Proxy, d Receiver, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.call(Sender{r.name, r.uuid}, p, d, msg);
	if err != nil {
		return err
	}

	return nil
}

// ProxyCast - send message throw another application without delivery guarantee
func (r *RPC) ProxyCast (p Proxy, d Receiver, message interface{}) error {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println("RPC: Protocol encode error")
		return err
	}

	err = r.cast(Sender{r.name, r.uuid}, p, d, msg);
	if err != nil {
		return err
	}

	return nil
}

// ProxyCallBinary - send binary message throw another application with delivery guarantee
func (r *RPC) ProxyCallBinary (p Proxy, d Receiver, message []byte) error {

	err := r.call(Sender{r.name, r.uuid}, p, d, message);
	if err != nil {
		return err
	}

	return nil
}

// ProxyCastBinary - send message throw another application without delivery guarantee
func (r *RPC) ProxyCastBinary (p Proxy, d Receiver, message []byte) error {

	err := r.cast(Sender{r.name, r.uuid}, p, d, message);
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

// Shutdown RPC instance
func (r *RPC) Shutdown () {

}