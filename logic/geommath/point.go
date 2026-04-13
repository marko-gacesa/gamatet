// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package geommath

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type P[T Integer] struct {
	X, Y T
}

func Abs[T Integer](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

func Manhattan[T Integer](p1, p2 P[T]) T {
	return Abs[T](p1.X-p2.X) + Abs[T](p1.Y-p2.Y)
}
