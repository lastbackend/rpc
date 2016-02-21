package rpc

func (s *Sender) Sign() []byte {
	var app [32]byte
	var tag [36]byte

	copy(app[:], s.name)
	copy(tag[:], s.uuid)

	var body []byte

	body = append(body[:], app[:]...)
	body = append(body[:], tag[:]...)

	return body
}

func (p *Proxy) Sign () []byte {

	var body []byte

	var name [16]byte
	var uuid [36]byte
	var hander [16]byte

	copy(name[:], p.name)
	copy(uuid[:], p.uuid)
	copy(hander[:], p.handler)

	body = append(body, name[:]...)
	body = append(body, uuid[:]...)
	body = append(body, hander[:]...)

	return body
}

func (r *Receiver) Sign () []byte {

	var body []byte

	var name [16]byte
	var uuid [36]byte
	var hander [16]byte

	copy(name[:], r.name)
	copy(uuid[:], r.uuid)
	copy(hander[:], r.handler)

	body = append(body, name[:]...)
	body = append(body, uuid[:]...)
	body = append(body, hander[:]...)

	return body
}