// Copyright (c) 2020 by Marko Gaćeša

package gutil

import (
	"math"
	"testing"
)

func TestCeilPow2(t *testing.T) {
	tests := []struct {
		input int
		exp   int
	}{
		{input: 0, exp: 0},
		{input: 1, exp: 1},
		{input: 2, exp: 2},
		{input: 3, exp: 4},
		{input: 4, exp: 4},
		{input: 235, exp: 256},
		{input: 858, exp: 1024},
		{input: 32769, exp: 65536},
		{input: 1000000000, exp: 1073741824},
		{input: 2000000000, exp: 2147483648},
	}

	for _, test := range tests {
		r := CeilPow2(test.input)
		if r != test.exp {
			t.Errorf("test failed for input=%d, expected=%d, but got=%d", test.input, test.exp, r)
		}
	}
}

func TestIsPow2(t *testing.T) {
	tests := []struct {
		input int
		exp   bool
	}{
		{input: 0, exp: true},
		{input: 1, exp: true},
		{input: 2, exp: true},
		{input: 3, exp: false},
		{input: 4, exp: true},
		{input: 235, exp: false},
		{input: 858, exp: false},
		{input: 65536, exp: true},
		{input: 1073741824, exp: true},
		{input: 2147483648, exp: true},
		{input: math.MaxInt32, exp: false},
	}

	for _, test := range tests {
		r := IsPow2(test.input)
		if r != test.exp {
			t.Errorf("test failed for input=%d, expected=%t, but got=%t", test.input, test.exp, r)
		}
	}
}

func TestLog2(t *testing.T) {
	tests := []struct {
		input int
		exp   int
	}{
		{input: 0, exp: 0},
		{input: 1, exp: 0},
		{input: 2, exp: 1},
		{input: 3, exp: 1},
		{input: 4, exp: 2},
		{input: 1023, exp: 9},
		{input: 1024, exp: 10},
		{input: 1025, exp: 10},
	}

	for _, test := range tests {
		r := Log2(test.input)
		if r != test.exp {
			t.Errorf("test failed for input=%d, expected=%d, but got=%d", test.input, test.exp, r)
		}
	}
}
