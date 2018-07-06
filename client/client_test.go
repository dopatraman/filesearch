package client

import (
	"bytes"
	"net/http"
	"testing"
)

func TestSendRequest(t *testing.T) {
	var c FileSearchClient
	c = FileSearchClient{}

	// Setup
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	go http.ListenAndServe(":8081", nil)

	t.Run("should return a value", func(t *testing.T) {
		ch := make(chan *bytes.Buffer)
		c.SendRequest("POST", "http://localhost:8081/", bytes.NewBuffer([]byte("Hello")), ch)
		v := <- ch
		if v.String() != "OK" {
			t.Error("No request sent")
		}
	})
	t.Run("should throw a fatal error", func(t *testing.T) {
		ch := make(chan *bytes.Buffer)
		c.SendRequest("POST", "http://localhost:8082/", bytes.NewBuffer([]byte("Hello")), ch)
		_, ok := <- ch
		if ok {
			t.Error("Client did not close response channel")
		}
	})
}