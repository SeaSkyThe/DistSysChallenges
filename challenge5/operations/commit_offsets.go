package operations

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type CommitOffsetRequest struct {
	Type    string         `json:"type"`
	Offsets map[string]int `json:"offsets"`
}

type CommitOffsetResponse struct {
	Type string `json:"type"`
}

func HandleCommitOffsets(n *maelstrom.Node, msg maelstrom.Message) error {
	var req CommitOffsetRequest

	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return err
	}

	for key, offset := range req.Offsets {
		GlobalState.commitOffset(key, offset)
	}

	response_body := CommitOffsetResponse{
		Type: "commit_offsets_ok",
	}

	return n.Reply(msg, response_body)
}
