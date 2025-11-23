// Copyright (c) 2020-2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

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

	// TypeRuby is immovable block type.
	TypeRuby

	// TypeAcid is types of a block that will melt the block below it.
	// Only as piece material. Can't appear in field.
	TypeAcid

	// TypeLava behaves as a TypeRock, but a Piece made of it will collapse after falling.
	// Only as piece material. Can't appear in field.
	TypeLava

	// TypeCurl "quantum tunnels" through blocks to fill the first hole it encounters.
	// Only as piece material. Can't appear in field.
	TypeCurl

	// TypeWave "quantum tunnels" through blocks to fill the last hole it encounters.
	// Only as piece material. Can't appear in field.
	TypeWave

	// TypeBomb immediately destroys killable blocks below it.
	// Only as piece material. Can't appear in field.
	TypeBomb

	// TypeGoal is special block type used as a goal.
	TypeGoal

	// TypeWall is used for walls around (and inside) the play area - indestructible and unmovable.
	TypeWall Type = 255
)

func (t Type) IsImmovable() bool     { return t == TypeWall || t == TypeRuby }
func (t Type) SupportsExBlock() bool { return t == TypeRock }

// Shootable returns true f the block type is directly destroyable by shooter pieces.
func (t Type) Shootable() bool { return t == TypeRock || t == TypeRuby }

func (t Type) NoSlide() bool {
	return t == TypeAcid || t == TypeLava || t == TypeCurl || t == TypeWave || t == TypeBomb
}

// Destructible returns true f the block type is directly destroyable.
func (t Type) Destructible() bool { return t == TypeRock || t == TypeRuby || t == TypeGoal }

const (
	// HardnessMax is special, maximum, value for block hardness that can't be reduced.
	HardnessMax byte = 0xFF
)

var (
	Rock = Block{Type: TypeRock, Hardness: 0, Color: 0xA0A0A0FF}
	Hard = Block{Type: TypeRock, Hardness: 1, Color: 0x808080FF}
	Iron = Block{Type: TypeIron, Hardness: HardnessMax, Color: 0xFFFFFFFF}
	Ruby = Block{Type: TypeRuby, Hardness: 0, Color: 0xA0A0A0FF}
	Acid = Block{Type: TypeAcid, Hardness: 0, Color: 0x00FF00FF}
	Lava = Block{Type: TypeLava, Hardness: 0, Color: 0xFF8000FF}
	Curl = Block{Type: TypeCurl, Hardness: 0, Color: 0xFF00C0FF}
	Wave = Block{Type: TypeWave, Hardness: 0, Color: 0x00C0FFFF}
	Bomb = Block{Type: TypeBomb, Hardness: 0, Color: 0xF0F0F0FF}
	Goal = Block{Type: TypeGoal, Hardness: 0, Color: 0xFF0000FF}
	Wall = Block{Type: TypeWall, Hardness: HardnessMax, Color: 0x606060FF}
)
