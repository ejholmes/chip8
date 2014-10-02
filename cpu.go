// Copyright 2014 Eric Holmes.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chip8 provides a Go implementation of the CHIP-8 emulator.
//
// CHIP-8 was most commonly implemented on 4K systems, such as the
// Cosmac VIP and the Telemac 1800. These machines had 4096 (0x1000)
// memory locations, all of which are 8 bits (a byte) which is where the
// term CHIP-8 originated. However, the CHIP-8 interpreter itself
// occupies the first 512 bytes of the memory space on these machines.
// For this reason, most programs written for the original system begin
// at memory location 512 (0x200) and do not access any of the memory
// below the location 512 (0x200). The uppermost 256 bytes (0xF00-0xFFF)
// are reserved for display refresh, and the 96 bytes below that
// (0xEA0-0XEFF) were reserved for call stack, internal use, and other
// variables.
package chip8

import (
	"fmt"
	"time"
)

// Sensible defaults
var (
	DefaultClockSpeed = time.Duration(60) // Hz
	DefaultOptions    = &Options{
		ClockSpeed: DefaultClockSpeed,
	}
)

// CPU represents a CHIP-8 CPU.
type CPU struct {
	// The 4096 bytes of memory.
	Memory [4096]byte

	// The address register, which is named I, is 16 bits wide and is used
	// with several opcodes that involve memory operations.
	I uint16

	// Program counter.
	PC uint16

	// CHIP-8 has 16 8-bit data registers named from V0 to VF. The VF
	// register doubles as a carry flag.
	V [16]byte

	// The stack is only used to store return addresses when subroutines are
	// called. The original 1802 version allocated 48 bytes for up to 12
	// levels of nesting; modern implementations normally have at least 16
	// levels.
	Stack [16]byte

	// Stack pointer.
	SP uint16

	// The CHIP-8 timers count down at 60 Hz, so we slow down the cpu clock
	// to only execute 60 opcodes per second.
	Clock <-chan time.Time
}

// Options provides a means of configuring the CPU.
type Options struct {
	ClockSpeed time.Duration
}

// NewCPU returns a new CPU instance.
func NewCPU(options *Options) *CPU {
	if options == nil {
		options = DefaultOptions
	}

	return &CPU{
		PC:    200,
		Clock: time.Tick(time.Second / options.ClockSpeed),
	}
}

// Step runs a single CPU cycle.
func (c *CPU) Step() error {
	// Simulate the clock speed of the CHIP-8 CPU.
	<-c.Clock

	// Dispatch the opcode.
	if err := c.Dispatch(c.op()); err != nil {
		return err
	}

	// Increment the program counter by 2.
	c.PC = c.PC + 2

	return nil
}

// Dispatch executes the given opcode.
func (c *CPU) Dispatch(op uint16) error {
	switch op & 0xF000 {
	//   0x0NNN
	case 0x0000:
		switch op {
		// Clears the screen.
		case 0x00E0:

		// Returns from a subroutine.
		case 0x00EE:

		// Calls RCA 1802 program at address NNNN
		default:
		}

	// Jumps to address NNN
	//   0x1NNN
	case 0x1000:

	// Calls subroutine at NNN.
	//   0x2NNN
	case 0x2000:

	// Skip the next instruction if VX equals NN.
	//   0x3XNN
	case 0x3000:

	// Skips the next instruction if VX doesn't equal NN.
	//   0x4XNN
	case 0x4000:

	// Skips the next instruction if VX equals VY.
	//   0x5XY0
	case 0x5000:

	// Sets VX to NN.
	//   0x6XNN
	case 0x6000:

	// Adds NN to VX
	//   0x7XNN
	case 0x7000:

	case 0x8000:

		switch op & 0x000F {
		// Sets VX to the value of VY.
		case 0x00:

		// Sets VX to VX or VY.
		case 0x01:

		// Sets VX to VX and VY.
		case 0x02:

		// Sets VX to VX xor VY.
		case 0x03:

		// Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
		case 0x04:

		// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
		case 0x05:
		// Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift.
		case 0x06:
		// Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
		case 0x07:
		// Shifts VX left by one. VF is set to the value of the most significant bit of VX before the shift.
		case 0x0E:
		}

	// Skips the next instruction if VX doesn't equal VY.
	//   0x9XY0
	case 0x9000:

	// Sets I to the address NNN.
	//   0xANNN
	case 0xA000:
		c.I = op & 0x0FFF
		break

	// Jumps to the address NNN plus V0.
	//   0xBNNN
	case 0xB000:

	// Sets VX to a random number and NN.
	//   0xCXNN
	case 0xC000:

	// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a
	// height of N pixels. Each row of 8 pixels is read as bit-coded (with the
	// most significant bit of each byte displayed on the left) starting from
	// memory location I; I value doesn't change after the execution of this
	// instruction. As described above, VF is set to 1 if any screen pixels are
	// flipped from set to unset when the sprite is drawn, and to 0 if that doesn't
	// happen.
	//
	//   0xDXYN
	case 0xD000:

	case 0xE000:
		switch op & 0x00FF {
		// Skips the next instruction if the key stored in VX is pressed.
		case 0x9E:

		// Skips the next instruction if the key stored in VX isn't pressed.
		case 0xA1:
		}
	case 0xF000:
		switch op & 0x00FF {
		// Sets VX to the value of the delay timer.
		case 0x07:

		// A key press is awaited, and then stored in VX.
		case 0x0A:

		// Sets the delay timer to VX.
		case 0x15:

		// Sets the sound timer to VX.
		case 0x18:

		// Adds VX to I.
		case 0x1E:

		// Sets I to the location of the sprite for the character in VX. Characters
		// 0-F (in hexadecimal) are represented by a 4x5 font.
		case 0x29:

		// Stores the Binary-coded decimal representation of VX, with the most
		// significant of three digits at the address in I, the middle digit at
		// I plus 1, and the least significant digit at I plus 2. (In other words,
		// take the decimal representation of VX, place the hundreds digit in
		// memory at location in I, the tens digit at location I+1, and the ones
		// digit at location I+2.)
		case 0x33:

		// Stores V0 to VX in memory starting at address I.
		case 0x55:

		// Fills V0 to VX with values from memory starting at address I.
		case 0x65:
		}
	default:
		return &UnknownOpcode{Opcode: op}
	}

	return nil
}

// op returns the next op code.
func (c *CPU) op() uint16 {
	return uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
}

// UnknownOpcode is return when the opcode is not recognized.
type UnknownOpcode struct {
	Opcode uint16
}

func (e *UnknownOpcode) Error() string {
	return fmt.Sprintf("chip8: unknown opcode: 0x%4X", e.Opcode)
}
