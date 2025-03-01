# Challenge \#3

## This challenge has five parts

A. [Single-Node Broadcast](https://fly.io/dist-sys/3a/)
B. [Multi-Node Broadcast](https://fly.io/dist-sys/3b/)
C. [Fault Tolerant Broadcast](https://fly.io/dist-sys/3c/)
D. [Efficient Broadcast, Part I](https://fly.io/dist-sys/3d/)

In this challenge, on the first part you’ll need to implement a globally-unique ID generation system that runs against Maelstrom’s unique-ids workload. Your service should be totally available, meaning that it can continue to operate even in the face of network partitions.

In the second part you will increment your single-node broadcast implementation to make it replicate our messages across a cluster that has no network partitions.

The third consists in making our multi-node broadcast implementation fault tolerant.

The fourth part is about improving our yet Fault-Tolerant, Multi-Node Broadcast implementation to be efficient. Now we are focusing on making our distributed system not only correct but also fast.
We have to achieve the following:

- Messages-per-operation: below `30`
- Median latency: below `400ms`
- Maximum latency: `600ms`

The fifth part making our broadcast even more efficient:

- Messages-per-operation: below `20`
- Median latency: below `1 second`
- Maximum latency: `2 seconds`

## Warning

The implementation here is the final result after doing the three parts.

## Run

You can run it by doing:

1. Install this package in your machine: `go install .`
2. Run the `malestrom` util to test each part:
   1. `maelstrom test -w broadcast --bin ~/go/bin/challenge3 --node-count 1 --time-limit 20 --rate 10`
   2. `maelstrom test -w broadcast --bin ~/go/bin/challenge3 --node-count 5 --time-limit 20 --rate 10`
   3. `maelstrom test -w broadcast --bin ~/go/bin/challenge3 --node-count 5 --time-limit 20 --rate 10 --nemesis partition`
   4. `maelstrom test -w broadcast --bin ~/go/bin/challenge3 --node-count 25 --time-limit 20 --rate 100 --latency 100`
   5. `maelstrom test -w broadcast --bin ~/go/bin/challenge3 --node-count 25 --time-limit 20 --rate 100 --latency 100`

Or, if you are familiar with Makefiles you can run:
1. `make a` for the part A test.
2. `make b` for the part B test.
3. `make c` for the part C test (PS: you have to uncomment the `replication` handler in the `main.go` file)
4. `make d` for the part D test.
5. `make e` for the part E test.

