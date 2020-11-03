// Copyright (c) 2020 by Marko Gaćeša

package block

// Type represents the type of field blocks. Different types are typically rendered differently.
type Type byte

const (
	// TypeEmpty is not actually a block but an empty space.
	TypeEmpty Type = iota

	// TypeRock is ordinary block type.
	TypeRock

	// TypeAcid is types of a block that will melt the block below it.
	TypeAcid

	// TypeLava behaves as a TypeRock, but a Piece made of it will collapse after falling.
	TypeLava

	TypeRuby

	// TypeWall is used for walls around (and inside) the play area.
	TypeWall Type = 255
)

func (t Type) IsImmovable() bool {
	return t == TypeWall
}

const (
	// HardnessMax is special, maximum, value for block hardness that can't be reduced.
	HardnessMax byte = 0xFF
)

var (
	Rock = Block{Type: TypeRock, Hardness: 0, Color: 0xA0A0A0FF}
	Hard = Block{Type: TypeRock, Hardness: 1, Color: 0x808080FF}
	Lava = Block{Type: TypeLava, Hardness: 0, Color: 0xFF8000FF}
	Acid = Block{Type: TypeAcid, Hardness: 0, Color: 0x00FF00FF}
	Wall = Block{Type: TypeWall, Hardness: HardnessMax, Color: 0x808080FF}
	Ruby = Block{Type: TypeRuby, Hardness: HardnessMax, Color: 0xFF0000FF}
)
