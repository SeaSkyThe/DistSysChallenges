package operations

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func HandleBroadcast(n *maelstrom.Node, msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Extract the value
	value := body["message"]
	v, ok := value.(float64)
	if !ok {
		return fmt.Errorf("broadcast: value type assertion failed")
	}

	messageValue := int(v)

	storeMessage(messageValue)

	response := map[string]any{
		"type":   "broadcast_ok",
		"msg_id": body["msg_id"],
	}

	return n.Reply(msg, response)
}
