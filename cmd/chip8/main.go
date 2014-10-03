package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ejholmes/chip8"
)

func main() {
	var (
		program = flag.String("program", "", "Path to the program to run.")
		clock   = flag.Int("clock", int(chip8.DefaultClockSpeed), "Clock speed in Hz.")
	)
	flag.Parse()

	chip8.DefaultOptions.ClockSpeed = time.Duration(*clock)

	c, err := chip8.NewCPU(nil)
	if err != nil {
		log.Fatal(err)
	}

	if *program == "" {
		flag.Usage()
		os.Exit(-2)
	}

	raw, err := ioutil.ReadFile(*program)
	if err != nil {
		log.Fatal(err)
	}

	c.LoadBytes(raw)
	log.Fatal(c.Run())
}
