.PHONY: install test a b c d

install:
	go install .

test: install
	maelstrom test -w kafka --bin ~/go/bin/challenge5 $(ARGS)

a: 
	$(MAKE) test ARGS="--node-count 1 --concurrency 2n --time-limit 20 --rate 1000"

b: 
	$(MAKE) test ARGS="--node-count 2 --concurrency 2n --time-limit 20 --rate 1000"
#
# c: 
# 	$(MAKE) test ARGS="--node-count 5 --time-limit 20 --rate 10 --nemesis partition"
#
# d: 
# 	$(MAKE) test ARGS="--node-count 25 --time-limit 20 --rate 100 --latency 100"
#
# e: 
# 	$(MAKE) test ARGS="--node-count 25 --time-limit 20 --rate 100 --latency 100"
