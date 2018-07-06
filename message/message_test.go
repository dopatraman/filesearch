package message

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
	t.Run("Should serialize a", func (t *testing.T) {
		t.Run("n unpackable message", func(t *testing.T) {
			var u *UnpackableMessage
			u = &UnpackableMessage{
				Data: []byte(`hello`),
			}
			packed := Pack(u)

			if !isJSONByteArray(packed.Data) {
				t.Error(fmt.Sprintf("Expected %v to be json", packed.Data))
			}
			testvalue, _ := json.Marshal(u)
			if (!reflect.DeepEqual(testvalue, packed.Data)) {
				t.Error(fmt.Sprintf("Expected %v to be %v", testvalue, packed.Data))
			}
		})
	})
}

func isJSONByteArray(b []byte) bool {
    var js interface{}
    return json.Unmarshal(b, &js) == nil
}