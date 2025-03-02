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

	n.Handle("sync", func(msg maelstrom.Message) error {
		return operations.HandleSynchronization(n, msg)
	})

	n.Handle("broadcast_ok", func(msg maelstrom.Message) error {
		return nil
	})

	n.Handle("read_ok", func(msg maelstrom.Message) error {
		return nil
	})

	n.Handle("topology_ok", func(msg maelstrom.Message) error {
		return nil
	})

	n.Handle("sync_ok", func(msg maelstrom.Message) error {
		return operations.HandleSyncOk(n, msg)
	})

	go operations.SyncData(n, 200)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
