// Copyright (c) 2020-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
)

func conjureBlock(p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimPop, 0, b))
}

func putBlock(p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimNo, 0, b))
}

func Init(q setup.FieldInit) func(f field.Reader, p event.Pusher) {
	switch q {
	default:
		fallthrough
	case setup.FieldInitEmpty:
		return nil
	case setup.FieldInitLowSparseBlocks:
		return InitLowSparseBlocks
	case setup.FieldInitLowDenseBlocks:
		return InitLowDenseBlocks
	case setup.FieldInitHighSparseBlocks:
		return InitHighSparseBlocks
	case setup.FieldInitHighDenseBlocks:
		return InitHighDenseBlocks
	case setup.FieldInitFunnel:
		return InitFunnel
	case setup.FieldInitTriangle:
		return InitTriangle
	}
}

func InitLowSparseBlocks(f field.Reader, p event.Pusher) {
	initRandomBlocks(0.3, 0.2, f, p)
}
func InitLowDenseBlocks(f field.Reader, p event.Pusher) {
	initRandomBlocks(0.7, 0.2, f, p)
}
func InitHighSparseBlocks(f field.Reader, p event.Pusher) {
	initRandomBlocks(0.3, 0.45, f, p)
}
func InitHighDenseBlocks(f field.Reader, p event.Pusher) {
	initRandomBlocks(0.7, 0.45, f, p)
}

func initRandomBlocks(fillPercent, heightPercent float32, f field.Reader, p event.Pusher) {
	fullWidth := f.GetWidth()
	playerWidth := f.CtrlWidth()
	playerCount := f.Ctrls()

	height := int(float32(f.GetHeight()) * heightPercent)
	blockPerLine := int(float32(playerWidth) * fillPercent)

	c := piece.NewRandomColor(setup.ColorRGB[:], 0)
	buffer := make([]uint, playerWidth)
	rnd := f.Random(0)

	for y := range height {
		for i := range buffer {
			buffer[i] = uint(i)
		}
		rnd.Perm(buffer)

		line := buffer[:blockPerLine]

		for i := range line {
			for k := range playerCount {
				x := k*playerWidth + int(line[i])
				if x < fullWidth {
					playerIdx := f.CtrlPlayerIndex(byte(k))
					b := block.Block{Type: block.TypeRock, Color: c.Color(uint(y*playerWidth+x), playerIdx)}
					conjureBlock(p, x, y, b)
				}
			}
		}
	}
}

func InitTriangle(f field.Reader, p event.Pusher) {
	playerWidth := f.CtrlWidth()
	fullWidth := f.GetWidth()
	height := f.GetHeight()

	c := piece.NewRandomColor(setup.ColorRGB[:], 0)

	for x := range fullWidth {
		var m int
		if playerWidth > height {
			m = x%playerWidth - playerWidth + height
		} else {
			m = x % playerWidth
		}
		for y := 0; y < m; y++ {
			idx := y*fullWidth + x
			putBlock(p, x, y, block.Block{
				Type:  block.TypeRock,
				Color: c.Color(uint(idx), byte(x/playerWidth)),
			})
		}
	}
}

func InitFunnel(f field.Reader, p event.Pusher) {
	playerWidth := f.CtrlWidth()
	playerCount := f.Ctrls()

	for k := range playerCount {
		for x := range playerWidth / 2 {
			for y := 0; y < playerWidth/2-x-1; y++ {
				putBlock(p, k*playerWidth+x, y, block.Wall)
				putBlock(p, (k+1)*playerWidth-x-1, y, block.Wall)
			}
		}
	}
}
