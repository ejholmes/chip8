.PHONY: cmd test

test:
	godep go test

cmd:
	godep go build -o ./build/chip8 ./cmd/chip8
