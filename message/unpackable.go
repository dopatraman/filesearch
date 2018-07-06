package message

import (
	"encoding/json"
	"log"
	"reflect"
)

type UnpackableMessage struct {
	Data []byte
}

func (u *UnpackableMessage) MessageType() reflect.Type {
	return reflect.TypeOf(u)
}

func (u *UnpackableMessage) Unpack() NodeMessage {
	if len(u.Data) == 0 {
		log.Fatal("empty data field")
		return nil
	}
	// check search
	var s *SearchMessage
	err := json.Unmarshal(u.Data, s)
	if err == nil {
		return s
	}
	// check relay
	var r *RelayMessage
	err = json.Unmarshal(u.Data, r)
	if err == nil {
		return r
	}
	// check confirmation
	var c *ConfirmationMessage
	err = json.Unmarshal(u.Data, c)
	if err == nil {
		return c
	}
	// check ping
	var v *PingMessage
	err = json.Unmarshal(u.Data, v)
	if err == nil {
		return v
	}
	panic("Unhandled message")
}

func (u *UnpackableMessage) String() string {
	return string(u.Data)
}
