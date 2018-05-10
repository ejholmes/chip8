package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
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

func printErr(err error) {
	fmt.Fprintf(os.Stderr, "error: %s\n", err)
}
