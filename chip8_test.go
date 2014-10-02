package chip8

import "testing"

var opcodeTests = map[string][]struct {
	op     uint16
	before func(*testing.T, *CPU)
	check  func(*testing.T, *CPU)
}{
	"2nnn - CALL addr": {
		{
			uint16(0x2100),
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "Stack[0]", c.Stack[0], uint16(0xC8))
				checkHex(t, "SP", uint16(c.SP), uint16(0x1))
				checkHex(t, "PC", c.PC, uint16(0x100))
			},
		},
	},

	"3xkk - SE Vx, byte": {
		{
			uint16(0x3123),
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(200))
			},
		},

		{
			uint16(0x3103),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(202))
			},
		},
	},

	"4xkk - SNE Vx, byte": {
		{
			uint16(0x4123),
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(202))
			},
		},

		{
			uint16(0x4103),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(200))
			},
		},
	},

	"5xy0 - SE Vx, Vy": {
		{
			uint16(0x5120),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x04
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(200))
			},
		},

		{
			uint16(0x5120),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(202))
			},
		},
	},

	"6xkk - LD Vx, byte": {
		{
			uint16(0x6102),
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x02))
			},
		},
	},

	"7xkk - ADD Vx, byte": {
		{
			uint16(0x7102),
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x02))
			},
		},

		{
			uint16(0x7102),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x03))
			},
		},
	},

	"8xy0 - LD Vx, Vy": {
		{
			uint16(0x8120),
			func(t *testing.T, c *CPU) {
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x01))
			},
		},
	},

	"8xy1 - OR Vx, Vy": {
		{
			uint16(0x8121),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x10
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x11))
			},
		},
	},

	"8xy2 - AND Vx, Vy": {
		{
			uint16(0x8122),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x10
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x00))
			},
		},
	},

	"8xy3 - XOR Vx, Vy": {
		{
			uint16(0x8123),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x00))
			},
		},
	},

	"8xy4 - ADD Vx, Vy": {
		{
			uint16(0x8124),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x2))
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x0))
			},
		},

		{
			uint16(0x8124),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0xFF
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x2))
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x1))
			},
		},
	},

	"8xy5 - SUB Vx, Vy": {
		{
			uint16(0x8125),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0xFF
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x1))
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0xFC))
			},
		},

		{
			uint16(0x8125),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x02
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x0))
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0xFF))
			},
		},
	},

	"8xy6 - SHR Vx {, Vy}": {
		{
			uint16(0x8126),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x1))
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x1))
			},
		},

		{
			uint16(0x8126),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x0))
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x1))
			},
		},
	},

	"8xy7 - SUBN Vx, Vy": {
		{
			uint16(0x8127),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0xFF
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x1))
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0xFC))
			},
		},

		{
			uint16(0x8127),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x0))
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0xFF))
			},
		},
	},

	"8xyE - SHL Vx {, Vy}": {
		{
			uint16(0x812E),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x0))
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x2))
			},
		},

		{
			uint16(0x812E),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x81
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", uint16(c.V[0xF]), uint16(0x1))
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x2))
			},
		},
	},

	"9xy0 - SNE Vx, Vy": {
		{
			uint16(0x9120),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(202))
			},
		},

		{
			uint16(0x9120),
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(200))
			},
		},
	},

	"Annn - LD I, addr": {
		{
			uint16(0xA100),
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "I", c.I, uint16(0x100))
			},
		},
	},

	"Bnnn - JP V0, addr": {
		{
			uint16(0xB210),
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(0x210))
			},
		},

		{
			uint16(0xB210),
			func(t *testing.T, c *CPU) {
				c.V[0] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, uint16(0x211))
			},
		},
	},

	"Cxkk - RND Vx, byte": {
		{
			uint16(0xC110),
			func(t *testing.T, c *CPU) {
				c.randByteFunc = func() byte {
					return 0x01
				}
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", uint16(c.V[1]), uint16(0x11))
			},
		},
	},
}

func checkHex(t *testing.T, subject string, got, want uint16) {
	if got != want {
		t.Errorf("%s => 0x%04X; want 0x%04X", subject, got, want)
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

func TestOpcodes(t *testing.T) {
	for i, tests := range opcodeTests {
		for _, tt := range tests {
			c := NewCPU(nil)
			if tt.before != nil {
				tt.before(t, c)
			}
			c.Dispatch(tt.op)
			tt.check(t, c)

			if t.Failed() {
				t.Logf("==============")
				t.Logf("Instruction: %s", i)
				t.Logf("Opcode: 0x%04X", tt.op)
				t.Logf("CPU: %v", c)
				t.Logf("==============")
				t.FailNow()
			}
		}
	}
}

func TestCPU_op(t *testing.T) {
	c := NewCPU(nil)
	c.Memory[200] = 0xA2
	c.Memory[201] = 0xF0

	checkHex(t, "op", c.op(), uint16(0xA2F0))
}
