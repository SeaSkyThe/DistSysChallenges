package main

import (
	"challenge5/operations"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := operations.GlobalState.Node

	n.Handle("send", func(msg maelstrom.Message) error {
		return operations.HandleSend(n, msg)
	})
	n.Handle("poll", func(msg maelstrom.Message) error {
		return operations.HandlePoll(n, msg)
	})
	n.Handle("commit_offsets", func(msg maelstrom.Message) error {
		return operations.HandleCommitOffsets(n, msg)
	})
	n.Handle("list_committed_offsets", func(msg maelstrom.Message) error {
		return operations.HandleListCommittedOffsets(n, msg)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
