// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"errors"
	"io"
	"math"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/serialize"
)

type polyominoFlip struct {
	shapeRect
	block block.Block // block material for the Piece
}

type polyominoFlipV polyominoFlip

var _ Piece = (*polyominoFlipV)(nil)

func (p *polyominoFlipV) Write(w io.Writer) error { return writeShapeRect(p.shapeRect, p.block, w) }
func (p *polyominoFlipV) Read(r io.Reader) error  { return readShapeRect(&p.shapeRect, &p.block, r) }

func (p *polyominoFlipV) Clone() Piece {
	clone := *p
	return &clone
}

func (p *polyominoFlipV) Equals(other Piece) bool {
	q, ok := other.(*polyominoFlipV)
	return ok && *p == *q
}

func (*polyominoFlipV) Type() Type         { return TypeFlipV }
func (p *polyominoFlipV) BlockCount() byte { return p.size }
func (p *polyominoFlipV) DimX() byte       { return p.width }
func (p *polyominoFlipV) DimY() byte       { return p.height }

func (*polyominoFlipV) CanActivate() bool     { return true }
func (*polyominoFlipV) ActivationCount() byte { return math.MaxUint8 }

func (p *polyominoFlipV) Activate() bool {
	p.shapeRect.data = p.shapeRect.data.flipV(p.width, p.height)
	return false
}

func (p *polyominoFlipV) UndoActivate() bool {
	p.shapeRect.data = p.shapeRect.data.flipV(p.width, p.height)
	return false
}

func (p *polyominoFlipV) WallKick() byte { return 0 }

func (p *polyominoFlipV) IsEmpty(x, y int) bool { return p.data.isEmpty(p.width, p.height, x, y) }

func (p *polyominoFlipV) Get(x, y int) block.Block {
	if p.IsEmpty(x, y) {
		return block.Block{Type: block.TypeEmpty}
	}
	return p.block
}

func (p *polyominoFlipV) LeftEmptyColumns() (empty byte)  { return 0 }
func (p *polyominoFlipV) RightEmptyColumns() (empty byte) { return 0 }
func (p *polyominoFlipV) TopEmptyRows() (empty byte)      { return 0 }
func (p *polyominoFlipV) BottomEmptyRows() (empty byte)   { return 0 }

func (p *polyominoFlipV) String() string { return p.shapeRect.String() }

type polyominoFlipH polyominoFlip

var _ Piece = (*polyominoFlipH)(nil)

func (p *polyominoFlipH) Write(w io.Writer) error { return writeShapeRect(p.shapeRect, p.block, w) }
func (p *polyominoFlipH) Read(r io.Reader) error  { return readShapeRect(&p.shapeRect, &p.block, r) }

func (p *polyominoFlipH) Clone() Piece {
	clone := *p
	return &clone
}

func (p *polyominoFlipH) Equals(other Piece) bool {
	q, ok := other.(*polyominoFlipH)
	return ok && *p == *q
}

func (*polyominoFlipH) Type() Type         { return TypeFlipH }
func (p *polyominoFlipH) BlockCount() byte { return p.size }
func (p *polyominoFlipH) DimX() byte       { return p.width }
func (p *polyominoFlipH) DimY() byte       { return p.height }

func (*polyominoFlipH) CanActivate() bool     { return true }
func (*polyominoFlipH) ActivationCount() byte { return math.MaxUint8 }

func (p *polyominoFlipH) Activate() bool {
	p.shapeRect.data = p.shapeRect.data.flipH(p.width, p.height)
	return false
}

func (p *polyominoFlipH) UndoActivate() bool {
	p.shapeRect.data = p.shapeRect.data.flipH(p.width, p.height)
	return false
}

func (p *polyominoFlipH) WallKick() byte { return 0 }

func (p *polyominoFlipH) IsEmpty(x, y int) bool { return p.data.isEmpty(p.width, p.height, x, y) }

func (p *polyominoFlipH) Get(x, y int) block.Block {
	if p.IsEmpty(x, y) {
		return block.Block{Type: block.TypeEmpty}
	}
	return p.block
}

func (p *polyominoFlipH) LeftEmptyColumns() (empty byte)  { return 0 }
func (p *polyominoFlipH) RightEmptyColumns() (empty byte) { return 0 }
func (p *polyominoFlipH) TopEmptyRows() (empty byte)      { return 0 }
func (p *polyominoFlipH) BottomEmptyRows() (empty byte)   { return 0 }

func (p *polyominoFlipH) String() string { return p.shapeRect.String() }

func writeShapeRect(s shapeRect, b block.Block, w io.Writer) error {
	var buffer [3]byte
	buffer[0] = s.width
	buffer[1] = s.height
	buffer[2] = s.size

	if _, err := w.Write(buffer[:]); err != nil {
		return err
	}

	if err := b.Write(w); err != nil {
		return err
	}

	if err := serialize.Write32(w, uint32(s.data)); err != nil {
		return err
	}

	return nil
}

func readShapeRect(s *shapeRect, b *block.Block, r io.Reader) (err error) {
	var buffer [3]byte

	n, err := r.Read(buffer[:])
	if err != nil {
		return
	}

	if n != 3 {
		err = errors.New("failed to read polyomino rectangle")
		return
	}

	s.width = buffer[0]
	s.height = buffer[1]
	s.size = buffer[2]

	err = b.Read(r)
	if err != nil {
		return
	}

	data, err := serialize.Read32(r)
	if err != nil {
		return
	}

	s.data = bitarray(data)

	return
}
