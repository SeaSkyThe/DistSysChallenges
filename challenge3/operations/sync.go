package operations

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// CRDTS = Conflict-Free replicated data type

func SyncData(n *maelstrom.Node, interval int) {
	time_interval := time.Duration(interval) * time.Millisecond
	for {
		if len(NEIGHBORHOOD) == 0 {
			log.Printf("no connected nodes, waiting for topology...")
		}
		sendSynchronization(n)
		time.Sleep(time_interval)
	}
}

func sendSynchronization(n *maelstrom.Node) error {
	messagesCopy := readMessages()

	for _, nodeID := range NEIGHBORHOOD {
		msgId, err := generateRandomMessageID()
		if err != nil {
			return fmt.Errorf("error generating random msg_id for replicate event: %v", err)
		}

		body := map[string]any{
			"type":     "sync",
			"msg_id":   msgId,
			"messages": messagesCopy,
		}

		if err := n.Send(nodeID, body); err != nil {
			log.Printf("Failed to send replication to %s: %v", nodeID, err)
		} else {
			log.Printf("Successfully sent replication to %s", nodeID)
		}
	}

	return nil
}

// Related to the message of type "sync"
func HandleSynchronization(n *maelstrom.Node, msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	rawMessages, ok := body["messages"].([]any)
	if !ok {
		return fmt.Errorf("sync: messages type assertion failed")
	}

	messages := []int{}
	for _, v := range rawMessages {
		num, ok := v.(float64) // JSON numbers default to float64
		if !ok {
			return fmt.Errorf("sync: message value type assertion failed")
		}

		msgValue := int(num)
		messages = append(messages, msgValue)

	}
	storeMessages(messages)
	unknownMessages := neighborAck(messages, msg.Src)

	response_body := map[string]any{
		"type":     "sync_ok",
		"msg_id":   body["msg_id"],
		"messages": unknownMessages,
	}

	return n.Reply(msg, response_body)
}

func HandleSyncOk(n *maelstrom.Node, msg maelstrom.Message) error {
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	rawMessages, ok := body["messages"].([]interface{})
	if !ok {
		return fmt.Errorf("expected messages to be a []interface{}, but got: %T", body["messages"])
	}

	messages := []int{}
	for _, v := range rawMessages {
		// Type assert each element in the slice to float64
		if msg, ok := v.(int); ok {
			messages = append(messages, msg)
		} else {
			return fmt.Errorf("expected message element to be float64, but got: %T", v)
		}
	}

	storeMessages(messages)
	neighborAck(messages, msg.Src)

	return nil
}
