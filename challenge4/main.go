package main

import (
	"challenge4/operations"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	operations.KV = *maelstrom.NewSeqKV(n)

	n.Handle("add", func(msg maelstrom.Message) error {
		return operations.HandleAdd(n, msg)
	})
	n.Handle("read", func(msg maelstrom.Message) error {
		return operations.HandleRead(n, msg)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
