.PHONY: cmd

cmd:
	godep go build -o ./build/chip8 ./cmd/chip8
