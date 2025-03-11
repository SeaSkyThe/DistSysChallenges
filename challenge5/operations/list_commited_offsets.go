package operations

import (
	"encoding/json"
	"fmt"

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

	offsets := make(map[string]int)
	for _, key := range req.Keys {
		offset, ok := GlobalState.getHighestCommitedOffset(key)
		if !ok {
			return fmt.Errorf("error when getting highest offset for key: %s", key)
		}
		offsets[key] = offset
	}

	responseBody := ListCommittedOffsetsResponse{
		Type:    "list_committed_offsets_ok",
		Offsets: offsets,
	}

	return n.Reply(msg, responseBody)
}
