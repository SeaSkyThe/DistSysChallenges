.PHONY: install test run

install:
	go install .

test: install
	maelstrom test -w g-counter --bin ~/go/bin/challenge4 $(ARGS)

run:
	$(MAKE) test ARGS="--node-count 3 --time-limit 20 --rate 100 --nemesis partition"
