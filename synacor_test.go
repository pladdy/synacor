package main

import (
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
