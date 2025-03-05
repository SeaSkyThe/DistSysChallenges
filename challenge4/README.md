# Challenge \#4 - [Grow-Only Counter](https://fly.io/dist-sys/4/)

In this challenge, you’ll need to implement a stateless, grow-only counter which will run against Maelstrom’s g-counter workload.
This challenge is different than before in that your nodes will rely on a sequentially-consistent key/value store service provided by Maelstrom.

## Solution

The solution was straightforward. For add messages, we increment the value of a node-specific key in the KV database by the given delta. For read messages, we simply return the sum of all values across all keys in the KV database.

## Run

1. Install this package in your machine: `go install .`
2. Run the `malestrom` test: `maelstrom test -w g-counter --bin ~/go/bin/challenge4 --node-count 3 --rate 100 --time-limit 20 --nemesis partition`

Or, you prefer using Makefiles: `make run`
