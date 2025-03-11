package operations

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type PollRequest struct {
	Type    string         `json:"type"`
	Offsets map[string]int `json:"offsets"`
}

type PollResponse struct {
	Type string              `json:"type"`
	Msgs map[string][][2]int `json:"msgs"`
}

func HandlePoll(n *maelstrom.Node, msg maelstrom.Message) error {
	var req PollRequest
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return err
	}

	responseMessages := make(map[string][][2]int)
	for key, offset := range req.Offsets {
		log, err := GlobalState.getLog(key, offset)
		if err != nil {
			continue
		}
		responseMessages[key] = append(responseMessages[key], [2]int{offset, log.Value})
	}

	resp := PollResponse{
		Msgs: responseMessages,
		Type: "poll_ok",
	}

	return n.Reply(msg, resp)
}
