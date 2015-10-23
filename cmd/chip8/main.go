package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "chip8"
	app.Usage = "Run chip8 programs using a Go based emulator"
	app.Commands = []cli.Command{
		cmdRun,
	}
	app.Run(os.Args)
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
}
