// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package random

import (
	"slices"
	"testing"
)

func TestRandomPerm(t *testing.T) {
	tests := []struct {
		z, w uint32
		n    int
	}{
		{z: 31233, w: 534526, n: 10},
		{z: 1567, w: 744, n: 20},
		{z: 763, w: 623569, n: 30},
		{z: 0, w: 0, n: 100},
	}

	for _, test := range tests {
		r := Random{z: test.z, w: test.w}

		perm := make([]int, test.n)
		r.Perm(perm)

		var sum int
		m := make(map[int]int, test.n)
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

		r2 := Random{z: test.z, w: test.w}
		perm2 := make([]int, test.n)
		r2.Perm(perm2)

		if !slices.Equal(perm, perm2) {
			t.Errorf("second perm with the same random produced different result: n=%d\n",
				test.n)
		}
	}
}

func BenchmarkPerm(b *testing.B) {
	b.ReportAllocs()
	r := Random{z: 42, w: 66}
	for i := 0; i < b.N; i++ {
		var a [100]int
		r.Perm(a[:])
	}
}
