package chip8_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ejholmes/chip8"
	termbox "github.com/nsf/termbox-go"
)

var programs = []string{
	"pong.ch8",
	"invaders.ch8",
}

func Test(t *testing.T) {
	for _, program := range programs {
		t.Run(program, runProgram(program))
	}
}

func runProgram(program string) func(t *testing.T) {
	return func(t *testing.T) {
		// Initialize peripherals.
		d, err := chip8.NewTermboxDisplay(
			termbox.ColorDefault, // Foreground
			termbox.ColorDefault, // Background
		)
		if err != nil {
			t.Fatal(err)
		}
		defer d.Close()

		k := chip8.NewTermboxKeypad()

		cpu, err := chip8.NewCPU(nil)
		if err != nil {
			t.Fatal(err)
		}

		cpu.Graphics.Display = d
		cpu.Keypad = k

		r, err := os.Open(filepath.Join("../programs", program))
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()

		if _, err := cpu.Load(r); err != nil {
			t.Fatal(err)
		}

		if err := cpu.Run(); err != nil {
			t.Fatal(err)
		}
	}
}
