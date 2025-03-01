package operations

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func HandleRead(n *maelstrom.Node, msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	messagesCopy := readMessages()

	response_body := map[string]any{
		"type":     "read_ok",
		"messages": messagesCopy,
		"msg_id":   body["msg_id"],
	}

	return n.Reply(msg, response_body)
}
