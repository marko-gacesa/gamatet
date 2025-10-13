// Copyright (c) 2020, 2025 by Marko Gaćeša

package piece

import "testing"

func TestBitArrayGet(t *testing.T) {
	tests := []struct {
		value    uint32
		bit      byte
		expected bool
	}{
		{value: 0b1001, bit: 0, expected: true},
		{value: 0b1001, bit: 1, expected: false},
		{value: 0b1001, bit: 2, expected: false},
		{value: 0b1001, bit: 3, expected: true},
		{value: 0x80008000, bit: 31, expected: true},
		{value: 0x80008000, bit: 30, expected: false},
		{value: 0x7F00F000, bit: 31, expected: false},
		{value: 0x7F00F000, bit: 30, expected: true},
	}

	for _, test := range tests {
		x := bitarray(test.value)
		got := x.get(test.bit)
		if test.expected != got {
			t.Errorf("test failed for %b bit=%d, expected=%t got=%t",
				test.value, test.bit, test.expected, got)
		}
	}
}

func TestBitArraySet(t *testing.T) {
	tests := []struct {
		initial  uint32
		bit      byte
		expected uint32
	}{
		{initial: 0, bit: 0, expected: 1},
		{initial: 0, bit: 1, expected: 0b10},
		{initial: 0b010101, bit: 3, expected: 0b011101},
		{initial: 0b010101, bit: 2, expected: 0b010101},
		{initial: 0, bit: 31, expected: 0x80000000},
	}

	for _, test := range tests {
		x := bitarray(test.initial)
		got := uint32(x.set(test.bit))
		if test.expected != got {
			t.Errorf("test failed for %b bit=%d, expected=%b got=%b",
				test.initial, test.bit, test.expected, got)
		}
	}
}

func TestBitArrayClear(t *testing.T) {
	tests := []struct {
		initial  uint32
		bit      byte
		expected uint32
	}{
		{initial: 0xFFFFFFFF, bit: 0, expected: 0xFFFFFFFE},
		{initial: 0xFFFFFFFF, bit: 31, expected: 0x7FFFFFFF},
		{initial: 0b010101, bit: 2, expected: 0b010001},
		{initial: 0b010101, bit: 3, expected: 0b010101},
	}

	for _, test := range tests {
		x := bitarray(test.initial)
		got := uint32(x.clear(test.bit))
		if test.expected != got {
			t.Errorf("test failed for %b bit=%d, expected=%b got=%b",
				test.initial, test.bit, test.expected, got)
		}
	}
}

func TestBitArrayExchange(t *testing.T) {
	tests := []struct {
		initial    uint32
		bit1, bit2 byte
		expected   uint32
	}{
		{initial: 0b10, bit1: 0, bit2: 1, expected: 0b01},
		{initial: 0b10, bit1: 1, bit2: 0, expected: 0b01},
		{initial: 0x7FFFFFFF, bit1: 31, bit2: 0, expected: 0xFFFFFFFE},
	}

	for _, test := range tests {
		x := bitarray(test.initial)
		got := uint32(x.exchange(test.bit1, test.bit2))
		if test.expected != got {
			t.Errorf("test failed for %b bit1=%d bit2=%d, expected=%b got=%b",
				test.initial, test.bit1, test.bit2, test.expected, got)
		}
	}
}
