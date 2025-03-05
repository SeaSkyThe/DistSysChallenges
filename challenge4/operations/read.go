package operations

import (
	"context"
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func HandleRead(n *maelstrom.Node, msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	ctx := context.Background()
	// Get from all nodes and sum
    totalValue := 0

    for _, nodeID := range n.NodeIDs() {
        value, err := KV.ReadInt(ctx, nodeID)
        if err != nil {
            // If key doesn't exist, return 0
            if _, ok := err.(*maelstrom.RPCError); ok && err.(*maelstrom.RPCError).Code == maelstrom.KeyDoesNotExist {
                value = 0
            } else {
                return err
            }
        }
        totalValue += value

    }

	response_body := map[string]any{
		"type":  "read_ok",
		"value": totalValue,
	}
	return n.Reply(msg, response_body)
}
