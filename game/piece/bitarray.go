// Copyright (c) 2020 by Marko Gaćeša

package piece

type bitarray uint32

func (a bitarray) get(idx int) bool {
	return a&(1<<byte(idx)) != 0
}

func (a bitarray) set(idx int) bitarray {
	return a | 1<<byte(idx)
}

func (a bitarray) clear(idx int) bitarray {
	return a & ^(1 << byte(idx))
}

func (a bitarray) exchange(idx1, idx2 int) bitarray {
	mask1 := bitarray(1 << byte(idx1))
	mask2 := bitarray(1 << byte(idx2))

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
