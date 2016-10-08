.PHONY: cmd

bin/chip8:
	go build -o $@ ./cmd/chip8
