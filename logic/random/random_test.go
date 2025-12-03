// Copyright (c) 2020, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package random

import (
	"slices"
	"testing"
)

func TestRandomPerm(t *testing.T) {
	tests := []struct {
		seed1, seed2 uint
		n            uint
	}{
		{seed1: 31233, seed2: 534526, n: 10},
		{seed1: 1567, seed2: 744, n: 20},
		{seed1: 763, seed2: 623569, n: 30},
		{seed1: 0, seed2: 0, n: 100},
	}

	for _, test := range tests {
		r := New(test.seed1, test.seed2)

		perm := make([]uint, test.n)
		r.Perm(perm)

		var sum uint
		m := make(map[uint]uint, test.n)
		for idx, num := range perm {
			if num < 0 || num >= test.n {
				t.Errorf("number outside of range: n=%d idx=%d num=%d\n",
					test.n, idx, num)
			}
			if m[num] > 0 {
				t.Errorf("number not unique in the permutation: n=%d idx=%d num=%d\n",
					test.n, idx, num)
			}

			sum += num
			m[num]++
		}
		if sum != test.n*(test.n-1)/2 {
			t.Errorf("sum does not match: n=%d, expected=%d, got=%d\n",
				test.n, test.n*(test.n-1)/2, sum)
		}

		r2 := New(test.seed1, test.seed2)
		perm2 := make([]uint, test.n)
		r2.Perm(perm2)

		if !slices.Equal(perm, perm2) {
			t.Errorf("second perm with the same random produced different result: n=%d\n",
				test.n)
		}
	}
}
