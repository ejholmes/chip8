package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ejholmes/chip8"
	"github.com/nsf/termbox-go"
	"github.com/urfave/cli"
)

var cmdRun = cli.Command{
	Name:   "run",
	Usage:  "Run a chip8 program",
	Action: withErrors(runRun),
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "log",
			Usage: "If provided, specifies a log file to write debug output to.",
		},
		cli.IntFlag{
			Name:  "clock",
			Usage: "Clock speed, in hz, to run at.",
			Value: int(chip8.DefaultClockSpeed),
		},
	},
}

func runRun(c *cli.Context) error {
	// Initialize peripherals.
	d, err := chip8.NewTermboxDisplay(
		termbox.ColorDefault, // Foreground
		termbox.ColorDefault, // Background
	)
	defer d.Close()
	if err != nil {
		return err
	}

	k := chip8.NewTermboxKeypad()

	// Initialize CPU.
	cpu, err := chip8.NewCPU(&chip8.Options{
		ClockSpeed: time.Duration(c.Int("clock")),
	})
	if err != nil {
		return err
	}
	cpu.Graphics.Display = d
	cpu.Keypad = k

	// If a log file is specified, create a logger and add it to the CPU.
	if fname := c.String("log"); fname != "" {
		f, err := os.Create(fname)
		if err != nil {
			return err
		}

		cpu.Logger = log.New(f, "", 0)
	}

	if c.Args().Present() {
		// Read program.
		program, err := ioutil.ReadFile(c.Args().First())
		if err != nil {
			return err
		}

		// Load program.
		_, err = cpu.LoadBytes(program)
		if err != nil {
			return err
		}
	} else {
		_, err = cpu.Load(os.Stdin)
		if err != nil {
			return err
		}
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		cpu.Stop()
	}()

	// Run it.
	err = cpu.Run()
	return err
}
