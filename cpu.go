package chip8

import "errors"

// ErrUnkownOpcode is returned when we try to execute an unkown opcode.
var ErrUnkownOpcode = errors.New("chip8: unknown opcode")

// CPU represents a Chip8 CPU.
type CPU struct {
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
}

// NewCPU returns a new CPU instance.
func NewCPU() *CPU {
	return &CPU{
		PC: 200,
	}
}

// Cycle runs a single CPU cycle.
func (c *CPU) Cycle() error {
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
	case 0xA000: // ANNN: Sets I to the address NNN
		c.I = op & 0x0FFF
		break
	default:
		return ErrUnkownOpcode
	}

	return nil
}

// op returns the next op code.
func (c *CPU) op() uint16 {
	return uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
}
