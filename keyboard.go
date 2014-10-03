// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chip8

var DefaultKeyboard = &keyboard{}

// Keyboard represents a CHIP-8 Keyboard.
type Keyboard interface {
	// Get waits for input on the keyboard and returns the key that was
	// pressed.
	Get() (byte, error)
}

// keyboard is a Keyboard implementation that maps keys
// from a standard keyboard to the CHIP-8 keyboard.
type keyboard struct{}

// Get waits for a keypress.
func (k *keyboard) Get() (byte, error) {
	//var b string
	//if _, err := fmt.Scan(&b); err != nil {
	//return 0, err
	//}

	return 0x01, nil
}
