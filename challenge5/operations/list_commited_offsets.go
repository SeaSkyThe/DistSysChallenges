package operations

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type ListCommittedOffsetsRequest struct {
	Type string   `json:"type"`
	Keys []string `json:"keys"`
}

type ListCommittedOffsetsResponse struct {
	Type    string         `json:"type"`
	Offsets map[string]int `json:"offsets"`
}

func HandleListCommittedOffsets(n *maelstrom.Node, msg maelstrom.Message) error {
	var req ListCommittedOffsetsRequest

	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return err
	}

	responseBody := ListCommittedOffsetsResponse{
		Type:    "list_committed_offsets_ok",
		Offsets: map[string]int{},
	}
	return n.Reply(msg, responseBody)
}
