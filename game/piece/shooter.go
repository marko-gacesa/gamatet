// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"io"

	"github.com/marko-gacesa/gamatet/game/block"
)

type shooter struct {
	bulletType block.Type
	ammo       byte
}

var _ Piece = (*shooter)(nil)

func Shooter(ammo byte, bulletType block.Type) Piece {
	return &shooter{
		bulletType: bulletType,
		ammo:       ammo,
	}
}

func (p *shooter) Write(w io.Writer) error {
	_, err := w.Write([]byte{byte(p.bulletType), p.ammo})
	return err
}

func (p *shooter) Read(r io.Reader) error {
	var buffer [2]byte
	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return err
	}

	p.bulletType = block.Type(buffer[0])
	p.ammo = buffer[1]

	return nil
}

func (p *shooter) Clone() Piece {
	clone := *p
	return &clone
}

func (p *shooter) Equals(other Piece) bool {
	q, ok := other.(*shooter)
	return ok && p.bulletType == q.bulletType && p.ammo == q.ammo
}

func (*shooter) Type() Type       { return TypeShooter }
func (*shooter) BlockCount() byte { return 1 }
func (*shooter) DimX() byte       { return 1 }
func (*shooter) DimY() byte       { return 1 }

func (p *shooter) CanActivate() bool     { return p.ammo > 0 }
func (p *shooter) ActivationCount() byte { return p.ammo }
func (p *shooter) Activate() bool {
	if p.ammo == 0 {
		return false
	}

	p.ammo--
	return true
}

func (p *shooter) UndoActivate() bool {
	p.ammo++
	return true
}

func (*shooter) WallKick() byte { return 0 }

func (p *shooter) IsEmpty(x, y int) bool {
	return x != 0 || y != 0
}

func (p *shooter) Get(x, y int) (b block.Block) {
	if p.IsEmpty(x, y) {
		return
	}

	switch p.bulletType {
	default:
		fallthrough
	case block.TypeRock:
		return block.Rock
	case block.TypeLava:
		return block.Lava
	case block.TypeAcid:
		return block.Acid
	case block.TypeCurl:
		return block.Curl
	case block.TypeWave:
		return block.Wave
	case block.TypeBomb:
		return block.Bomb
	}
}

func (p *shooter) LeftEmptyColumns() byte  { return 0 }
func (p *shooter) RightEmptyColumns() byte { return 0 }
func (p *shooter) TopEmptyRows() byte      { return 0 }
func (p *shooter) BottomEmptyRows() byte   { return 0 }

func (p *shooter) String() string {
	switch p.bulletType {
	default:
		fallthrough
	case block.TypeRock:
		return "RR"
	case block.TypeLava:
		return "LL"
	case block.TypeAcid:
		return "AA"
	case block.TypeCurl:
		return "CC"
	case block.TypeWave:
		return "WW"
	case block.TypeBomb:
		return "BB"
	}
}
