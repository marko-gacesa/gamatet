// Copyright (c) 2020-2024 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
	"io"
)

type shooter struct {
	bulletType block.Type
	ammo       int
}

var _ Piece = (*shooter)(nil)

func Shooter(ammo int, bulletType block.Type) Piece {
	return &shooter{
		bulletType: bulletType,
		ammo:       ammo,
	}
}

func (p *shooter) Write(w io.Writer) error {
	_, err := w.Write([]byte{byte(p.bulletType), byte(p.ammo)})
	return err
}

func (p *shooter) Read(r io.Reader) error {
	var buffer [2]byte
	_, err := io.ReadFull(r, buffer[:])
	if err != nil {
		return err
	}

	p.bulletType = block.Type(buffer[0])
	p.ammo = int(buffer[1])

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

func (*shooter) Type() Type      { return TypeShooter }
func (*shooter) BlockCount() int { return 1 }
func (*shooter) DimX() int       { return 1 }
func (*shooter) DimY() int       { return 1 }

func (p *shooter) CanActivate() bool        { return p.ammo > 0 }
func (p *shooter) GetActivationCount() int  { return p.ammo }
func (p *shooter) SetActivationCount(n int) { p.ammo = n }

func (*shooter) CurrentRot() int { return 0 }
func (*shooter) Rots() int       { return 0 }
func (*shooter) RotateCW() bool  { return false }
func (*shooter) RotateCCW() bool { return false }
func (*shooter) WallKick() int   { return 0 }

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

func (p *shooter) IsRowEmpty(r int) bool    { return r != 0 }
func (p *shooter) IsColumnEmpty(c int) bool { return c != 0 }
func (p *shooter) LeftEmptyColumns() int    { return 0 }
func (p *shooter) RightEmptyColumns() int   { return 0 }
func (p *shooter) TopEmptyRows() int        { return 0 }
func (p *shooter) BottomEmptyRows() int     { return 0 }

func (p *shooter) String() string {
	switch p.bulletType {
	default:
		fallthrough
	case block.TypeLava:
		return "LL"
	case block.TypeAcid:
		return "AA"
	case block.TypeWave:
		return "CC"
	}
}
