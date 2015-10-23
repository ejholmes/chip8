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
	// GetKey waits for input on the keyboard and returns the key that was
	// pressed.
	GetKey() (byte, error)
}

// Keypad func can be used to wrap a function that returns a byte as a Keypad.
type KeypadFunc func() (byte, error)

func (f KeypadFunc) GetKey() (byte, error) {
	return f()
}

// NullKeypad is a Keypad that just returns an error.
var NullKeypad = KeypadFunc(func() (byte, error) {
	return 0x00, errors.New("null keypad not usable")
})

// TermboxKeypad is a Keypad implementation that maps keys from a standard
// keyboard to the CHIP-8 keyboard and uses termbox to poll for events.
type TermboxKeypad struct{}

func NewTermboxKeypad() *TermboxKeypad {
	return &TermboxKeypad{}
}

var keyMap = map[rune]byte{
	'1': 0x01, '2': 0x02, '3': 0x03, '4': 0x0C,
	'q': 0x04, 'w': 0x05, 'e': 0x06, 'r': 0x0D,
	'a': 0x07, 's': 0x08, 'd': 0x09, 'f': 0x0E,
	'z': 0x0A, 'x': 0x00, 'c': 0x0B, 'v': 0x0F,
}

// Get waits for a keypress.
func (k *TermboxKeypad) GetKey() (byte, error) {
	event := termbox.PollEvent()
	return keyMap[event.Ch], nil
}
