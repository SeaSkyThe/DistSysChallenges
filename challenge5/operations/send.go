package operations

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type SendRequest struct {
	Type string `json:"type"`
	Key  string `json:"key"`
	Msg  int    `json:"msg"`
}

type SendResponse struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
}

func HandleSend(n *maelstrom.Node, msg maelstrom.Message) error {
	var req SendRequest

	if err := json.Unmarshal(msg.Body, &req); err != nil {
		return err
	}

	log, err := GlobalState.addLog(req.Key, req.Msg)
	if err != nil {
		return err
	}

	responseBody := SendResponse{
		Type:   "send_ok",
		Offset: log.Offset,
	}

	return n.Reply(msg, responseBody)
}
