package server

import (
	"bytes"
	"net/http"
	"testing"
)

func TestListen(t *testing.T) {
	s := FileSearchServer{}
	t.Run("Should listen for requests", func(t *testing.T) {
		ch := s.Listen(8080)
		go request()
		v := <-ch
		if v.String() != "Hello" {
			t.Error("Request not received properly")
		}
	})
}

func request() {
	req, _ := http.NewRequest("POST", "http://localhost:8080", bytes.NewBuffer([]byte(`Hello`)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, _ = client.Do(req)
}
