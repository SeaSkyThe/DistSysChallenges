package main

import (
	"challenge3/operations"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		return operations.HandleBroadcast(n, msg)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		return operations.HandleRead(n, msg)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		return operations.HandleTopology(n, msg)
	})

    // FOR TEST C TO WORK, UNCOMMENT THIS BELOW
	n.Handle("replicate", func(msg maelstrom.Message) error {
		return operations.HandleReplicate(n, msg)
	})

	go operations.ReplicateData(n, 1000)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
