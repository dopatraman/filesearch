package message

import (
	"reflect"
)

type SearchMessage struct {
	FileHash      string
	SenderAddress string
}

func (s *SearchMessage) MessageType() reflect.Type {
	return reflect.TypeOf(s)
}

func (s *SearchMessage) String() string {
	return ""
}
