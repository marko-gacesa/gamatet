// Copyright (c) 2025 by Marko Gaćeša

package piece

import (
	"fmt"
	"strings"
)

type shapeRect struct {
	width  byte     // Piece X dimension
	height byte     // Piece Y dimension
	size   byte     // number of blocks in the Piece
	data   bitarray // width x height matrix, 1 if has a block, 0 if not
}

func (p shapeRect) String() string {
	return polyominoToStr(p.data, p.width, p.height)
}

func (p shapeRect) def() string {
	return fmt.Sprintf("{width: %d, height: %d, size: %d, data: %d},", p.width, p.height, p.size, p.data)
}

type shapeSquare struct {
	dim  byte     // Piece dimension
	size byte     // number of blocks in the Piece
	rots byte     // total number of rotation positions
	data bitarray // dim x dim matrix, 1 if has a block, 0 if not
}

func (p shapeSquare) String() string {
	return polyominoToStr(p.data, p.dim, p.dim)
}

func (p shapeSquare) def() string {
	return fmt.Sprintf("{dim: %d, size: %d, rots: %d, data: %d},", p.dim, p.size, p.rots, p.data)
}

func polyominoToStr(b bitarray, w, h byte) string {
	sb := strings.Builder{}

	for j := range h {
		for i := range w {
			if b.get(j*w + i) {
				sb.WriteString("[]")
			} else {
				sb.WriteString(". ")
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
