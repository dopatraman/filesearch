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
	mux sync.RWMutex
	client.FileSearchClient
	server.FileSearchServer
	neighboringNodes []string
	searchCache      map[string]chan interface{}
	senderCache      map[string]chan *bytes.Buffer
	responseCache    map[string]chan *bytes.Buffer
	listener         chan *bytes.Buffer
}

func (n *FileSearchNode) Connect(addr string) {
	// eventually some authentication should happen here
	// this method can only be called once per address...
	// @TODO: make sure this method is only ever called once
	// @TODO: senderCache should be a map of address strings : tuple(sender, receiver)
	//			will cut down on multiple attempts to access
	senderChan, responseChan := n.createSenderAndResponseChannels(addr)
	go n.connect(addr, senderChan, responseChan)

	var u *message.UnpackableMessage
	confirmation := <-responseChan
	err := json.Unmarshal(confirmation.Bytes(), u)
	if err != nil {
		panic("Unexpected response")
	}
	n.handleMessage(u.Unpack())
}

func (n *FileSearchNode) SendMessage(addr string, msg message.NodeMessage) *bytes.Buffer {
	var u *message.UnpackableMessage
	u = message.Pack(msg)

	// @TODO: see above, sender cache should be a map of address strings to tuple(sender, receiver)
	senderChan, responseChan := n.findSenderAndResponseChannels(addr)
	go func(ch chan *bytes.Buffer) { ch <- bytes.NewBuffer(u.Data) }(senderChan)
	return <-responseChan
}

func (n *FileSearchNode) HandleMessages() {
	// @TODO: check to see if listener already exists
	n.listener = n.Listen(8080)
	for rawBuffer := range n.listener {
		var u *message.UnpackableMessage
		err := json.Unmarshal(rawBuffer.Bytes(), &u)
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
	case *message.ConfirmationMessage:
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

func (n *FileSearchNode) handleConfirmation(msg *message.ConfirmationMessage) {
	// @TODO: some validation should happen here
	if msg.Sender != "" {
		newNodes := append(n.neighboringNodes, msg.Sender)
		n.neighboringNodes = newNodes
	}
}

func (n *FileSearchNode) createSenderAndResponseChannels(addr string) (chan *bytes.Buffer, chan *bytes.Buffer) {
	senderChan, responseChan := n.findSenderAndResponseChannels(addr)
	if senderChan == nil {
		senderChan = make(chan *bytes.Buffer)
		n.mux.Lock()
		n.senderCache[addr] = senderChan
		n.mux.Unlock()
	}
	if responseChan == nil {
		responseChan = make(chan *bytes.Buffer)
		n.mux.Lock()
		n.responseCache[addr] = responseChan
		n.mux.Unlock()
	}
	return senderChan, responseChan
}

func (n *FileSearchNode) findSenderAndResponseChannels(addr string) (chan *bytes.Buffer, chan *bytes.Buffer) {
	n.mux.RLock()
	senderChan := n.senderCache[addr]
	n.mux.RUnlock()

	n.mux.RLock()
	responseChan := n.responseCache[addr]
	n.mux.RUnlock()
	return senderChan, responseChan
}
