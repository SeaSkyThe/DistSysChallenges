package operations

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// CRDTS = Conflict-Free replicated data type

func sendReplications(n *maelstrom.Node) error {
	messagesCopy := readMessages()

	var wg sync.WaitGroup
	for _, nodeID := range NEIGHBOR_NODES {
		msg_id, err := generateRandomMessageID()
		if err != nil {
			return fmt.Errorf("error generating random msg_id for replicate event: %v", err)
		}
		go func(nodeID string) {
			body := map[string]any{
				"type":     "replicate",
				"msg_id":   msg_id,
				"messages": messagesCopy,
			}
			wg.Add(1)

			defer wg.Done()
			if err := n.Send(nodeID, body); err != nil {
				log.Printf("Failed to send replication to %s: %v", nodeID, err)
			} else {
				log.Printf("Successfully sent replication to %s", nodeID)
			}
		}(nodeID)
	}

	wg.Wait()
	return nil
}

func ReplicateData(n *maelstrom.Node, interval int) {
	go func() {
		for i := 0; i < 50; i++ {
			time_interval := time.Duration(interval) * time.Millisecond
			if len(NEIGHBOR_NODES) == 0 {
				log.Printf("no connected nodes, waiting for topology...")
			}
			if err := sendReplications(n); err != nil {
				log.Printf("problem sending replications")
			}
			time.Sleep(time_interval)
			interval = int(interval / 2)
		}
	}()
}

func HandleReplicate(n *maelstrom.Node, msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	rawMessages, ok := body["messages"].([]any)
	if !ok {
		return fmt.Errorf("replicate: messages type assertion failed")
	}

	messages := []int{}
	for _, v := range rawMessages {
		num, ok := v.(float64) // JSON numbers default to float64
		if !ok {
			return fmt.Errorf("replicate: message value type assertion failed")
		}

		msgValue := int(num)
		messages = append(messages, msgValue)

		storeMessages(messages)
	}

	response_body := map[string]any{
		"type":   "replicate_ok",
		"msg_id": body["msg_id"],
	}

	return n.Reply(msg, response_body)
}
