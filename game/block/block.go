// Copyright (c) 2020, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package block

import "fmt"

type Block struct {
	Type     Type
	Hardness byte
	Color    uint32
}

func (b Block) String() string {
	return fmt.Sprintf("[t=%d,h=%d,c=%08x]", b.Type, b.Hardness, b.Color)
}

type XY struct {
	X, Y int
}

func (xy XY) String() string {
	return fmt.Sprintf("(%d,%d)", xy.X, xy.Y)
}

type XYB struct {
	XY
	Block
}

func (xyb XYB) String() string {
	return xyb.Block.String() + "@" + xyb.XY.String()
}
