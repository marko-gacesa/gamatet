// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"fmt"
)

type ShapeAny interface {
	_shape()
}

type shapeRect struct {
	width  byte     // Piece X dimension
	height byte     // Piece Y dimension
	size   byte     // number of blocks in the Piece
	data   bitarray // width x height matrix, 1 if has a block, 0 if not
}

func (p shapeRect) String() string {
	return p.data.rectangleString(p.width, p.height)
}

func (p shapeRect) def() string {
	return fmt.Sprintf("{width: %d, height: %d, size: %d, data: %d},", p.width, p.height, p.size, p.data)
}

func (p shapeRect) _shape() {}

type shapeSquare struct {
	dim  byte     // Piece dimension
	size byte     // number of blocks in the Piece
	rots byte     // total number of rotation positions
	data bitarray // dim x dim matrix, 1 if has a block, 0 if not
}

func (p shapeSquare) String() string {
	return p.data.rectangleString(p.dim, p.dim)
}

func (p shapeSquare) def() string {
	return fmt.Sprintf("{dim: %d, size: %d, rots: %d, data: %d},", p.dim, p.size, p.rots, p.data)
}

func (p shapeSquare) _shape() {}

type ShapeShooter struct{}

func (p ShapeShooter) _shape() {}
