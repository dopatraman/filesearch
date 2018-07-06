package message

import (
	"fmt"
	"reflect"
)

type PingMessage struct {
	Ping bool
}

func (m *PingMessage) MessageType() reflect.Type {
	return reflect.TypeOf(m)
}

func (m *PingMessage) String() string {
	return fmt.Sprintf("Ping: Content: %v", m.Ping)
}
