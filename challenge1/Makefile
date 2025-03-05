.PHONY: install test run

install:
	go install .

test: install
	maelstrom test -w echo --bin ~/go/bin/challenge1 $(ARGS)

run:
	$(MAKE) test ARGS="--node-count 1 --time-limit 10"
