package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type FileSearchClient struct {}

func (c FileSearchClient) SendRequest(method string, url string, payload io.Reader, responseChan chan *bytes.Buffer) {
	req, _ := http.NewRequest(method, url, payload)
	req.Header.Set("Content-Type", "application/json")
	go c.send(req, responseChan)
}

func (c FileSearchClient) send(req *http.Request, responseChan chan *bytes.Buffer) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		close(responseChan)
		log.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	responseChan <- bytes.NewBuffer(body)
	resp.Body.Close()
}