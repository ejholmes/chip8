package chip8

import "testing"

func TestCPU_Step(t *testing.T) {
	c := NewCPU(nil)
	c.Memory[200] = 0xA1
	c.Memory[201] = 0x00

	if err := c.Step(); err != nil {
		t.Fatal(err)
	}

	if got, want := c.PC, uint16(202); got != want {
		t.Errorf("PC => %d; want %d", got, want)
	}
}

func TestCPU_Dispatch(t *testing.T) {
	tests := []struct {
		op    uint16
		check func(*CPU)
	}{
		{
			uint16(0xA100),
			func(c *CPU) {
				if got, want := c.I, uint16(0x100); got != want {
					t.Errorf("I => %x; want %x", got, want)
				}
			},
		},
	}

	for _, tt := range tests {
		c := NewCPU(nil)
		c.Dispatch(tt.op)
		tt.check(c)
	}
}

func TestCPU_op(t *testing.T) {
	c := NewCPU(nil)
	c.Memory[200] = 0xA2
	c.Memory[201] = 0xF0

	if got, want := c.op(), uint16(0xA2F0); got != want {
		t.Errorf("op => %x; want %x", got, want)
	}
}
