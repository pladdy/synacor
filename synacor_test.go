package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		value    uint16
		expected bool
	}{
		{0, true},
		{32775, true},
		{32776, false},
	}

	for _, test := range tests {
		result := isValid(test.value)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected, "Value:", test.value)
		}
	}
}

func TestIsLiteralValue(t *testing.T) {
	tests := []struct {
		value    uint16
		expected bool
	}{
		{0, true},
		{32767, true},
		{32768, false},
		{32769, false},
	}

	for _, test := range tests {
		result := isLiteralValue(test.value)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected, "Value:", test.value)
		}
	}
}

func TestIsRegister(t *testing.T) {
	tests := []struct {
		value    uint16
		expected bool
	}{
		{32768, true},
		{32775, true},
		{32776, false},
	}

	for _, test := range tests {
		result := isRegister(test.value)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected, "Value:", test.value)
		}
	}
}

func TestOpCodes(t *testing.T) {
	if opHalt != 0 {
		t.Error("Opcode halt should be 0")
	}
	if opOut != 19 {
		t.Error("Opcode out should be 19")
	}
	if opNoop != 21 {
		t.Error("Opcode noop should be 21")
	}
}

func TestProgramGetNext(t *testing.T) {
	tests := []struct {
		p             program
		reg           registers
		expectedIndex int
		expectedValue uint16
	}{
		{program{index: 0, memory: []uint16{0, 1}}, registers{}, 1, 1},
		{program{index: 1, memory: []uint16{0, 1, 2}}, registers{}, 2, 2},
		{program{index: 0, memory: []uint16{0, register0}}, registers{42}, 1, 42},
	}

	for _, test := range tests {
		result := test.p.getNext(&test.reg)

		if test.p.index != test.expectedIndex {
			t.Error("Got:", test.p.index, "Expected:", test.expectedIndex)
		}
		if result != test.expectedValue {
			t.Error("Got:", result, "Expected:", test.expectedValue)
		}
	}
}

func TestProgramGetNextRaw(t *testing.T) {
	tests := []struct {
		p             program
		expectedIndex int
		expectedValue uint16
	}{
		{program{index: 0, memory: []uint16{0, 1}}, 1, 1},
		{program{index: 1, memory: []uint16{0, 1, 2}}, 2, 2},
		{program{index: 0, memory: []uint16{0, register0}}, 1, register0},
	}

	for _, test := range tests {
		result := test.p.getNextRaw()

		if test.p.index != test.expectedIndex {
			t.Error("Got:", test.p.index, "Expected:", test.expectedIndex)
		}
		if result != test.expectedValue {
			t.Error("Got:", result, "Expected:", test.expectedValue)
		}
	}
}

func TestProgramLoad(t *testing.T) {
	file, err := os.Create("test.bin")
	if err != nil {
		t.Error("Failed to create test file:", err)
	}

	err = binary.Write(file, binary.LittleEndian, uint16(19))
	if err != nil {
		t.Error("Failed to write to test file", err)
	}

	err = file.Close()
	if err != nil {
		t.Error("Failed to close file", err)
	}

	p := program{}
	p.load("test.bin")
	memoryLen := len(p.memory)
	expected := 1

	if len(p.memory) != expected {
		t.Error("Got:", memoryLen, "Expected:", expected)
	}
}

func TestReadNext(t *testing.T) {
	tests := []struct {
		testInt uint16
		err     error
	}{
		{testInt: 18097, err: nil},
		{testInt: 0, err: nil},
	}

	for _, test := range tests {
		binInt := make([]byte, 2)
		binary.LittleEndian.PutUint16(binInt, test.testInt)

		le, err := readNext(bytes.NewReader(binInt))
		if le != test.testInt {
			t.Error("Got:", le, "Expected:", test.testInt)
		}
		if err != test.err {
			t.Error("Got:", err, "Expected:", test.err)
		}
	}
}

func TestReadNextErr(t *testing.T) {
	binInt := make([]byte, 2)
	binary.LittleEndian.PutUint16(binInt, 1)

	reader := bytes.NewReader(binInt)
	_, err := readNext(reader)
	if err != nil {
		t.Error("Got:", err, "Expected:", nil)
	}

	_, err = readNext(reader)
	if err != io.EOF {
		t.Error("Got:", err, "Expected:", io.EOF)
	}
}

func TestRegisterGet(t *testing.T) {
	tests := []struct {
		reg      registers
		index    uint16
		expected uint16
	}{
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, 0, 0},
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, register0, 0},
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, register1, 1},
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, register2, 2},
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, register3, 3},
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, register4, 4},
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, register5, 5},
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, register6, 6},
		{registers{0, 1, 2, 3, 4, 5, 6, 7}, register7, 7},
	}

	for _, test := range tests {
		result := test.reg.get(test.index)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestRegisterSet(t *testing.T) {
	reg := registers{0, 0, 0, 0, 0, 0, 0, 0}
	tests := []struct {
		reg      registers
		index    uint16
		expected uint16
	}{
		{reg, 0, 42},
		{reg, register0, 13},
		{reg, register1, 12},
		{reg, register2, 29},
		{reg, register3, 101},
		{reg, register4, 4},
		{reg, register5, 99},
		{reg, register6, 5},
		{reg, register7, 1},
	}

	for _, test := range tests {
		test.reg.set(test.index, test.expected)
		result := test.reg.get(test.index)
		if result != test.expected {
			t.Error("Got:", result, "Expected:", test.expected)
		}
	}
}

func TestStackIsEmpty(t *testing.T) {
	tests := []struct {
		s        stack
		expected bool
	}{
		{stack{0}, false},
		{stack{}, true},
	}

	for _, test := range tests {
		if test.s.isEmpty() != test.expected {
			t.Error("Got:", test.s.isEmpty(), "Expected:", test.expected)
		}
	}
}

func TestStackPop(t *testing.T) {
	tests := []struct {
		s        stack
		expected uint16
	}{
		{stack{12}, 12},
		{stack{}, 0},
	}

	for _, test := range tests {
		if test.s.pop() != test.expected {
			t.Error("Got:", test.s.pop(), "Expected:", test.expected)
		}
	}
}

func TestStackPush(t *testing.T) {
	tests := []struct {
		oldStack stack
		value    uint16
		expected uint16
	}{
		{stack{12}, 12, 12},
		{stack{0}, 0, 0},
	}

	for _, test := range tests {
		test.oldStack.push(test.value)
		if test.oldStack[1] != test.expected {
			t.Error("Got:", test.oldStack[1], "Expected:", test.expected)
		}
	}
}
