// Copyright (c) 2020, 2025 by Marko Gaćeša

package piece

import (
	"fmt"
	"gamatet/game/block"
	"gamatet/game/serialize"
	"io"
)

type Type byte

const (
	// TypeStandard have all blocks of the same type (Rock, Lava or Acid) and color.
	TypeStandard Type = iota
	// TypeShooter is 1x1 block that shoots bullets, and disappears after falling.
	TypeShooter
)

type Piece interface {
	Write(w io.Writer) error
	Read(r io.Reader) error

	Clone() Piece
	Equals(Piece) bool

	Type() Type
	BlockCount() byte
	DimX() byte
	DimY() byte

	CanActivate() bool
	GetActivationCount() byte
	SetActivationCount(count byte)

	FlipV()
	FlipH()
	RotateCW() bool
	RotateCCW() bool

	// WallKick is maximum distance the piece is allowed to move left or right if a wall prevents rotation.
	WallKick() byte

	IsEmpty(x, y int) bool
	Get(x, y int) block.Block

	LeftEmptyColumns() byte
	RightEmptyColumns() byte
	TopEmptyRows() byte
	BottomEmptyRows() byte

	fmt.Stringer
}

func GetBlocks(p Piece, blocks []block.XYB) []block.XYB {
	dimX := int(p.DimX())
	dimY := int(p.DimY())
	for j := 0; j < dimY; j++ {
		for i := 0; i < dimX; i++ {
			pBlock := p.Get(i, j)
			if pBlock.Type == block.TypeEmpty {
				continue
			}

			blocks = append(blocks, block.XYB{
				XY:    block.XY{X: i, Y: -j},
				Block: pBlock,
			})
		}
	}

	return blocks
}

func Write(w io.Writer, p Piece) (err error) {
	switch v := p.(type) {
	case *polyominoRot:
		err = serialize.Write8(w, 'P')
		if err != nil {
			return
		}
		err = v.Write(w)
	case *shooter:
		err = serialize.Write8(w, 'S')
		if err != nil {
			return
		}
		err = v.Write(w)
	default:
		err = fmt.Errorf("unsupported piece type %T", p)
	}

	return
}

func Read(r io.Reader) (p Piece, err error) {
	code, err := serialize.Read8(r)
	if err != nil {
		return
	}

	switch code {
	case 'P':
		p = &polyominoRot{}
		err = p.Read(r)
	case 'S':
		p = &shooter{}
		err = p.Read(r)
	default:
		err = fmt.Errorf("unrecognized piece code: %d", code)
	}

	return
}
