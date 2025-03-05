# [Challenge \#2](https://fly.io/dist-sys/2/) - Unique ID Generation

In this challenge, you’ll need to implement a globally-unique ID generation system that runs against Maelstrom’s unique-ids workload. 
Your service should be totally available, meaning that it can continue to operate even in the face of network partitions.

## Run

You can run it by doing:

1. Install _maelstrom_ go library: `go get github.com/jepsen-io/maelstrom/demo/go`
2. Install this package in your machine: `go install .`
3. Run the `malestrom` util: `maelstrom test -w unique-ids --bin ~/go/bin/challenge2 --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition`.
