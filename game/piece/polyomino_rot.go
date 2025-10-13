// Copyright (c) 2020, 2025 by Marko Gaćeša

package piece

import (
	"errors"
	"gamatet/game/block"
	"gamatet/game/serialize"
	"io"
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

func (*polyominoRot) Type() Type         { return TypeStandard }
func (p *polyominoRot) BlockCount() byte { return p.size }
func (p *polyominoRot) DimX() byte       { return p.dim }
func (p *polyominoRot) DimY() byte       { return p.dim }

func (*polyominoRot) CanActivate() bool        { return false }
func (*polyominoRot) GetActivationCount() byte { return 0 }
func (*polyominoRot) SetActivationCount(byte)  {}

func (p *polyominoRot) FlipV() {
	dim := p.dim
	y1 := p.TopEmptyRows()
	y2 := dim - p.BottomEmptyRows()
	flipDim := y2 - y1
	flipDim2 := flipDim >> 1
	for i := byte(0); i < flipDim2; i++ {
		for x := byte(0); x < dim; x++ {
			idx0 := (y1+i)*dim + x
			idx1 := (flipDim+y1-i-1)*dim + x
			p.data = p.data.exchange(idx0, idx1)
		}
	}
}

func (p *polyominoRot) FlipH() {
	dim := p.dim
	x1 := p.LeftEmptyColumns()
	x2 := dim - p.RightEmptyColumns()
	flipDim := x2 - x1
	flipDim2 := flipDim >> 1
	for y := byte(0); y < dim; y++ {
		for i := byte(0); i < flipDim2; i++ {
			idx0 := y*dim + x1 + i
			idx1 := y*dim + (x1 + flipDim - i - 1)
			p.data = p.data.exchange(idx0, idx1)
		}
	}
}

func (p *polyominoRot) RotateCW() (inverted bool) {
	switch p.rots {
	default:
		return
	case 2:
		if p.rot != 0 {
			p.rot--
		} else {
			p.RotateCCW()
			inverted = true
			return
		}
	case 4:
		p.rot = (p.rot + 1) & 3
	}

	dim := p.dim
	dim1 := dim - 1
	dim2 := p.dim >> 1
	for j := byte(0); j < dim2; j++ {
		for i := j; i < dim1-j; i++ {
			idx0 := j*dim + i
			idx1 := i*dim + dim1 - j
			idx2 := (dim1-j)*dim + dim1 - i
			idx3 := (dim1-i)*dim + j
			p.data = p.data.exchange(idx0, idx3).exchange(idx3, idx2).exchange(idx2, idx1)
		}
	}

	return
}

func (p *polyominoRot) RotateCCW() (inverted bool) {
	switch p.rots {
	default:
		return
	case 2:
		if p.rot == 0 {
			p.rot++
		} else {
			p.RotateCW()
			inverted = true
			return
		}
		//p.rot = (p.rot + 1) & 1
	case 4:
		p.rot = (p.rot - 1) & 3
	}

	dim := p.dim
	dim1 := dim - 1
	dim2 := p.dim >> 1
	for j := byte(0); j < dim2; j++ {
		for i := j; i < dim1-j; i++ {
			idx0 := j*dim + i
			idx1 := i*dim + dim1 - j
			idx2 := (dim1-j)*dim + dim1 - i
			idx3 := (dim1-i)*dim + j
			p.data = p.data.exchange(idx0, idx1).exchange(idx1, idx2).exchange(idx2, idx3)
		}
	}

	return
}

func (p *polyominoRot) WallKick() byte {
	return p.dim / 2
}

func (p *polyominoRot) IsEmpty(x, y int) bool {
	d := int(p.dim)
	if x < 0 || x >= d || y < 0 || y >= d {
		return true
	}

	idx := y*d + x
	return !p.data.get(byte(idx))
}

func (p *polyominoRot) Get(x, y int) (b block.Block) {
	if p.IsEmpty(x, y) {
		return
	}
	b = p.block
	return
}

func (p *polyominoRot) isRowEmpty(r byte) bool {
	d := p.dim
	lim := (r + 1) * d
	for idx := r * d; idx < lim; idx++ {
		if p.data.get(idx) {
			return false
		}
	}
	return true
}

func (p *polyominoRot) isColumnEmpty(c byte) bool {
	d := p.dim
	lim := d * d
	for idx := c; idx < lim; idx += d {
		if p.data.get(idx) {
			return false
		}
	}
	return true
}

func (p *polyominoRot) LeftEmptyColumns() (empty byte) {
	d := p.dim
	for i := byte(0); i < d; i++ {
		if p.isColumnEmpty(i) {
			empty++
		} else {
			return
		}
	}
	return
}

func (p *polyominoRot) RightEmptyColumns() (empty byte) {
	for i := p.dim; i > 0; i-- {
		if p.isColumnEmpty(i - 1) {
			empty++
		} else {
			return
		}
	}
	return
}

func (p *polyominoRot) TopEmptyRows() (empty byte) {
	d := p.dim
	for i := byte(0); i < d; i++ {
		if p.isRowEmpty(i) {
			empty++
		} else {
			return
		}
	}
	return
}

func (p *polyominoRot) BottomEmptyRows() (empty byte) {
	for i := p.dim; i > 0; i-- {
		if p.isRowEmpty(i - 1) {
			empty++
		} else {
			return
		}
	}
	return
}

func (p *polyominoRot) String() string {
	return p.shapeSquare.String()
}
