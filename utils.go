package rpc

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"strconv"
	"strings"
)

var (
	ERRINVALIDLENGTH = errors.New("Invalid length of message field")
	ERRINVALIDTOKEN  = errors.New("Invalid authentication token")
)

func (s *Sender) Sign() ([]byte, error) {

	var body []byte

	if (len(s.Name) > 255) || (len(s.UUID) > 255) {
		return body, ERRINVALIDLENGTH
	}

	cn := [3]byte{}
	cu := [3]byte{}

	copy(cn[:], strconv.Itoa(len(s.Name)))
	copy(cu[:], strconv.Itoa(len(s.UUID)))

	body = append(body[:], cn[:]...)
	body = append(body[:], cu[:]...)

	body = append(body[:], []byte(s.Name)[:]...)
	body = append(body[:], []byte(s.UUID)[:]...)

	return body, nil
}

func (p *Receiver) Sign() ([]byte, error) {

	var body []byte
	if (len(p.Name) > 255) || (len(p.UUID) > 255) || (len(p.Handler) > 255) {
		return body, ERRINVALIDLENGTH
	}

	cn := [3]byte{}
	cu := [3]byte{}
	ch := [3]byte{}

	copy(cn[:], strconv.Itoa(len(p.Name)))
	copy(cu[:], strconv.Itoa(len(p.UUID)))
	copy(ch[:], strconv.Itoa(len(p.Handler)))

	body = append(body[:], cn[:]...)
	body = append(body[:], cu[:]...)
	body = append(body[:], ch[:]...)

	body = append(body[:], []byte(p.Name)[:]...)
	body = append(body[:], []byte(p.UUID)[:]...)
	body = append(body[:], []byte(p.Handler)[:]...)

	return body, nil
}

func (d *Destination) Sign() ([]byte, error) {

	var body []byte

	if (len(d.Name) > 255) || (len(d.UUID) > 255) || (len(d.Handler) > 255) {
		return body, ERRINVALIDLENGTH
	}

	cn := [3]byte{}
	cu := [3]byte{}
	ch := [3]byte{}

	copy(cn[:], strconv.Itoa(len(d.Name)))
	copy(cu[:], strconv.Itoa(len(d.UUID)))
	copy(ch[:], strconv.Itoa(len(d.Handler)))

	body = append(body[:], cn[:]...)
	body = append(body[:], cu[:]...)
	body = append(body[:], ch[:]...)

	body = append(body[:], []byte(d.Name)[:]...)
	body = append(body[:], []byte(d.UUID)[:]...)
	body = append(body[:], []byte(d.Handler)[:]...)

	return body, nil
}

func (r *RPC) encode(s Sender, d Destination, p Receiver, data []byte) ([]byte, error) {
	var body []byte

	if len(r.token) > 96 {
		return body, ERRINVALIDLENGTH
	}

	ct := [2]byte{}
	copy(ct[:], strconv.Itoa(len(r.token)))
	body = append(body[:], ct[:]...)
	body = append(body[:], []byte(r.token)[:]...)

	sender, err := s.Sign()
	if err != nil {
		return body, err
	}
	body = append(body, sender[:]...)

	destinaiton, err := d.Sign()
	if err != nil {
		return body, err
	}
	body = append(body, destinaiton[:]...)

	receiver, err := p.Sign()
	if err != nil {
		return body, err
	}
	body = append(body, receiver[:]...)
	body = append(body, data[:]...)

	return body, nil
}

func (r *RPC) decode(data []byte) (Sender, Destination, Receiver, []byte, error) {

	s := Sender{}
	d := Destination{}
	p := Receiver{}

	if len(data) == 0 {
		return s, d, p, []byte{}, errors.New("Body is empty")
	}

	tc, err := r.parseInt(data[0:2], 2)
	if err != nil {
		log.Error("Parse message error:", err)
		return s, d, p, []byte{}, err
	}

	var token = string(data[2:tc])

	if token != r.token {
		return s, d, p, []byte{}, ERRINVALIDTOKEN
	}

	data = data[tc:]
	snl, _ := r.parseInt(data[0:3], 6)
	sul, _ := r.parseInt(data[3:6], snl)

	s.Name = string(data[6:snl])
	s.UUID = string(data[snl:sul])

	data = data[sul:]

	dnl, _ := r.parseInt(data[0:3], 9)
	dul, _ := r.parseInt(data[3:6], dnl)
	dhl, _ := r.parseInt(data[6:9], dul)

	if dnl > 9 {
		d.Name = string(data[9:dnl])
	}
	if dul > dnl {
		d.UUID = string(data[dnl:dul])
	}

	if dhl > dul {
		d.Handler = string(data[dul:dhl])
	}

	data = data[dhl:]

	pnl, _ := r.parseInt(data[0:3], 9)
	pul, _ := r.parseInt(data[3:6], pnl)
	phl, _ := r.parseInt(data[6:9], pul)

	if pnl > 9 {
		p.Name = string(data[9:pnl])
	}
	if pul > pnl {
		p.UUID = string(data[pnl:pul])
	}
	if phl > pul {
		p.Handler = string(data[pul:phl])
	}
	data = data[phl:]

	return s, d, p, data, nil
}

func (r *RPC) parseInt(data []byte, start int) (int, error) {

	num := string(data)
	if strings.Index(num, "\x00") >= 0 {
		num = num[:strings.Index(string(num), "\x00")]
	}

	i, err := strconv.Atoi(num)
	if err != nil {
		return 0, err
	}

	i += start
	return i, nil
}
