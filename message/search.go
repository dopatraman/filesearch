package message

import (
	"fmt"
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
	return fmt.Sprintf("{FileHash: %v, SenderAddress: %v}", s.FileHash, s.SenderAddress)
}
