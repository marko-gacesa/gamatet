// Copyright (c) 2025 by Marko Gaćeša

package piece

import (
	"math/rand/v2"
)

type Color interface {
	Color(idx, playerIdx int) uint32
}

type DefaultColor struct{}

func (c DefaultColor) Color(idx, playerIdx int) uint32 {
	return _colors[idx%len(_colors)]
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
	rand       rand.Rand
}

func NewRandomColor(colorTable [][3]float32, seed int) *RandomColor {
	c := &RandomColor{
		seed:       seed,
		colorTable: colorTable,
		pcg:        *rand.NewPCG(0, 0),
	}
	c.rand = *rand.New(&c.pcg)
	return c
}

func (c *RandomColor) Color(idx, playerIdx int) uint32 {
	rgb := c.colorTable[playerIdx%len(c.colorTable)]
	r := uint16(256 * rgb[0])
	g := uint16(256 * rgb[1])
	b := uint16(256 * rgb[2])
	c.pcg.Seed(uint64(c.seed), uint64(idx))
	v := uint16(c.rand.Uint64())&0xFF | 0b11000000
	return uint32(r*v)&0xFF00<<16 | uint32(g*v)&0xFF00<<8 | uint32(b*v)&0xFF00
}
