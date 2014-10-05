package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ejholmes/chip8"
)

var (
	Options = chip8.DefaultOptions
	Display = chip8.DefaultDisplay
	Keypad  = chip8.DefaultKeypad
)

func main() {
	var (
		program = flag.String("program", "", "Path to the program to run.")
		clock   = flag.Int("clock", int(chip8.DefaultClockSpeed), "Clock speed in Hz.")
	)
	flag.Parse()

	f, err := os.Create("log")
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(f, "", 0)
	chip8.DefaultLogger = logger

	Options.ClockSpeed = time.Duration(*clock)
	Display.Init()
	defer Display.Close()

	Keypad.Init()

	c, err := chip8.NewCPU(nil)
	if err != nil {
		logger.Println(err)
		os.Exit(-2)
	}

	if *program == "" {
		flag.Usage()
		os.Exit(-2)
	}

	raw, err := ioutil.ReadFile(*program)
	if err != nil {
		logger.Println(err)
		os.Exit(-2)
	}

	// Load the program into RAM.
	c.LoadBytes(raw)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		Display.Close()
		os.Exit(-2)
	}()

	// Run it.
	if err := c.Run(); err != nil {
		logger.Println(err)
	}
}
