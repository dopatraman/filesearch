package message

import (
	"encoding/json"
	"testing"
)

func TestUnpack(t *testing.T) {
	t.Run("Should unpack a", func(t *testing.T) {
		var u *UnpackableMessage
		t.Run("Search message", func(t *testing.T) {
			var m SearchMessage
			m = SearchMessage{
				FileHash:      "wefnwlefnwe",
				SenderAddress: "192.168.2.1",
			}
			mJson, _ := json.Marshal(m)
			u = &UnpackableMessage{
				Data: mJson,
			}
			var s SearchMessage
			json.Unmarshal(mJson, &s)
			newM := u.Unpack()
			_, ok := newM.(*SearchMessage)
			if !ok {
				t.Error("no")
			}
		})
	})
}
