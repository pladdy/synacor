package synacor

import (
	"encoding/binary"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 1, 1}}, r, 2},
		{program{index: 0, memory: []uint16{0, register0, 1, 0}}, r, 1},
		{program{index: 0, memory: []uint16{0, register0, 32766, 7}}, r, 5},
	}

	for _, test := range tests {
		add(&test.p, &test.r, &stack{})

		if test.p.index != 4 {
			t.Error("Got:", test.p.index, "Expected:", 4)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestAnd(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 1, 1}}, r, 1},
		{program{index: 0, memory: []uint16{0, register0, 1, 0}}, r, 0},
	}

	for _, test := range tests {
		and(&test.p, &test.r, &stack{})

		if test.p.index != 4 {
			t.Error("Got:", test.p.index, "Expected:", 3)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestCall(t *testing.T) {
	tests := []struct {
		p        program
		s        stack
		expected int
	}{
		{program{index: 0, memory: []uint16{0, 12}}, stack{}, 12},
		{program{index: 0, memory: []uint16{0, 15}}, stack{}, 15},
	}

	for _, test := range tests {
		call(&test.p, &registers{}, &test.s)

		if test.p.index != test.expected {
			t.Error("Got:", test.p.index, "Expected:", test.expected)
		}

		result := test.s.pop()
		if int(result) != 2 {
			t.Error("Got:", result, "Expected:", 2)
		}
	}
}

func TestEq(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 1, 1}}, r, 1},
		{program{index: 0, memory: []uint16{0, register0, 1, 0}}, r, 0},
	}

	for _, test := range tests {
		eq(&test.p, &test.r, &stack{})

		if test.p.index != 4 {
			t.Error("Got:", test.p.index, "Expected:", 4)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestGt(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 1, 1}}, r, 0},
		{program{index: 0, memory: []uint16{0, register0, 1, 0}}, r, 1},
	}

	for _, test := range tests {
		gt(&test.p, &test.r, &stack{})

		if test.p.index != 4 {
			t.Error("Got:", test.p.index, "Expected:", 3)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestInRegisters(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 1, 1}, input: []uint16{101}}, r, 101},
	}

	for _, test := range tests {
		in(&test.p, &test.r, &stack{})

		if test.p.index != 2 {
			t.Error("Got:", test.p.index, "Expected:", 2)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

// Jump will get the jump location from the next memory location (like all optimize
// functions do), but returns it decremented (since the VM loop will immediately
// loop to the next iteration).
func TestJump(t *testing.T) {
	// build a memory data set to verify register values are accessed if specified
	fullMem := []uint16{}
	for i := 0; i <= registerEnd; i++ {
		fullMem = append(fullMem, uint16(i))
	}

	tests := []struct {
		p        program
		r        registers
		expected int
	}{
		{program{index: 0, memory: []uint16{10, 11, 12}}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 11},
		{program{index: 1, memory: []uint16{10, 11, 12}}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 12},
		{program{index: 0, memory: []uint16{0, registerStart, 12}}, registers{200, 1, 2, 3, 4, 5, 6, 7}, 200},
	}

	for _, test := range tests {
		jump(&test.p, &test.r, &stack{})
		if test.p.index != test.expected {
			t.Error("Got:", test.p.index, "Expected:", test.expected)
		}
	}
}

// See TestJump docstring for why expected return index is 1 minus the index
// position.
func TestJumpFalse(t *testing.T) {
	tests := []struct {
		p        program
		r        registers
		expected int
	}{
		// i, a, b, a == 0, jump to b, return index of b
		{program{index: 0, memory: []uint16{0, 0, 3, 4, 5}}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 3},
		// i, a, b, a != 0, no jump to b, next index is 2
		{program{index: 0, memory: []uint16{0, 1, 0, 0, 0}}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 3},
	}

	for _, test := range tests {
		jumpFalse(&test.p, &test.r, &stack{})
		if test.p.index != test.expected {
			t.Error("Got:", test.p.index, "Expected:", test.expected)
		}
	}
}

// See TestJump docstring for why expected return index is 1 minus the index
// position.
func TestJumpTrue(t *testing.T) {
	tests := []struct {
		p        program
		r        registers
		expected int
	}{
		// i, a, b, a >= 0, jump to b, return index of b
		{program{index: 0, memory: []uint16{0, 1, 3, 4, 5}}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 3},
		// i, a, b, a !> 0, no jump to b -> 2
		{program{index: 0, memory: []uint16{0, 0, 0, 0, 0}}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 3},
	}

	for _, test := range tests {
		jumpTrue(&test.p, &test.r, &stack{})
		if test.p.index != test.expected {
			t.Error("Got:", test.p.index, "Expected:", test.expected)
		}
	}
}

func TestMod(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 1, 1}}, r, 0},
		{program{index: 0, memory: []uint16{0, register0, 1, 2}}, r, 1},
		{program{index: 0, memory: []uint16{0, register0, 32766, 7}}, r, 6},
	}

	for _, test := range tests {
		mod(&test.p, &test.r, &stack{})

		if test.p.index != 4 {
			t.Error("Got:", test.p.index, "Expected:", 4)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestMult(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 1, 1}}, r, 1},
		{program{index: 0, memory: []uint16{0, register0, 1, 0}}, r, 0},
		{program{index: 0, memory: []uint16{0, register0, 4, 9}}, r, 36},
		{program{index: 0, memory: []uint16{0, register0, 32766, 7}}, r, 32754},
	}

	for _, test := range tests {
		mult(&test.p, &test.r, &stack{})

		if test.p.index != 4 {
			t.Error("Got:", test.p.index, "Expected:", 4)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestNoop(t *testing.T) {
	tests := []struct {
		p        program
		r        registers
		expected int
	}{
		{program{index: 0, memory: []uint16{0, 1, 2}}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 1},
	}

	for _, test := range tests {
		noop(&test.p, &test.r, &stack{})
		if test.p.index != test.expected {
			t.Error("Got:", test.p.index, "Expected:", test.expected)
		}
	}
}

func TestNot(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 2, 1}}, r, 32765},
		{program{index: 0, memory: []uint16{0, register0, 1, 0}}, r, 32766},
	}

	for _, test := range tests {
		not(&test.p, &test.r, &stack{})

		if test.p.index != 3 {
			t.Error("Got:", test.p.index, "Expected:", 3)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestOr(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 2, 1}}, r, 3},
		{program{index: 0, memory: []uint16{0, register0, 1, 0}}, r, 1},
	}

	for _, test := range tests {
		or(&test.p, &test.r, &stack{})

		if test.p.index != 4 {
			t.Error("Got:", test.p.index, "Expected:", 4)
		}

		result := test.r.get(register0)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestOut(t *testing.T) {
	tests := []struct {
		p        program
		r        registers
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, 65}}, registers{0, 1, 2, 3, 4, 5, 6, 7}, 65},
	}

	for _, test := range tests {
		// Redirect stdout via Pipe()
		r, w, _ := os.Pipe()
		backupStdout := os.Stdout
		os.Stdout = w

		out(&test.p, &test.r, &stack{})

		buf := make([]byte, 2)
		_, err := r.Read(buf)
		if err != nil {
			t.Error(err)
		}

		// Convert buffer and restore stdout
		result := binary.LittleEndian.Uint16(buf)
		w.Close()
		os.Stdout = backupStdout

		if uint16(result) != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		p        program
		s        stack
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, 14}}, stack{}, 14},
	}

	for _, test := range tests {
		push(&test.p, &registers{}, &test.s)
		result := test.s[0]

		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestPop(t *testing.T) {
	tests := []struct {
		p        program
		r        registers
		s        stack
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0}}, registers{}, stack{14}, 14},
	}

	for _, test := range tests {
		pop(&test.p, &test.r, &test.s)
		result := test.r.get(register0)

		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestRet(t *testing.T) {
	tests := []struct {
		p        program
		s        stack
		expected int
	}{
		{program{index: 0, memory: []uint16{0, 12}}, stack{27}, 27},
		{program{index: 0, memory: []uint16{0, 15}}, stack{14}, 14},
	}

	for _, test := range tests {
		ret(&test.p, &registers{}, &test.s)

		if test.p.index != test.expected {
			t.Error("Got:", test.p.index, "Expected:", test.expected)
		}
	}
}

func TestRmem(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		s        stack
		expected int
	}{
		{program{index: 0, memory: []uint16{0, 2, 2, 3}}, r, stack{27}, 0},
		{program{index: 0, memory: []uint16{0, 2, 2, 4}}, r, stack{14}, 0},
	}

	for _, test := range tests {
		register := test.p.memory[1]
		rmem(&test.p, &registers{}, &test.s)

		if test.p.index != 3 {
			t.Error("Got:", test.p.index, "Expected:", 3)
		}

		result := test.r.get(register)
		if int(result) != test.expected {
			t.Error("Got:", result, "Expected:", test.expected, "Register:", register)
		}
	}
}

func TestSet(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		s        stack
		expected uint16
	}{
		{program{index: 0, memory: []uint16{0, register0, 42}}, r, stack{}, 42},
	}

	for _, test := range tests {
		register := test.p.memory[1]
		set(&test.p, &test.r, &stack{})

		if test.p.index != 3 {
			t.Error("Got:", test.p.index, "Expected:", 3)
		}

		result := test.r.get(register)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected, "Register:", register)
		}
	}
}

func TestWmem(t *testing.T) {
	r := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		p        program
		r        registers
		s        stack
		expected int
	}{
		{program{index: 0, memory: []uint16{0, 2, 2, 3}}, r, stack{27}, 0},
		{program{index: 0, memory: []uint16{0, 2, 2, 4}}, r, stack{14}, 0},
	}

	for _, test := range tests {
		register := test.p.memory[1]
		wmem(&test.p, &registers{}, &test.s)

		if test.p.index != 3 {
			t.Error("Got:", test.p.index, "Expected:", 3)
		}

		result := test.r.get(register)
		if int(result) != test.expected {
			t.Error("Got:", result, "Expected:", test.expected, "Register:", register)
		}
	}
}
