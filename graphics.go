// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chip8

import "github.com/nsf/termbox-go"

const (
	GraphicsWidth  = 64 // Pixels
	GraphicsHeight = 32 // Pixels
)

var (
	DefaultDisplay = Display(&display{})
)

// Display represents the output display for the CHIP-8 graphics array.
type Display interface {
	// Turn on the display and do any initialization.
	Init() error

	// Turn off the display and cleanup.
	Close()

	// Render should render the current graphics array to the display.
	Render(*Graphics) error
}

// Graphics represents the graphics array for the CHIP-8.
type Graphics struct {
	// The raw pixels of the graphics array.
	Pixels [GraphicsWidth * GraphicsHeight]byte

	// The display to render to.
	Display
}

// Draw draws the graphics array to the Display.
func (g *Graphics) Draw() error {
	return g.display().Render(g)
}

func (g *Graphics) display() Display {
	if g.Display == nil {
		return DefaultDisplay
	}

	return g.Display
}

// display is an implementation of the Display interface that renders
// the graphics array to the terminal.
type display struct{}

func (d *display) Init() error {
	if err := termbox.Init(); err != nil {
		return err
	}

	termbox.HideCursor()

	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorWhite); err != nil {
		return err
	}

	return termbox.Flush()
}

func (d *display) Close() {
	termbox.Close()
}

func (d *display) Render(g *Graphics) error {
	for x := 0; x < GraphicsWidth; x++ {
		for y := 0; y < GraphicsHeight; y++ {
			v := ' '

			if g.Pixels[x*y] == 0x01 {
				v = '*'
			}

			termbox.SetCell(
				x,
				y,
				v,
				termbox.ColorBlack,
				termbox.ColorWhite,
			)
		}
	}

	return termbox.Flush()
}
