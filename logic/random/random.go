// Copyright (c) 2020, 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package random

import "math/rand/v2"

type Random struct {
	src rand.PCG
}

func New(seed1, seed2 uint64) *Random {
	src := rand.NewPCG(seed1, seed2)
	return &Random{
		src: *src,
	}
}

func (r *Random) Int(n int) int {
	return int(uint(r.src.Uint64()) % uint(n))
}

func (r *Random) UInt(n uint) uint {
	return uint(r.src.Uint64()) % n
}

func (r *Random) Perm(m []uint) {
	n := uint(len(m))
	for i := range n {
		j := r.UInt(i + 1)
		m[i] = m[j]
		m[j] = i
	}
}

func Shuffle[T any](r *Random, a []T) {
	if len(a) == 0 {
		return
	}
	for i := uint(len(a) - 1); i > 0; i-- {
		j := r.UInt(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}
