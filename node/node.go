package node

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/user/file_search/client"
	"github.com/user/file_search/message"
	"github.com/user/file_search/server"
)

type FileSearchNode struct {
	mux sync.Mutex
	client.FileSearchClient
	server.FileSearchServer
	neighbors     []string
	searchCache   map[string]chan interface{}
	senderCache   map[string]chan *bytes.Buffer
	responseCache map[string]chan *bytes.Buffer
	listener      chan *bytes.Buffer
}

func (n *FileSearchNode) Connect(addr string) error {
	// eventually some authentication should happen here
	// this method can only be called once per address...
	// @TODO: make sure this method is only ever called once
	// @TODO: senderCache should be a map of address strings : tuple(sender, receiver)
	//			will cut down on multiple attempts to access

	// lock + release mutex?
	n.mux.Lock()
	senderChan := n.senderCache[addr]
	n.mux.Unlock()
	if senderChan == nil {
		senderChan = make(chan *bytes.Buffer)
	}
	n.mux.Lock()
	responseChan := n.responseCache[addr]
	n.mux.Unlock()
	if responseChan == nil {
		responseChan = make(chan *bytes.Buffer)
	}
	go n.connect(addr, senderChan, responseChan)

	n.mux.Lock()
	n.senderCache[addr] = senderChan
	n.mux.Unlock()

	n.mux.Lock()
	n.responseCache[addr] = responseChan
	n.mux.Unlock()
	return nil
}

func (n *FileSearchNode) SendMessage(addr string, msg message.NodeMessage) *bytes.Buffer {
	var u *message.UnpackableMessage
	u = message.Pack(msg)

	// @TODO: see above, sender cache should be a map of address strings to tuple(sender, receiver)
	n.mux.Lock()
	senderChan := n.senderCache[addr]
	n.mux.Unlock()

	n.mux.Lock()
	responseChan := n.responseCache[addr]
	n.mux.Unlock()

	go func(ch chan *bytes.Buffer) { ch <- bytes.NewBuffer(u.Data) }(senderChan)
	return <-responseChan
}

func (n *FileSearchNode) HandleMessages() {
	// @TODO: check to see if listener already exists
	n.listener = n.Listen(8080)
	for rawBuffer := range n.listener {
		var u message.UnpackableMessage
		var err error
		err = json.Unmarshal(rawBuffer.Bytes(), &u)
		if err != nil {
			panic("Unrecognized message format")
		}
		n.handleMessage(u.Unpack())
	}
}

func (n *FileSearchNode) handleMessage(msg message.NodeMessage) {
	// this will likely be a switch-case of all messages and how to handle them
	switch msg.(type) {
	case *message.SearchMessage:
		searchMessage, ok := msg.(*message.SearchMessage)
		if !ok {
			panic("Could not cast NodeMessage to SearchMessage")
		}
		n.handleSearch(searchMessage)
	case *message.RelayMessage:
		return
	case *message.PingMessage:
		return
	default:
		return
	}
}

func (n *FileSearchNode) connect(addr string, s chan *bytes.Buffer, r chan *bytes.Buffer) {
	for {
		msg, ok := <-s
		if !ok {
			break
		}
		n.SendRequest("POST", addr, msg, r)
	}
}

func (n *FileSearchNode) handleSearch(msg *message.SearchMessage) {
	filepath.Walk("~/", func(path string, info os.FileInfo, err error) error {
		// if file found, encrypt
		file := ""
		n.handleMessage(&message.RelayMessage{
			FileData: file,
			Address:  msg.SenderAddress,
		})
		return nil
	})
}
