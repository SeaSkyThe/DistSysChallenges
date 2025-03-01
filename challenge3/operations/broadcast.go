package operations

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func sendViaBroadcast(n *maelstrom.Node, target string, msgBody map[string]any) error {
	backoff := time.Millisecond * 50
	for i := 0; i < 5; i++ {
		if err := n.Send(target, msgBody); err == nil {
			break
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return fmt.Errorf("failed to send message to %s after retries", target)
}

func propagateBroadcast(n *maelstrom.Node, messageValue int, msgSource string) {
	for _, nodeID := range NEIGHBOR_NODES {
		if nodeID != msgSource {
			msg_id, err := generateRandomMessageID()
			if err != nil {
				log.Printf("Error generating random msg_id: %v", err)
			}

			msg_body := map[string]any{
				"type":    "broadcast",
				"message": messageValue,
				"msg_id":  msg_id,
			}

			go func(nodeID string) {
				if err := sendViaBroadcast(n, nodeID, msg_body); err != nil {
					log.Printf("Failed to send broadcast to %s: %v", nodeID, err)
				} else {
					log.Printf("Successfully sent broadcast to %s", nodeID)
				}
			}(nodeID)
		}
	}
}

func HandleBroadcast(n *maelstrom.Node, msg maelstrom.Message) error { // Unmarshall the message body as an loosely-typed map.
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

	if !wasSeen(messageValue) {
		storeMessage(messageValue)

		// WaitGroup to ensure all broadcasts to complete
		propagateBroadcast(n, messageValue, msg.Src)
	}

	response := map[string]any{
		"type":   "broadcast_ok",
		"msg_id": body["msg_id"],
	}

	return n.Reply(msg, response)
}
