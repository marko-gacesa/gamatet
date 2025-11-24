// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"io"

	"github.com/marko-gacesa/gamatet/game/block"
)

type polyominoDumb struct {
	shapeRect
	block block.Block // block material for the Piece
}

var _ Piece = (*polyominoDumb)(nil)

func (p *polyominoDumb) Write(w io.Writer) error { return writeShapeRect(p.shapeRect, p.block, w) }
func (p *polyominoDumb) Read(r io.Reader) error  { return readShapeRect(&p.shapeRect, &p.block, r) }

func (p *polyominoDumb) Clone() Piece {
	clone := *p
	return &clone
}

func (p *polyominoDumb) Equals(other Piece) bool {
	q, ok := other.(*polyominoDumb)
	return ok && *p == *q
}

func (*polyominoDumb) Type() Type         { return TypeDumb }
func (p *polyominoDumb) BlockCount() byte { return p.size }
func (p *polyominoDumb) DimX() byte       { return p.width }
func (p *polyominoDumb) DimY() byte       { return p.height }

func (*polyominoDumb) CanActivate() bool     { return false }
func (*polyominoDumb) ActivationCount() byte { return 0 }

func (*polyominoDumb) Activate() bool     { return false }
func (*polyominoDumb) UndoActivate() bool { return false }

func (*polyominoDumb) WallKick() byte { return 0 }

func (p *polyominoDumb) IsEmpty(x, y int) bool { return p.data.isEmpty(p.width, p.height, x, y) }

func (p *polyominoDumb) Get(x, y int) block.Block {
	if p.IsEmpty(x, y) {
		return block.Block{Type: block.TypeEmpty}
	}
	return p.block
}

func (*polyominoDumb) LeftEmptyColumns() (empty byte)  { return 0 }
func (*polyominoDumb) RightEmptyColumns() (empty byte) { return 0 }
func (*polyominoDumb) TopEmptyRows() (empty byte)      { return 0 }
func (*polyominoDumb) BottomEmptyRows() (empty byte)   { return 0 }

func (p *polyominoDumb) String() string { return p.shapeRect.String() }
