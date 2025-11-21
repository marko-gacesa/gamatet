// Copyright (c) 2020, 2025 by Marko Gaćeša

package random

type Random struct {
	z, w uint32
}

func New(z, w uint32) *Random {
	return &Random{z: z, w: w}
}

func (r *Random) gen() uint32 {
	// Found this algorithm at: https://www.codeproject.com/Articles/25172/Simple-Random-Number-Generation.
	// Apparently it's written by https://en.wikipedia.org/wiki/George_Marsaglia.
	r.z = 36969*(r.z&65535) + r.z>>16
	r.w = 18000*(r.w&65535) + r.w>>16
	return r.z<<16 + r.w
}

func (r *Random) Int(n int) int {
	return int(r.gen() % uint32(n))
}

func (r *Random) Perm(m []int) {
	n := len(m)
	for i := range n {
		j := r.Int(i + 1)
		m[i] = m[j]
		m[j] = i
	}
}
