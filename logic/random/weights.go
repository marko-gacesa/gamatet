// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package random

type Weights[T integer] []T

func (w *Weights[T]) Add(amount T) {
	*w = append(*w, amount)
}

func (w *Weights[T]) Total() T {
	var sum T
	for _, amount := range *w {
		sum += amount
	}

	return sum
}

func (w *Weights[T]) Weighted(n T) T {
	sum := w.Total()

	for i := T(len(*w)) - 1; i >= 0; i-- {
		sum -= (*w)[i]
		if n >= sum {
			return i
		}
	}

	panic("unreachable")
}

func (w *Weights[T]) Random(r *Random) T {
	return w.Weighted(T(r.UInt(uint(w.Total()))))
}

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
