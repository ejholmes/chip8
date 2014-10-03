.PHONY: cmd test

test:
	go test

cmd:
	go build -o ./build/chip8 ./cmd/chip8
