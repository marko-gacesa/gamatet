// Copyright (c) 2020, 2025 by Marko Gaćeša

package piece

import (
	"errors"
	"io"
	"math"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/serialize"
)

type polyominoRot struct {
	shapeSquare
	rot   byte        // current rotation position
	block block.Block // block material for the Piece
}

var _ Piece = (*polyominoRot)(nil)

func (p *polyominoRot) Write(w io.Writer) error {
	var buffer [4]byte
	buffer[0] = p.dim
	buffer[1] = p.size
	buffer[2] = p.rot
	buffer[3] = p.rots

	if _, err := w.Write(buffer[:]); err != nil {
		return err
	}

	if err := p.block.Write(w); err != nil {
		return err
	}

	if err := serialize.Write32(w, uint32(p.data)); err != nil {
		return err
	}

	return nil
}

func (p *polyominoRot) Read(r io.Reader) (err error) {
	var buffer [4]byte

	n, err := r.Read(buffer[:])
	if err != nil {
		return
	}

	if n != 4 {
		err = errors.New("failed to read polyomino square")
		return
	}

	p.dim = buffer[0]
	p.size = buffer[1]
	p.rot = buffer[2]
	p.rots = buffer[3]

	err = p.block.Read(r)
	if err != nil {
		return
	}

	data, err := serialize.Read32(r)
	if err != nil {
		return
	}

	p.data = bitarray(data)

	return
}

func (p *polyominoRot) Clone() Piece {
	clone := *p
	return &clone
}

func (p *polyominoRot) Equals(other Piece) bool {
	q, ok := other.(*polyominoRot)
	return ok && *p == *q
}

func (*polyominoRot) Type() Type         { return TypeRotation }
func (p *polyominoRot) BlockCount() byte { return p.size }
func (p *polyominoRot) DimX() byte       { return p.dim }
func (p *polyominoRot) DimY() byte       { return p.dim }

func (*polyominoRot) CanActivate() bool     { return true }
func (*polyominoRot) ActivationCount() byte { return math.MaxUint8 }

// Activate is rotate CCW
func (p *polyominoRot) Activate() (inverted bool) {
	switch p.rots {
	case 2:
		if p.rot == 0 {
			p.rot++
			p.data = p.data.rotateCCW(p.dim)
		} else {
			p.rot = 0
			p.data = p.data.rotateCW(p.dim)
			inverted = true
		}
	case 4:
		p.rot = (p.rot - 1) & 3
		p.data = p.data.rotateCCW(p.dim)
	}
	return
}

// UndoActivate is rotate CW
func (p *polyominoRot) UndoActivate() (inverted bool) {
	switch p.rots {
	case 2:
		if p.rot != 0 {
			p.rot = 0
			p.data = p.data.rotateCW(p.dim)
		} else {
			p.rot++
			p.data = p.data.rotateCCW(p.dim)
			inverted = true
		}
	case 4:
		p.rot = (p.rot + 1) & 3
		p.data = p.data.rotateCW(p.dim)
	}
	return
}

func (p *polyominoRot) WallKick() byte { return p.dim / 2 }

func (p *polyominoRot) IsEmpty(x, y int) bool { return p.data.isEmpty(p.dim, p.dim, x, y) }

func (p *polyominoRot) Get(x, y int) block.Block {
	if p.IsEmpty(x, y) {
		return block.Block{Type: block.TypeEmpty}
	}
	return p.block
}

func (p *polyominoRot) LeftEmptyColumns() byte  { return p.data.countSquareLeftEmptyColumns(p.dim) }
func (p *polyominoRot) RightEmptyColumns() byte { return p.data.countSquareRightEmptyColumns(p.dim) }
func (p *polyominoRot) TopEmptyRows() byte      { return p.data.countSquareTopEmptyRows(p.dim) }
func (p *polyominoRot) BottomEmptyRows() byte   { return p.data.countSquareBottomEmptyRows(p.dim) }

func (p *polyominoRot) String() string { return p.shapeSquare.String() }
