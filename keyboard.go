// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chip8

import (
	"errors"

	"github.com/nsf/termbox-go"
)

var DefaultKeyboard = &keyboard{}

var (
	ErrQuit = errors.New("quit key pressed")
)

// Keyboard represents a CHIP-8 Keyboard.
type Keyboard interface {
	// Do any initialization.
	Init() error

	// Get waits for input on the keyboard and returns the key that was
	// pressed.
	Get() (byte, error)
}

// keyboard is a Keyboard implementation that maps keys
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
