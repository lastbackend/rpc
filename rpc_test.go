package rpc

import (
	"fmt"
	"testing"
)

func TestRegister(t *testing.T) {

	var (
		name  = "test"
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

	if r.name != name {
		t.Error("Expected name: %s got %s", name, r.name)
	}

	if r.uuid != uuid {
		t.Error("Expected uuid: %s got %s", uuid, r.uuid)
	}

	if r.token != token {
		t.Error("Expected token: %s got %s", token, r.token)
	}

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	r.SetURI(uri)
	if r.uri != uri {
		t.Error("Expected uri: %s got %s", uri, r.uri)
	}

	r.Listen()
	select {
	case <-r.connected:
		return
	}

}
