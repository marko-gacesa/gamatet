// Copyright (c) 2025 by Marko Gaćeša

package piece

import (
	"errors"
	"gamatet/game/block"
	"gamatet/game/serialize"
	"io"
)

type polyominoFlip struct {
	shapeRect
	block block.Block // block material for the Piece
}

var _ Piece = (*polyominoFlip)(nil)

func (p *polyominoFlip) Write(w io.Writer) error {
	var buffer [3]byte
	buffer[0] = p.width
	buffer[1] = p.height
	buffer[2] = p.size

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

func (p *polyominoFlip) Read(r io.Reader) (err error) {
	var buffer [3]byte

	n, err := r.Read(buffer[:])
	if err != nil {
		return
	}

	if n != 4 {
		err = errors.New("failed to read polyomino rectangle")
		return
	}

	p.width = buffer[0]
	p.height = buffer[1]
	p.size = buffer[2]

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

func (p *polyominoFlip) Clone() Piece {
	clone := *p
	return &clone
}

func (p *polyominoFlip) Equals(other Piece) bool {
	q, ok := other.(*polyominoFlip)
	return ok && *p == *q
}

func (*polyominoFlip) Type() Type         { return TypeStandard }
func (p *polyominoFlip) BlockCount() byte { return p.size }
func (p *polyominoFlip) DimX() byte       { return p.width }
func (p *polyominoFlip) DimY() byte       { return p.height }

func (*polyominoFlip) CanActivate() bool        { return false }
func (*polyominoFlip) GetActivationCount() byte { return 0 }
func (*polyominoFlip) SetActivationCount(byte)  {}

func (p *polyominoFlip) FlipV() {
	w := p.width
	h := p.height
	h2 := h >> 1
	for y := byte(0); y < h2; y++ {
		for x := byte(0); x < w; x++ {
			idx0 := y*w + x
			idx1 := (h-y-1)*w + x
			p.data = p.data.exchange(idx0, idx1)
		}
	}
}

func (p *polyominoFlip) FlipH() {
	w := p.width
	h := p.height
	w2 := w >> 1
	for y := byte(0); y < h; y++ {
		for x := byte(0); x < w2; x++ {
			idx0 := y*w + x
			idx1 := y*w + (w - x - 1)
			p.data = p.data.exchange(idx0, idx1)
		}
	}
}

func (p *polyominoFlip) RotateCW() (inverted bool)  { return false }
func (p *polyominoFlip) RotateCCW() (inverted bool) { return false }
func (p *polyominoFlip) WallKick() byte             { return 0 }

func (p *polyominoFlip) IsEmpty(x, y int) bool {
	w := int(p.width)
	h := int(p.height)
	if x < 0 || x >= w || y < 0 || y >= h {
		return true
	}

	idx := y*w + x
	return !p.data.get(byte(idx))
}

func (p *polyominoFlip) Get(x, y int) (b block.Block) {
	if p.IsEmpty(x, y) {
		return
	}
	b = p.block
	return
}

func (p *polyominoFlip) LeftEmptyColumns() (empty byte)  { return 0 }
func (p *polyominoFlip) RightEmptyColumns() (empty byte) { return 0 }
func (p *polyominoFlip) TopEmptyRows() (empty byte)      { return 0 }
func (p *polyominoFlip) BottomEmptyRows() (empty byte)   { return 0 }

func (p *polyominoFlip) String() string { return p.shapeRect.String() }
