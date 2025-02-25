# [Challenge \#1](https://fly.io/dist-sys/1/)

This challenge is more of a "getting started guide".

# Run

You can run it by doing:

1. Install _maelstrom_ go library: `go get github.com/jepsen-io/maelstrom/demo/go`
2. Install this package in your machine: `go install .`
3. Run the `malestrom` util: `maelstrom test -w echo --bin ~/go/bin/challenge1 --node-count 1 --time-limit 10`
