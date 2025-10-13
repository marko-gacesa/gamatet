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
