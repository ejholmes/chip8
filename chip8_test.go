package chip8

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func init() {
	// mock out randByte func to return a deterministic value.
	randByte = func() byte {
		return 0x01
	}
}

var opcodeTests = map[string][]struct {
	op     int
	before func(*testing.T, *CPU)
	check  func(*testing.T, *CPU)
}{
	"2nnn - CALL addr": {
		{
			0x2100,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "Stack[1]", c.Stack[1], 0x200)
				checkHex(t, "SP", c.SP, 0x1)
				checkHex(t, "PC", c.PC, 0x100)
			},
		},
	},

	"3xkk - SE Vx, byte": {
		{
			0x3123,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},

		{
			0x3103,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x204)
			},
		},
	},

	"4xkk - SNE Vx, byte": {
		{
			0x4123,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x204)
			},
		},

		{
			0x4103,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},
	},

	"5xy0 - SE Vx, Vy": {
		{
			0x5120,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x04
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},

		{
			0x5120,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x204)
			},
		},
	},

	"6xkk - LD Vx, byte": {
		{
			0x6102,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x02)
			},
		},
	},

	"7xkk - ADD Vx, byte": {
		{
			0x7102,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x02)
			},
		},

		{
			0x7102,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x03)
			},
		},
	},

	"8xy0 - LD Vx, Vy": {
		{
			0x8120,
			func(t *testing.T, c *CPU) {
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x01)
			},
		},
	},

	"8xy1 - OR Vx, Vy": {
		{
			0x8121,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x10
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x11)
			},
		},
	},

	"8xy2 - AND Vx, Vy": {
		{
			0x8122,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x10
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x00)
			},
		},
	},

	"8xy3 - XOR Vx, Vy": {
		{
			0x8123,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x00)
			},
		},
	},

	"8xy4 - ADD Vx, Vy": {
		{
			0x8124,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x2)
				checkHex(t, "VF", c.V[0xF], 0x0)
			},
		},

		{
			0x8124,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0xFF
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x2)
				checkHex(t, "VF", c.V[0xF], 0x1)
			},
		},
	},

	"8xy5 - SUB Vx, Vy": {
		{
			0x8125,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0xFF
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x1)
				checkHex(t, "V[1]", c.V[1], 0xFC)
			},
		},

		{
			0x8125,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x02
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x0)
				checkHex(t, "V[1]", c.V[1], 0xFF)
			},
		},
	},

	"8xy6 - SHR Vx {, Vy}": {
		{
			0x8126,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x1)
				checkHex(t, "V[1]", c.V[1], 0x1)
			},
		},

		{
			0x8126,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x0)
				checkHex(t, "V[1]", c.V[1], 0x1)
			},
		},
	},

	"8xy7 - SUBN Vx, Vy": {
		{
			0x8127,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0xFF
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x1)
				checkHex(t, "V[1]", c.V[1], 0xFC)
			},
		},

		{
			0x8127,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x0)
				checkHex(t, "V[1]", c.V[1], 0xFF)
			},
		},
	},

	"8xyE - SHL Vx {, Vy}": {
		{
			0x812E,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x0)
				checkHex(t, "V[1]", c.V[1], 0x2)
			},
		},

		{
			0x812E,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x81
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x1)
				checkHex(t, "V[1]", c.V[1], 0x2)
			},
		},
	},

	"9xy0 - SNE Vx, Vy": {
		{
			0x9120,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x204)
			},
		},

		{
			0x9120,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},
	},

	"Annn - LD I, addr": {
		{
			0xA100,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "I", c.I, 0x100)
			},
		},
	},

	"Bnnn - JP V0, addr": {
		{
			0xB210,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x210)
			},
		},

		{
			0xB210,
			func(t *testing.T, c *CPU) {
				c.V[0] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x211)
			},
		},
	},

	"Cxkk - RND Vx, byte": {
		{
			0xC110,
			func(t *testing.T, c *CPU) {
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x11)
			},
		},
	},

	"Dxyn - DRW Vx, Vy, nibble": {
		{
			0xD001,
			func(t *testing.T, c *CPU) {
				c.I = 0x200
				c.Memory[0x200] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "Pixel", c.Pixels[0], 0x00)
				checkHex(t, "VF", c.V[0xF], 0x0)
			},
		},

		{
			0xD001,
			func(t *testing.T, c *CPU) {
				c.I = 0x200
				c.Memory[0x200] = 0x01
				c.Pixels[0x0] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "Pixel", c.Pixels[0], 0x01)
				checkHex(t, "VF", c.V[0xF], 0x1)
			},
		},

		{
			0xD005,
			func(t *testing.T, c *CPU) {
				c.I = 0x0
			},
			func(t *testing.T, c *CPU) {
				checkGraphics(t, &c.Graphics, "0926caffcf1848383015adc82f712f671dd1e0edc3f19560b4a44a8552eea2ba")
				checkHex(t, "VF", c.V[0xF], 0x00)
			},
		},

		{
			0xD015,
			func(t *testing.T, c *CPU) {
				c.V[0] = 0x00
				c.V[1] = 0x1C
			},
			func(t *testing.T, c *CPU) {
				checkGraphics(t, &c.Graphics, "796405cda1fa18bbd6e42dd2643af022793a37bc917b24c4bc8f88c242122a93")
			},
		},
		// Wrap y
		{
			0xD015,
			func(t *testing.T, c *CPU) {
				c.V[0] = 0x00
				c.V[1] = GraphicsHeight + 28
			},
			func(t *testing.T, c *CPU) {
				checkGraphics(t, &c.Graphics, "796405cda1fa18bbd6e42dd2643af022793a37bc917b24c4bc8f88c242122a93")
			},
		},
	},

	"Ex9E - SKP Vx": {
		{
			0xE19E,
			func(t *testing.T, c *CPU) {
				c.Keypad = KeypadFunc(func() (byte, error) {
					return 0x01, nil
				})
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},

		{
			0xE19E,
			func(t *testing.T, c *CPU) {
				c.V[0x01] = 0x02
				c.Keypad = KeypadFunc(func() (byte, error) {
					return 0x02, nil
				})
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x204)
			},
		},
	},

	"ExA1 - SKNP Vx": {
		{
			0xE1A1,
			func(t *testing.T, c *CPU) {
				c.Keypad = KeypadFunc(func() (byte, error) {
					return 0x01, nil
				})
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x204)
			},
		},

		{
			0xE1A1,
			func(t *testing.T, c *CPU) {
				c.V[0x01] = 0x02
				c.Keypad = KeypadFunc(func() (byte, error) {
					return 0x02, nil
				})
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},
	},

	"Fx07 - LD Vx, DT": {
		{
			0xF107,
			func(t *testing.T, c *CPU) {
				c.DT = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x01)
			},
		},
	},

	"Fx15 - LD DT, Vx": {
		{
			0xF115,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x05
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "DT", c.DT, 0x05)
			},
		},
	},

	"Fx18 - LD ST, Vx": {
		{
			0xF118,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x05
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "ST", c.ST, 0x05)
			},
		},
	},

	"Fx29 - LD F, Vx": {
		{
			0xF029,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "I", c.I, 0x00)
			},
		},

		{
			0xF129,
			func(t *testing.T, c *CPU) {
				c.V[0x01] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "I", c.I, 0x05)
			},
		},

		{
			0xF229,
			func(t *testing.T, c *CPU) {
				c.V[0x02] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "I", c.I, 0x0A)
			},
		},
	},

	"Fx33 - LD B, Vx": {
		{
			0xF033,
			func(t *testing.T, c *CPU) {
				c.V[0] = 0xFF
				c.I = 0x200
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "Memory[0x200]", c.Memory[0x200], 0x02)
				checkHex(t, "Memory[0x201]", c.Memory[0x201], 0x05)
				checkHex(t, "Memory[0x201]", c.Memory[0x201], 0x05)
			},
		},
	},

	"Fx55 - LD [I], Vx": {
		{
			0xF255,
			func(t *testing.T, c *CPU) {
				c.V[0] = 0x03
				c.V[1] = 0x02
				c.V[2] = 0x01
				c.I = 0x220
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "Memory[0x220]", c.Memory[0x220], 0x03)
				checkHex(t, "Memory[0x221]", c.Memory[0x221], 0x02)
				checkHex(t, "Memory[0x222]", c.Memory[0x222], 0x01)
			},
		},
	},

	"Fx65 - LD Vx, [I]": {
		{
			0xF165,
			func(t *testing.T, c *CPU) {
				c.Memory[0x200] = 0x01
				c.Memory[0x201] = 0x02
				c.I = 0x200
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[0]", c.V[0], 0x01)
				checkHex(t, "V[1]", c.V[1], 0x02)
				checkHex(t, "V[2]", c.V[2], 0x00)
			},
		},

		{
			0xF265,
			func(t *testing.T, c *CPU) {
				c.Memory[0x200] = 0x01
				c.Memory[0x201] = 0x02
				c.I = 0x200
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[0]", c.V[0], 0x01)
				checkHex(t, "V[1]", c.V[1], 0x02)
			},
		},
	},
}

func TestCPU_Step(t *testing.T) {
	c, _ := NewCPU(nil)
	c.Memory[0x200] = 0xA1
	c.Memory[0x201] = 0x00

	if _, err := c.Step(); err != nil {
		t.Fatal(err)
	}

	checkHex(t, "PC", c.PC, 0x202)
}

func TestOpcodes(t *testing.T) {
	for instruction, tests := range opcodeTests {
		for i, tt := range tests {
			t.Run(fmt.Sprintf("%s/%d", instruction, i), func(t *testing.T) {
				c := newCPU(t)
				if tt.before != nil {
					tt.before(t, c)
				}
				c.Dispatch(uint16(tt.op))
				tt.check(t, c)

				if t.Failed() {
					t.Logf("==============")
					t.Logf("Instruction: %s", instruction)
					t.Logf("Opcode: 0x%04X", tt.op)
					t.Logf("CPU: %v", c)
					t.Logf("==============")
					t.FailNow()
				}
			})
		}
	}
}

func TestCPU_Load(t *testing.T) {
	c := newCPU(t)
	p := []byte{0x01, 0x02}

	n, err := c.LoadBytes(p)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(p) {
		t.Fatal(err)
	}

	checkHex(t, "Memory[0x200]", c.Memory[0x200], 0x01)
	checkHex(t, "Memory[0x201]", c.Memory[0x201], 0x02)
}

func TestCPU_decodeOp(t *testing.T) {
	c := newCPU(t)
	c.Memory[0x200] = 0xA2
	c.Memory[0x201] = 0xF0

	checkHex(t, "op", c.decodeOp(), 0xA2F0)
}

func newCPU(t testing.TB) *CPU {
	c, err := NewCPU(nil)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func tryUint16(v interface{}) uint16 {
	switch v := v.(type) {
	case byte:
		return uint16(v)
	case uint16:
		return v
	case int:
		return uint16(v)
	case uint32:
		return uint16(v)
	}

	return 0
}

func checkHex(t *testing.T, subject string, got, want interface{}) {
	t.Helper()

	g := tryUint16(got)
	w := tryUint16(want)

	if g != w {
		t.Errorf("%s => 0x%04X; want 0x%04X", subject, g, w)
	}
}

func checkGraphics(t *testing.T, g *Graphics, hash string) {
	t.Helper()

	h := fmt.Sprintf("%x", sha256.Sum256(g.Pixels[:]))

	if h != hash {
		t.Errorf("Expected graphics hash to be %s, got %s", hash, h)
		t.Log("Graphics Array:")
		t.Log(g.Pixels)
	}
}
