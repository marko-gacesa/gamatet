// Copyright (c) 2020, 2025 by Marko Gaćeša

package piece

type bitarray uint32

func (a bitarray) get(idx byte) bool {
	return a&(1<<idx) != 0
}

func (a bitarray) set(idx byte) bitarray {
	return a | 1<<idx
}

func (a bitarray) clear(idx byte) bitarray {
	return a & ^(1 << idx)
}

func (a bitarray) exchange(idx1, idx2 byte) bitarray {
	mask1 := bitarray(1 << idx1)
	mask2 := bitarray(1 << idx2)

	if a&mask1 != 0 {
		if a&mask2 == 0 {
			return a & ^mask1 | mask2
		}
	} else {
		if a&mask2 != 0 {
			return a & ^mask2 | mask1
		}
	}

	return a
}

func (a bitarray) flipV(w, h byte) bitarray {
	h2 := h >> 1
	for y := range h2 {
		for x := range w {
			idx0 := y*w + x
			idx1 := (h-y-1)*w + x
			a = a.exchange(idx0, idx1)
		}
	}
	return a
}

func (a bitarray) flipH(w, h byte) bitarray {
	w2 := w >> 1
	for y := range h {
		for x := range w2 {
			idx0 := y*w + x
			idx1 := y*w + (w - x - 1)
			a = a.exchange(idx0, idx1)
		}
	}
	return a
}

func (a bitarray) rotateCW(dim byte) bitarray {
	dim1 := dim - 1
	dim2 := dim >> 1
	for j := range dim2 {
		for i := j; i < dim1-j; i++ {
			idx0 := j*dim + i
			idx1 := i*dim + dim1 - j
			idx2 := (dim1-j)*dim + dim1 - i
			idx3 := (dim1-i)*dim + j
			a = a.exchange(idx0, idx3).exchange(idx3, idx2).exchange(idx2, idx1)
		}
	}
	return a
}

func (a bitarray) rotateCCW(dim byte) bitarray {
	dim1 := dim - 1
	dim2 := dim >> 1
	for j := range dim2 {
		for i := j; i < dim1-j; i++ {
			idx0 := j*dim + i
			idx1 := i*dim + dim1 - j
			idx2 := (dim1-j)*dim + dim1 - i
			idx3 := (dim1-i)*dim + j
			a = a.exchange(idx0, idx1).exchange(idx1, idx2).exchange(idx2, idx3)
		}
	}
	return a
}

func (a bitarray) isSquareRowEmpty(dim, r byte) bool {
	lim := (r + 1) * dim
	for idx := r * dim; idx < lim; idx++ {
		if a.get(idx) {
			return false
		}
	}
	return true
}

func (a bitarray) isSquareColumnEmpty(dim, c byte) bool {
	lim := dim * dim
	for idx := c; idx < lim; idx += dim {
		if a.get(idx) {
			return false
		}
	}
	return true
}

func (a bitarray) countSquareLeftEmptyColumns(dim byte) (empty byte) {
	for i := range dim {
		if a.isSquareColumnEmpty(dim, i) {
			empty++
		} else {
			return
		}
	}
	return
}

func (a bitarray) countSquareRightEmptyColumns(dim byte) (empty byte) {
	for i := dim; i > 0; i-- {
		if a.isSquareColumnEmpty(dim, i-1) {
			empty++
		} else {
			return
		}
	}
	return
}

func (a bitarray) countSquareTopEmptyRows(dim byte) (empty byte) {
	for i := range dim {
		if a.isSquareRowEmpty(dim, i) {
			empty++
		} else {
			return
		}
	}
	return
}

func (a bitarray) countSquareBottomEmptyRows(dim byte) (empty byte) {
	for i := dim; i > 0; i-- {
		if a.isSquareRowEmpty(dim, i-1) {
			empty++
		} else {
			return
		}
	}
	return
}

func (a bitarray) isEmpty(w, h byte, x, y int) bool {
	if x < 0 || x >= int(w) || y < 0 || y >= int(h) {
		return true
	}

	return !a.get(byte(y)*w + byte(x))
}
