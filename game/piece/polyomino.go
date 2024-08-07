// Copyright (c) 2020 by Marko Gaćeša

package piece

import (
	"errors"
	"gamatet/game/block"
	"gamatet/game/serialize"
	"io"
	"strings"
)

type polyominoShape struct {
	dim  int      // Piece dimension
	size byte     // number of blocks in the Piece
	rots byte     // total number of rotation positions
	data bitarray // dim x dim matrix, 1 if has a block, 0 if not
}

func (p *polyominoShape) String() string {
	sb := strings.Builder{}

	dim := p.dim
	for j := 0; j < dim; j++ {
		for i := 0; i < dim; i++ {
			if p.data.get(j*dim + i) {
				sb.WriteString("[]")
			} else {
				sb.WriteString(". ")
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}

type polyomino struct {
	polyominoShape
	rot   byte        // current rotation position
	block block.Block // block material for the Piece
}

var _ Piece = (*polyomino)(nil)

func (p *polyomino) Write(w io.Writer) error {
	var buffer [4]byte
	buffer[0] = byte(p.dim)
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

func (p *polyomino) Read(r io.Reader) (err error) {
	var buffer [4]byte

	n, err := r.Read(buffer[:])
	if err != nil {
		return
	}

	if n != 4 {
		err = errors.New("failed to read polyomino")
		return
	}

	p.dim = int(buffer[0])
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

func (p *polyomino) Clone() Piece {
	clone := *p
	return &clone
}

func (p *polyomino) Equals(other Piece) bool {
	q, ok := other.(*polyomino)
	return ok && *p == *q
}

func (*polyomino) Type() Type        { return TypeStandard }
func (p *polyomino) BlockCount() int { return int(p.size) }
func (p *polyomino) DimX() int       { return p.dim }
func (p *polyomino) DimY() int       { return p.dim }

func (*polyomino) CanActivate() bool       { return false }
func (*polyomino) GetActivationCount() int { return 0 }
func (*polyomino) SetActivationCount(int)  {}

func (p *polyomino) CurrentRot() int { return int(p.rot) }
func (p *polyomino) Rots() int       { return int(p.rots) }

func (p *polyomino) RotateCW() (inverted bool) {
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
	for j := 0; j < dim2; j++ {
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

func (p *polyomino) RotateCCW() (inverted bool) {
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
	for j := 0; j < dim2; j++ {
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

func (p *polyomino) WallKick() int {
	return p.dim / 2
}

func (p *polyomino) IsEmpty(x, y int) bool {
	if x < 0 || x >= p.dim || y < 0 || y >= p.dim {
		return true
	}

	idx := y*p.dim + x
	return !p.data.get(idx)
}

func (p *polyomino) Get(x, y int) (b block.Block) {
	if p.IsEmpty(x, y) {
		return
	}
	b = p.block
	return
}

func (p *polyomino) IsRowEmpty(r int) bool {
	lim := (r + 1) * p.dim
	for idx := r * p.dim; idx < lim; idx++ {
		if p.data.get(idx) {
			return false
		}
	}
	return true
}

func (p *polyomino) IsColumnEmpty(c int) bool {
	lim := p.dim * p.dim
	for idx := c; idx < lim; idx += p.dim {
		if p.data.get(idx) {
			return false
		}
	}
	return true
}

func (p *polyomino) LeftEmptyColumns() (empty int) {
	n := p.dim
	for i := 0; i < n; i++ {
		if p.IsColumnEmpty(i) {
			empty++
		} else {
			return
		}
	}
	return
}

func (p *polyomino) RightEmptyColumns() (empty int) {
	for i := p.dim - 1; i >= 0; i-- {
		if p.IsColumnEmpty(i) {
			empty++
		} else {
			return
		}
	}
	return
}

func (p *polyomino) TopEmptyRows() (empty int) {
	n := p.dim
	for i := 0; i < n; i++ {
		if p.IsRowEmpty(i) {
			empty++
		} else {
			return
		}
	}
	return
}

func (p *polyomino) BottomEmptyRows() (empty int) {
	for i := p.dim - 1; i >= 0; i-- {
		if p.IsRowEmpty(i) {
			empty++
		} else {
			return
		}
	}
	return
}

func (p *polyomino) String() string {
	sb := strings.Builder{}

	dim := p.dim
	for j := 0; j < dim; j++ {
		for i := 0; i < dim; i++ {
			if p.data.get(j*dim + i) {
				sb.WriteString("[]")
			} else {
				sb.WriteString(". ")
			}
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
