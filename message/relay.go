package message

import (
	"fmt"
	"reflect"
)

type RelayMessage struct {
	FileData string
	Address  string
}

func (m *RelayMessage) MessageType() reflect.Type {
	return reflect.TypeOf(m)
}

func (m *RelayMessage) String() string {
	return fmt.Sprintf("RelayMessage: FileData: %v, Address: %v", m.FileData, m.Address)
}
