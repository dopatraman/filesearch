package message

import (
	"encoding/json"
	"log"
	"reflect"
)

type NodeMessage interface {
	MessageType() reflect.Type
}

func Pack(msg NodeMessage) *UnpackableMessage {
	var b []byte
	var err error
	b, err = json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	u := &UnpackableMessage{
		Data: b,
	}
	return u
}
