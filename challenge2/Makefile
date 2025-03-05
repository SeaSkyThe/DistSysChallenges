.PHONY: install test run

install:
	go install .

test: install
	maelstrom test -w unique-ids --bin ~/go/bin/challenge2 $(ARGS)

run:
	$(MAKE) test ARGS="--node-count 3 --time-limit 30 --rate 100 --availability total --nemesis partition"
