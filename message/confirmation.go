package message

import (
	"fmt"
	"reflect"
)

type ConfirmationMessage struct {
	Sender string
}

func (m *ConfirmationMessage) MessageType() reflect.Type {
	return reflect.TypeOf(m)
}

func (m *ConfirmationMessage) String() string {
	return fmt.Sprintf("Confirmation Message: {Sender: %v", m.Sender)
}
