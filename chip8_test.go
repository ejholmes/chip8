package chip8

import "testing"

func checkHex(t *testing.T, subject string, got, want uint16) {
	if got != want {
		t.Errorf("%s => 0x%04X; want 0x%04x", subject, got, want)
	}
}

func TestCPU_Step(t *testing.T) {
	c := NewCPU(nil)
	c.Memory[200] = 0xA1
	c.Memory[201] = 0x00

	if err := c.Step(); err != nil {
		t.Fatal(err)
	}

	checkHex(t, "PC", c.PC, uint16(202))
}

func TestCPU_Dispatch(t *testing.T) {
	tests := []struct {
		op     uint16
		before func(*CPU)
		check  func(*CPU)
	}{
		{
			uint16(0xA100),
			nil,
			func(c *CPU) {
				checkHex(t, "I", c.I, uint16(0x100))
			},
		},

		{
			uint16(0x2100),
			nil,
			func(c *CPU) {
				checkHex(t, "Stack[0]", c.Stack[0], uint16(0xC8))
				checkHex(t, "SP", uint16(c.SP), uint16(0x1))
				checkHex(t, "PC", c.PC, uint16(0x100))
			},
		},

		{
			uint16(0x3123),
			nil,
			func(c *CPU) {
				checkHex(t, "PC", c.PC, uint16(200))
			},
		},

		{
			uint16(0x3123),
			func(c *CPU) {
				c.V[1] = 0x03
			},
			func(c *CPU) {
				checkHex(t, "PC", c.PC, uint16(202))
			},
		},

		{
			uint16(0x4123),
			nil,
			func(c *CPU) {
				checkHex(t, "PC", c.PC, uint16(202))
			},
		},

		{
			uint16(0x4123),
			func(c *CPU) {
				c.V[1] = 0x03
			},
			func(c *CPU) {
				checkHex(t, "PC", c.PC, uint16(200))
			},
		},

		{
			uint16(0x5120),
			func(c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x04
			},
			func(c *CPU) {
				checkHex(t, "PC", c.PC, uint16(200))
			},
		},

		{
			uint16(0x5120),
			func(c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x03
			},
			func(c *CPU) {
				checkHex(t, "PC", c.PC, uint16(202))
			},
		},
	}

	for _, tt := range tests {
		c := NewCPU(nil)
		if tt.before != nil {
			tt.before(c)
		}
		c.Dispatch(tt.op)
		tt.check(c)
	}
}

func TestCPU_op(t *testing.T) {
	c := NewCPU(nil)
	c.Memory[200] = 0xA2
	c.Memory[201] = 0xF0

	checkHex(t, "op", c.op(), uint16(0xA2F0))
}
