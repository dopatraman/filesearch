package message

import (
	"fmt"
	"reflect"
)

type StopMessage struct{}

func (m *StopMessage) MessageType() reflect.Type {
	return reflect.TypeOf(m)
}

func (m *StopMessage) String() string {
	return fmt.Sprintf("StopMessage")
}
