// Copyright (c) 2020-2024 by Marko Gaćeša

package block

// Type represents the type of field blocks. Different types are typically rendered differently.
type Type byte

const (
	// TypeEmpty is not actually a block but an empty space.
	TypeEmpty Type = iota

	// TypeRock is ordinary block type.
	TypeRock

	// TypeIron is indestructible block type.
	TypeIron

	// TypeRuby is special block type used as a goal.
	TypeRuby

	// TypeAcid is types of a block that will melt the block below it. Used only as piece material.
	TypeAcid

	// TypeLava behaves as a TypeRock, but a Piece made of it will collapse after falling. Used only as piece material.
	TypeLava

	// TypeWave "quantum tunnels" through blocks to fill the first hole it encounters. Used only as piece material.
	TypeWave

	// TypeWall is used for walls around (and inside) the play area - indestructible and unmovable.
	TypeWall Type = 255
)

func (t Type) IsImmovable() bool { return t == TypeWall }
func (t Type) SupportsExBlock() bool {
	return t == TypeRock || t == TypeAcid || t == TypeLava || t == TypeWave
}
func (t Type) Meltable() bool { return t == TypeRock }

const (
	// HardnessMax is special, maximum, value for block hardness that can't be reduced.
	HardnessMax byte = 0xFF
)

var (
	Rock = Block{Type: TypeRock, Hardness: 0, Color: 0xA0A0A0FF}
	Hard = Block{Type: TypeRock, Hardness: 1, Color: 0x808080FF}
	Lava = Block{Type: TypeLava, Hardness: 0, Color: 0xFF8000FF}
	Acid = Block{Type: TypeAcid, Hardness: 0, Color: 0x00FF00FF}
	Wave = Block{Type: TypeWave, Hardness: 0, Color: 0xFF00C0FF}
	Wall = Block{Type: TypeWall, Hardness: HardnessMax, Color: 0x808080FF}
	Iron = Block{Type: TypeIron, Hardness: HardnessMax, Color: 0xFFFFFFFF}
	Ruby = Block{Type: TypeRuby, Hardness: 0, Color: 0xFF0000FF}
)
