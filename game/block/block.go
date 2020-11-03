// Copyright (c) 2020 by Marko Gaćeša

package block

type Block struct {
	Type     Type
	Hardness byte
	Color    uint32
}

type XY struct {
	X, Y int
}

type XYB struct {
	XY
	Block
}

func SliceEqual(a, b []Block) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
