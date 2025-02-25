package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func randUUID() string {
	// Im trusting rand.Int :X

	possibleChars := [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	// Size 36 because: 8 dash 4 dash 4 dash 4 dash 12 (397a5719-763e-4402-aa58-c2753f80cdbd)
	randomID := make([]byte, 36)
	for i := 0; i < 36; i++ {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			randomID[i] = '-'
		} else {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(possibleChars))))
			randomID[i] = possibleChars[n.Int64()]
		}
	}

	return string(randomID[:])
}

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		// Unmarshall the message body as an loosely-typed map.
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "generate_ok"
		body["id"] = randUUID()

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
