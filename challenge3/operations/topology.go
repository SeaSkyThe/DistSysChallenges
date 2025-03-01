package operations

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func HandleTopology(n *maelstrom.Node, msg maelstrom.Message) error {
	// Unmarshall the message body as an loosely-typed map.
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	response := map[string]any{
		"type":   "topology_ok",
		"msg_id": body["msg_id"],
	}

	topology, ok := body["topology"].(map[string]any)

	if !ok {
		return fmt.Errorf("topology type assertion failed")
	}

	if neighbors, exists := topology[n.ID()]; exists {
		neighborsSlice, ok := neighbors.([]any)
		if !ok {
			return fmt.Errorf("neighbors type assertion failed")
		}
		temp := make([]string, len(neighborsSlice))
		for i, v := range neighborsSlice {
			id, ok := v.(string)
			if !ok {
				return fmt.Errorf("failed to convert neighbor to string")
			}
			temp[i] = id
		}

		NEIGHBOR_NODES = temp
	}

	return n.Reply(msg, response)
}
