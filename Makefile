.PHONY: bin/chip8

bin/chip8:
	go build -o $@ ./cmd/chip8
