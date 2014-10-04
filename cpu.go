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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

// Sensible defaults
var (
	// DefaultLogger is the default logger to use. Defaults to logging to /dev/null
	DefaultLogger = log.New(ioutil.Discard, "", 0)

	// DefaultClockSpeed is the default clock speed of the CPU. The CHIP-8
	// operated at 60 Hz.
	DefaultClockSpeed = time.Duration(60) // Hz

	// DefaultOptions is the default set of options that's used when calling
	// NewCPU.
	DefaultOptions = &Options{
		ClockSpeed: DefaultClockSpeed,
	}
)

// CPU represents a CHIP-8 CPU.
type CPU struct {
	// The 4096 bytes of memory.
	//
	// Memory Map:
	// +---------------+= 0xFFF (4095) End of Chip-8 RAM
	// |               |
	// |               |
	// |               |
	// |               |
	// |               |
	// | 0x200 to 0xFFF|
	// |     Chip-8    |
	// | Program / Data|
	// |     Space     |
	// |               |
	// |               |
	// |               |
	// +- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
	// |               |
	// |               |
	// |               |
	// +---------------+= 0x200 (512) Start of most Chip-8 programs
	// | 0x000 to 0x1FF|
	// | Reserved for  |
	// |  interpreter  |
	// +---------------+= 0x000 (0) Start of Chip-8 RAM
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
	Stack [16]uint16

	// Stack pointer.
	SP byte

	// The CHIP-8 timers count down at 60 Hz, so we slow down the cpu clock
	// to only execute 60 opcodes per second.
	Clock <-chan time.Time

	// The graphics array.
	Graphics

	// The connected Keypad. The zero value is the DefaultKeypad.
	Keypad Keypad

	// A logger to log information about the CPU while it's executing. The
	// zero value is the DefaultLogger.
	Logger *log.Logger

	// Settable in tests.
	randByteFunc func() byte
}

// Options provides a means of configuring the CPU.
type Options struct {
	ClockSpeed time.Duration
}

// NewCPU returns a new CPU instance.
func NewCPU(options *Options) (*CPU, error) {
	if options == nil {
		options = DefaultOptions
	}

	c := &CPU{
		PC:    0x200,
		Clock: time.Tick(time.Second / options.ClockSpeed),
	}

	return c, c.init()
}

// Load reads from the reader and loads the bytes into memory starting at
// address 200.
func (c *CPU) Load(r io.Reader) (int, error) {
	return c.load(0x200, r)
}

// LoadBytes loads the bytes into memory.
func (c *CPU) LoadBytes(p []byte) (int, error) {
	return c.Load(bytes.NewReader(p))
}

func (c *CPU) load(offset int, r io.Reader) (int, error) {
	return r.Read(c.Memory[offset:])
}

// init loads initalizes the cpu by loading the fontset into RAM.
func (c *CPU) init() error {
	if _, err := c.load(0, bytes.NewReader(FontSet)); err != nil {
		return fmt.Errorf("chip8: could not load font set: %s", err.Error())
	}

	return nil
}

// Step runs a single CPU cycle.
func (c *CPU) Step() (uint16, error) {
	// Simulate the clock speed of the CHIP-8 CPU.
	<-c.Clock

	// Decode the opcode.
	op := c.op()

	// Dispatch the opcode.
	if err := c.Dispatch(c.op()); err != nil {
		return op, err
	}

	return op, nil
}

// Run does the thing.
func (c *CPU) Run() error {
	for {
		op, err := c.Step()
		if err != nil {
			return err
		}

		c.logger().Printf("op=0x%04X %s\n", op, c)
	}

	return nil
}

// Dispatch executes the given opcode.
func (c *CPU) Dispatch(op uint16) error {
	// In these listings, the following variables are used:
	//
	// nnn or addr - A 12-bit value, the lowest 12 bits of the instruction
	// n or nibble - A 4-bit value, the lowest 4 bits of the instruction
	// x - A 4-bit value, the lower 4 bits of the high byte of the instruction
	// y - A 4-bit value, the upper 4 bits of the low byte of the instruction
	// kk or byte - An 8-bit value, the lowest 8 bits of the instruction

	switch op & 0xF000 {
	// 0nnn - SYS addr
	case 0x0000:
		switch op {
		// 00E0 - CLS
		case 0x00E0:
			c.Graphics.Clear()

			c.PC += 2

			break

		// 00EE - RET
		case 0x00EE:
			// Return from a subroutine.
			//
			// The interpreter sets the program counter to the
			// address at the top of the stack, then subtracts 1
			// from the stack pointer.

			c.PC = c.Stack[c.SP]
			c.SP--

			break

		default:
			// Jump to a machine code routine at nnn.
			//
			// This instruction is only used on the old computers on
			// which Chip-8 was originally implemented. It is
			// ignored by modern interpreters.

			return &UnknownOpcode{Opcode: op}
		}

		break

	// 1nnn - JP addr
	case 0x1000:
		// Jump to location nnn.
		//
		// The interpreter sets the program counter to nnn.

		c.PC = op & 0x0FFF

		break

	// 2nnn - CALL addr
	case 0x2000:
		// Call subroutine at nnn.
		//
		// The interpreter increments the stack pointer, then puts the
		// current PC on the top of the stack. The PC is then set to
		// nnn.

		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = op & 0x0FFF

		break

	// 3xkk - SE Vx, byte
	case 0x3000:
		// Skip next instruction if Vx = kk.
		//
		// The interpreter compares register Vx to kk, and if they are
		// equal, increments the program counter by 2.

		x := (op & 0x0F00) >> 8
		kk := byte(op)

		c.PC += 2
		if c.V[x] == kk {
			c.PC += 2
		}

		break

	// 4xkk - SNE Vx, byte
	case 0x4000:
		// Skip next instruction if Vx != kk.
		//
		// The interpreter compares register Vx to kk, and if they are
		// not equal, increments the program counter by 2.

		x := (op & 0x0F00) >> 8
		kk := byte(op)

		c.PC += 2
		if c.V[x] != kk {
			c.PC += 2
		}

		break

	// 5xy0 - SE Vx, Vy
	case 0x5000:
		switch op & 0xF00F {
		case 0x5000:
			// Skip next instruction if Vx = Vy.
			//
			// The interpreter compares register Vx to register Vy, and if
			// they are equal, increments the program counter by 2.

			x := (op & 0x0F00) >> 8
			y := (op & 0x00F0) >> 4

			c.PC += 2
			if c.V[x] == c.V[y] {
				c.PC += 2
			}

			break
		default:
			return &UnknownOpcode{Opcode: op}
		}

		break

	// 6xkk - LD Vx, byte
	case 0x6000:
		// Set Vx = kk.
		//
		// The interpreter puts the value kk into register Vx.

		x := (op & 0x0F00) >> 8
		kk := byte(op)

		c.V[x] = kk

		c.PC += 2

		break

	// 7xkk - ADD Vx, byte
	case 0x7000:
		// Set Vx = Vx + kk.
		//
		// Adds the value kk to the value of register Vx, then stores
		// the result in Vx.

		x := (op & 0x0F00) >> 8
		kk := byte(op)

		c.V[x] = c.V[x] + kk

		c.PC += 2

		break

	case 0x8000:
		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4

		switch op & 0x000F {
		// 8xy0 - LD Vx, Vy
		case 0x0000:
			// Set Vx = Vy.
			//
			// Stores the value of register Vy in register Vx.

			c.V[x] = c.V[y]

			c.PC += 2

			break

		// 8xy1 - OR Vx, Vy
		case 0x0001:
			// Set Vx = Vx OR Vy.
			//
			// Performs a bitwise OR on the values of Vx and Vy,
			// then stores the result in Vx. A bitwise OR compares
			// the corrseponding bits from two values, and if either
			// bit is 1, then the same bit in the result is also 1.
			// Otherwise, it is 0.

			c.V[x] = c.V[y] | c.V[x]

			c.PC += 2

			break

		// 8xy2 - AND Vx, Vy
		case 0x0002:
			// Set Vx = Vx AND Vy.
			//
			// Performs a bitwise AND on the values of Vx and Vy,
			// then stores the result in Vx. A bitwise AND compares
			// the corrseponding bits from two values, and if both
			// bits are 1, then the same bit in the result is also 1.
			// Otherwise, it is 0.

			c.V[x] = c.V[y] & c.V[x]

			c.PC += 2

			break

		// 8xy3 - XOR Vx, Vy
		case 0x0003:
			// Set Vx = Vx XOR Vy.
			//
			// Performs a bitwise exclusive OR on the values of Vx
			// and Vy, then stores the result in Vx. An exclusive OR
			// compares the corrseponding bits from two values, and
			// if the bits are not both the same, then the
			// corresponding bit in the result is set to 1.
			// Otherwise, it is 0.

			c.V[x] = c.V[y] ^ c.V[x]

			c.PC += 2

			break

		// 8xy4 - ADD Vx, Vy
		case 0x0004:
			// Set Vx = Vx + Vy, set VF = carry.
			//
			// The values of Vx and Vy are added together. If the
			// result is greater than 8 bits (i.e., > 255,) VF is
			// set to 1, otherwise 0. Only the lowest 8 bits of the
			// result are kept, and stored in Vx.

			r := uint16(c.V[x]) + uint16(c.V[y])

			var cf byte
			if r > 0xFF {
				cf = 1
			}
			c.V[0xF] = cf

			c.V[x] = byte(r)

			c.PC += 2

			break

		// 8xy5 - SUB Vx, Vy
		case 0x0005:
			// Set Vx = Vx - Vy, set VF = NOT borrow.
			//
			// If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy
			// is subtracted from Vx, and the results stored in Vx.

			var cf byte
			if c.V[x] > c.V[y] {
				cf = 1
			}
			c.V[0xF] = cf

			c.V[x] = c.V[x] - c.V[y]

			c.PC += 2

			break

		// 8xy6 - SHR Vx {, Vy}
		case 0x0006:
			// Set Vx = Vx SHR 1.
			//
			// If the least-significant bit of Vx is 1, then VF is
			// set to 1, otherwise 0. Then Vx is divided by 2.

			var cf byte
			if (c.V[x] & 0x01) == 0x01 {
				cf = 1
			}
			c.V[0xF] = cf

			c.V[x] = c.V[x] / 2

			c.PC += 2

			break

		// 8xy7 - SUBN Vx, Vy
		case 0x0007:
			// Set Vx = Vy - Vx, set VF = NOT borrow.
			//
			// If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx
			// is subtracted from Vy, and the results stored in Vx.

			var cf byte
			if c.V[y] > c.V[x] {
				cf = 1
			}
			c.V[0xF] = cf

			c.V[x] = c.V[y] - c.V[x]

			c.PC += 2

			break

		// 8xyE - SHL Vx {, Vy}
		case 0x000E:
			// Set Vx = Vx SHL 1.
			//
			// If the most-significant bit of Vx is 1, then VF is
			// set to 1, otherwise to 0. Then Vx is multiplied by 2.

			var cf byte
			if (c.V[x] & 0x80) == 0x80 {
				cf = 1
			}
			c.V[0xF] = cf

			c.V[x] = c.V[x] * 2

			c.PC += 2

			break
		}

		break

	// Skips the next instruction if VX doesn't equal VY.
	//   0x9XY0
	case 0x9000:
		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4

		switch op & 0x000F {
		// 9xy0 - SNE Vx, Vy
		case 0x0000:
			// Skip next instruction if Vx != Vy.
			//
			// The values of Vx and Vy are compared, and if they are
			// not equal, the program counter is increased by 2.

			c.PC += 2
			if c.V[x] != c.V[y] {
				c.PC += 2
			}

			break
		default:
			return &UnknownOpcode{Opcode: op}
		}

		break

	// Annn - LD I, addr
	case 0xA000:
		// Set I = nnn.
		//
		// The value of register I is set to nnn.

		c.I = op & 0x0FFF
		c.PC += 2

		break

	// Bnnn - JP V0, addr
	case 0xB000:
		// Jump to location nnn + V0.
		//
		// The program counter is set to nnn plus the value of V0.

		c.PC = op&0x0FFF + uint16(c.V[0])

		break

	// Cxkk - RND Vx, byte
	case 0xC000:
		// Set Vx = random byte AND kk.
		//
		// The interpreter generates a random number from 0 to 255,
		// which is then ANDed with the value kk. The results are stored
		// in Vx. See instruction 8xy2 for more information on AND.

		x := (op & 0x0F00) >> 8
		kk := byte(op)

		c.V[x] = kk + c.randByte()

		c.PC += 2

		break

	// Dxyn - DRW Vx, Vy, nibble
	case 0xD000:
		// Display n-byte sprite starting at memory location I at (Vx,
		// Vy), set VF = collision.
		//
		// The interpreter reads n bytes from memory, starting at the
		// address stored in I. These bytes are then displayed as
		// sprites on screen at coordinates (Vx, Vy). Sprites are XORed
		// onto the existing screen. If this causes any pixels to be
		// erased, VF is set to 1, otherwise it is set to 0. If the
		// sprite is positioned so part of it is outside the coordinates
		// of the display, it wraps around to the opposite side of the
		// screen. See instruction 8xy3 for more information on XOR, and
		// section 2.4, Display, for more information on the Chip-8
		// screen and sprites.

		var cf byte

		// The starting X coordinate on the graphics array.
		x := c.V[(op&0x0F00)>>8]

		// The starting Y coordinate on the graphics array.
		y := c.V[(op&0x00F0)>>4]

		// The height of the sprite.
		n := op & 0x000F

		if c.Graphics.WriteSprite(c.Memory[c.I:c.I+n], x, y) {
			cf = 0x01
		}

		c.V[0xF] = cf
		c.PC += 2

		c.Graphics.Draw()

		break

	case 0xE000:
		switch op & 0x00FF {
		// Skips the next instruction if the key stored in VX is pressed.
		case 0x9E:

		// Skips the next instruction if the key stored in VX isn't pressed.
		case 0xA1:
		}
	case 0xF000:
		x := (op & 0x0F00) >> 8

		switch op & 0x00FF {
		// Sets VX to the value of the delay timer.
		case 0x07:

		// Fx0A - LD Vx, K
		case 0x0A:
			// Wait for a key press, store the value of the key in
			// Vx.
			//
			// All execution stops until a key is pressed, then the
			// value of that key is stored in Vx.

			b, err := c.getKey()
			if err != nil {
				return err
			}

			c.V[x] = b

			c.PC += 2

			break

		// Sets the delay timer to VX.
		case 0x15:

		// Sets the sound timer to VX.
		case 0x18:

		// Fx1E - ADD I, Vx
		case 0x1E:
			// Set I = I + Vx.
			//
			// The values of I and Vx are added, and the results are
			// stored in I.

			c.I = c.I + uint16(c.V[x])

			break

		// Fx29 - LD F, Vx
		case 0x29:
			// Set I = location of sprite for digit Vx.
			//
			// The value of I is set to the location for the
			// hexadecimal sprite corresponding to the value of Vx.
			// See section 2.4, Display, for more information on the
			// Chip-8 hexadecimal font.

			c.I = uint16(c.V[x]) * uint16(0x05)

			c.PC += 2

			break

		// Fx33 - LD B, Vx
		case 0x33:
			// Store BCD representation of Vx in memory locations I,
			// I+1, and I+2.
			//
			// The interpreter takes the decimal value of Vx, and
			// places the hundreds digit in memory at location in I,
			// the tens digit at location I+1, and the ones digit at
			// location I+2.

			c.Memory[c.I] = c.V[x] / 100
			c.Memory[c.I+1] = (c.V[x] / 10) % 10
			c.Memory[c.I+2] = (c.V[x] % 100) % 10

			c.PC += 2

			break

		// Stores V0 to VX in memory starting at address I.
		case 0x55:

		// Fx65 - LD Vx, [I]
		case 0x65:
			// Read registers V0 through Vx from memory starting at
			// location I.
			//
			// The interpreter reads values from memory starting at
			// location I into registers V0 through Vx.

			for i := 0; byte(i) <= byte(x); i++ {
				c.V[uint16(i)] = c.Memory[c.I+uint16(i)]
			}

			c.PC += 2

			break
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

func (c *CPU) randByte() byte {
	if c.randByteFunc == nil {
		return randByte()
	}

	return c.randByteFunc()
}

func (c *CPU) getKey() (byte, error) {
	b, err := c.keypad().Get()
	if err != nil {
		return b, fmt.Errorf("chip8: unable to get key from keypad: %s", err.Error())
	}

	return b, nil
}

func (c *CPU) keypad() Keypad {
	if c.Keypad == nil {
		return DefaultKeypad
	}

	return c.Keypad
}

// String implements the fmt.Stringer interface.
func (c *CPU) String() string {
	return fmt.Sprintf(
		"I=0x%04X pc=0x%04X V[x]=%v stack=%v SP=0x%04X",
		c.I, c.PC, c.V, c.Stack, c.SP,
	)
}

// logger returns the logger to use for debugging.
func (c *CPU) logger() *log.Logger {
	if c.Logger == nil {
		return DefaultLogger
	}

	return c.Logger
}

// UnknownOpcode is return when the opcode is not recognized.
type UnknownOpcode struct {
	Opcode uint16
}

func (e *UnknownOpcode) Error() string {
	return fmt.Sprintf("chip8: unknown opcode: 0x%04X", e.Opcode)
}

// randByte returns a random byte.
func randByte() byte {
	return byte(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(255))
}
