package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type FileSearchServer struct{}

func (s FileSearchServer) Listen(port int) chan *bytes.Buffer {
	ch := make(chan *bytes.Buffer)
	http.HandleFunc("/", handleRoot(ch))
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return ch
}

func handleRoot(ch chan *bytes.Buffer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			ch <- bytes.NewBuffer([]byte("Error!"))
			log.Fatal("Error!")
		}
		w.Write(b)
		ch <- bytes.NewBuffer(b)
	}
}

func handleConnect(ch chan *bytes.Buffer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
