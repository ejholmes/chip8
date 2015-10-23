package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/ejholmes/chip8"
)

var cmdRun = cli.Command{
	Name:   "run",
	Usage:  "Run a chip8 program",
	Action: runRun,
}

func runRun(c *cli.Context) {
	if !c.Args().Present() {
		cli.ShowAppHelp(c)
		return
	}

	// Read program.
	program, err := ioutil.ReadFile(c.Args().First())
	must(err)

	// Initialize termbox.
	err = chip8.TermboxInit()
	defer chip8.TermboxCleanup()
	must(err)

	// Initialize peripherals.
	d := chip8.NewTermboxDisplay()
	k := chip8.NewTermboxKeypad()

	// Initialize CPU.
	cpu, err := chip8.NewCPU(nil)
	must(err)
	cpu.Graphics.Display = d
	cpu.Keypad = k

	// Load program.
	_, err = cpu.LoadBytes(program)
	must(err)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		cpu.Stop()
	}()

	// Run it.
	err = cpu.Run()
	must(err)
}
