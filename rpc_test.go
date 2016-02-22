package rpc

import (
	"testing"
	"fmt"
)

const (
	name  = "test"
	uuid  = "uuid"
	token = "token"
)

func TestRegister(t *testing.T) {

	var err error

	_, err = Register(name, uuid, token)

	if err != nil {
		t.Error("Register APP error", err)
	}

	if rpc.name != name {
		t.Error("Expected name: %s got %s", name, rpc.name)
	}

	if rpc.uuid != uuid {
		t.Error("Expected uuid: %s got %s", uuid, rpc.uuid)
	}

	if rpc.token != token {
		t.Error("Expected token: %s got %s", token, rpc.token)
	}

	uri := fmt.Sprintf("amqp://guest:guest@localhost:5672")
	rpc = SetURI(uri)
	if rpc.uri != uri {
		t.Error("Expected uri: %s got %s", uri, rpc.uri)
	}

	Listen()

	<- rpc.connected
}