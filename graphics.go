// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chip8

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

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

	// The display to render to. The nil value is the DefaultDisplay.
	Display
}

// DrawSprite draws a sprite to the graphics array starting at coording x, y.
// If there is a collision, WriteSprite returns true.
func (g *Graphics) WriteSprite(sprite []byte, x, y byte) (collision bool) {
	n := len(sprite)

	for yl := 0; yl < n; yl++ {
		// A row of sprite data.
		r := sprite[yl]

		for xl := 0; xl < 8; xl++ {
			// This represents a mask for the bit that we
			// care about for this coordinate.
			i := 0x80 >> byte(xl)

			var v byte

			// Whether the bit should be set or not
			if (r & byte(i)) == byte(i) {
				v = 0x01
			}

			// The address for this bit of data on the
			// graphics array.
			a := uint16(x) + uint16(xl) + ((uint16(y) + uint16(yl)) * GraphicsWidth)

			// If there's a collision, set the carry flag.
			if g.Pixels[a] == 0x01 {
				collision = true
			}

			// XOR the bit value.
			g.Pixels[a] = g.Pixels[a] ^ v
		}
	}

	return
}

func (g *Graphics) String() string {
	var s string

	for y := 0; y < GraphicsHeight-1; y++ {
		for x := 0; x < GraphicsWidth-1; x++ {
			c := y*GraphicsWidth + x

			var v string
			if g.Pixels[c] == 0x01 {
				v = "X"
			}

			s += fmt.Sprintf("%s ", v)
		}
		fmt.Printf("\n")
	}

	return s
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

var (
	fg = termbox.ColorBlack
	bg = termbox.ColorDefault
)

// display is an implementation of the Display interface that renders
// the graphics array to the terminal.
type display struct{}

func (d *display) Init() error {
	if err := termbox.Init(); err != nil {
		return err
	}

	termbox.HideCursor()

	if err := termbox.Clear(bg, bg); err != nil {
		return err
	}

	return termbox.Flush()
}

func (d *display) Close() {
	termbox.Close()
}

func (d *display) Render(g *Graphics) error {
	for y := 0; y < GraphicsHeight-1; y++ {
		for x := 0; x < GraphicsWidth-1; x++ {
			c := y*GraphicsWidth + x
			v := ' '

			if g.Pixels[c] == 0x01 {
				v = '*'
			}

			termbox.SetCell(
				x,
				y,
				v,
				fg,
				bg,
			)
		}
	}

	return termbox.Flush()
}
