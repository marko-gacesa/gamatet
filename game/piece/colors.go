// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"math/rand/v2"

	"github.com/marko-gacesa/gamatet/game/block"
)

type Color interface {
	Color(idx uint, playerIdx byte) uint32
}

type RockColor struct{}

func (c RockColor) Color(uint, byte) uint32 {
	return block.Rock.Color
}

type DefaultColor struct{}

func (c DefaultColor) Color(idx uint, playerIdx byte) uint32 {
	return _colors[idx%uint(len(_colors))]
}

var _colors = []uint32{
	0xFFFF0000, // yellow
	0x00FFFF00, // cyan
	0xFF00FF00, // magenta
	0x0000FF00, // blue
	0xFF7F0000, // orange
	0x00FF0000, // green
	0xFF000000, // red
	0xFFD70000, // gold
	0x7F00FF00, // violet
	0x7FFF0000, // lime
	0x007FFF00, // azure
	0x3F00FF00, // indigo
	0xF5F5DC00, // beige
	0x964B0000, // brown
	0xFFC0CB00, // pink
	0xFF007F00, // pink-red
	0x00FF7F00, // spring-green
	0xFF007F00, // electric-magenta
	0xC3B09100, // khaki
	0xFFE5B400, // peach
	0x80008000, // purple
	0x00808000, // teal
	0x80800000, // olive
	0x40E0D000, // turquoise
}

type RandomColor struct {
	seed       int
	colorTable [][3]float32
	pcg        rand.PCG
}

func NewRandomColor(colorTable [][3]float32, seed int) *RandomColor {
	c := &RandomColor{
		seed:       seed,
		colorTable: colorTable,
		pcg:        *rand.NewPCG(0, 0),
	}
	return c
}

func (c *RandomColor) Color(idx uint, playerIdx byte) uint32 {
	random := rand.New(&c.pcg)
	rgb := c.colorTable[playerIdx%byte(len(c.colorTable))]
	r := clamp(rgb[0] + 0.4*(rand.Float32()-0.5))
	g := clamp(rgb[1] + 0.4*(rand.Float32()-0.5))
	b := clamp(rgb[2] + 0.4*(rand.Float32()-0.5))
	r16 := uint16(256 * r)
	g16 := uint16(256 * g)
	b16 := uint16(256 * b)
	c.pcg.Seed(uint64(c.seed), uint64(idx))
	v := uint16(random.Uint64())&0xFF | 0b11000000
	return uint32(r16*v)&0xFF00<<16 | uint32(g16*v)&0xFF00<<8 | uint32(b16*v)&0xFF00
}

func clamp(f float32) float32 {
	return max(min(1, f), 0)
}
