# CHIP-8 [![Build Status](https://travis-ci.org/ejholmes/chip8.svg?branch=master)](https://travis-ci.org/ejholmes/chip8)

[Godoc](https://godoc.org/github.com/ejholmes/chip8)

CHIP-8 emulator written in Go.

## Usage

This comes with a `chip8` package that can be used as a library for executing CHIP-8 binary programs, and also a `chip8` reference command.

You can install it with:

```console
go get github.com/ejholmes/chip8/cmd/chip8
```

And run a binary like so:

```console
chip8 run myprog.ch8
```

The default display implementation uses [go-termbox](https://github.com/nsf/termbox-go) so the program runs entirely inside your terminal.

## Reference

http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
