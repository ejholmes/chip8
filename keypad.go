// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chip8

import (
	"errors"

	"github.com/nsf/termbox-go"
)

// ErrQuit can be returned when we need to indicate that the user wants to quit.
var ErrQuit = errors.New("quit key pressed")

// Keypad represents a CHIP-8 Keypad.
type Keypad interface {
	// Do any initialization.
	Init() error

	// Get waits for input on the keyboard and returns the key that was
	// pressed.
	Get() (byte, error)
}

// Keypad func can be used to wrap a function that returns a byte as a Keypad.
type KeypadFunc func() (byte, error)

func (f KeypadFunc) Init() error {
	return nil
}

func (f KeypadFunc) Get() (byte, error) {
	return f()
}

// keyboard is a Keypad implementation that maps keys
// from a standard keyboard to the CHIP-8 keyboard.
type keyboard struct{}

func (k *keyboard) Init() error {
	return nil
}

// Get waits for a keypress.
func (k *keyboard) Get() (byte, error) {
	event := termbox.PollEvent()

	if event.Ch == 'q' {
		return 0x00, ErrQuit
	}

	return 0x01, nil
}

// polls for keyboard events.
func (k *keyboard) poll() {
}
