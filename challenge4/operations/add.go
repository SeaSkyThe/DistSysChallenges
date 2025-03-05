package operations

import (
	"context"
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func HandleAdd(n *maelstrom.Node, msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	d, ok := body["delta"].(float64)

	if !ok {
		return fmt.Errorf("add: invalid delta")
	}

	delta := int(d)

	ctx := context.Background()

	currentValue := 0
	value, err := KV.ReadInt(ctx, n.ID())

	if err != nil {
		if rpcErr, ok := err.(*maelstrom.RPCError); !ok || rpcErr.Code != maelstrom.KeyDoesNotExist {
			return fmt.Errorf("add: failed to read from KV store: %w", err)
		}
	} else {
		currentValue = value
	}

	newValue := currentValue + delta

	KV.CompareAndSwap(ctx, n.ID(), currentValue, newValue, true)

	response_body := map[string]any{
		"type": "add_ok",
	}

	return n.Reply(msg, response_body)
}
