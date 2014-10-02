// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chip8

var (
	DefaultDisplay  = Display(&terminalDisplay{})
	DefaultGraphics = &Graphics{}
)

// Display represents the output display for the CHIP-8 graphics array.
type Display interface {
	// Render should render the current graphics array to the display.
	Render(*Graphics) error
}

// Graphics represents the graphics array for the CHIP-8.
type Graphics struct {
	// The raw pixels of the graphics array.
	Pixels [64 * 32]byte

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

// terminalDisplay is an implementation of the Display interface that renders
// the graphics array to the terminal.
type terminalDisplay struct {
}

func (d *terminalDisplay) Render(g *Graphics) error {
	return nil
}
