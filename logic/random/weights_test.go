// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package random

import (
	"testing"
)

func TestWeights(t *testing.T) {
	tests := []struct {
		name     string
		weights  []uint
		expected []uint
	}{
		{
			name:     "empty",
			weights:  []uint{},
			expected: []uint{},
		},
		{
			name:     "zeros",
			weights:  []uint{0, 0, 0},
			expected: []uint{},
		},
		{
			name:     "even",
			weights:  []uint{2, 2, 2},
			expected: []uint{0, 0, 1, 1, 2, 2},
		},
		{
			name:     "variable",
			weights:  []uint{1, 3, 2, 4},
			expected: []uint{0, 1, 1, 1, 2, 2, 3, 3, 3, 3},
		},
		{
			name:     "with_zero",
			weights:  []uint{2, 0, 1},
			expected: []uint{0, 0, 2},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var w Weights[uint]
			for _, k := range test.weights {
				w.Add(k)
			}

			n := w.Total()

			for i := uint(0); i < n; i++ {
				if want, got := test.expected[i], w.Weighted(i); want != got {
					t.Errorf("idx %d: want: %d, got: %d", i, want, got)
				}
			}
		})
	}
}
