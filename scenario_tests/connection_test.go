package scenario_tests

import (
	"testing"

	"github.com/user/file_search/client"
	"github.com/user/file_search/node"
	"github.com/user/file_search/server"
)

func TestConnectionToNode(t *testing.T) {
	var actorNode, targetNode *node.FileSearchNode
	actorNode = createNode()
	targetNode = createNode()
	_ = targetNode.Listen(8080)
	t.Run("Connect to a node", func(t *testing.T) {
		addr := "http://localhost:8080"
		err := actorNode.Connect(addr)
		if err != nil {
			t.Error("Could not connect to target node")
		}
	})
}

func createNode() *node.FileSearchNode {
	c := client.FileSearchClient{}
	s := server.FileSearchServer{}
	return &node.FileSearchNode{
		FileSearchClient: c,
		FileSearchServer: s,
	}
}
