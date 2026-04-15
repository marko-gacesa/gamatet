// Copyright (c) 2020-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

type Analyzer struct {
	Field *field.Field

	blocks delta
	stats  deltaStats
	shots  []block.Type

	endMode *field.Mode
}

type delta struct {
	added     int
	removed   int
	hardened  int
	softened  int
	goal      byte
	gnawKills []block.XYB
}

type deltaStats struct {
	removed  int
	softened int
}

func (a *Analyzer) Analyze(e event.Event) {
	switch v := e.(type) {
	case *op.FieldBlockSet:
		switch v.Op {
		case op.TypeSet:
			a.blocks.added++
		case op.TypeClear:
			a.blocks.removed++
			switch v.Block.Type {
			case block.TypeGoal:
				a.blocks.goal++
			case block.TypeGnaw:
				a.blocks.gnawKills = append(a.blocks.gnawKills, block.XYB{
					XY:    block.XY{X: int(v.Col), Y: int(v.Row)},
					Block: v.Block,
				})
			}
		}
	case *op.FieldBlockHardness:
		switch {
		case v.Hardness > 1:
			a.blocks.hardened += int(v.Hardness)
		case v.Hardness < 1:
			a.blocks.softened -= int(v.Hardness)
		}
	case *op.FieldDestroyRow:
		a.blocks.removed += a.Field.GetWidth()
		for col, b := range v.Blocks {
			switch b.Type {
			case block.TypeGoal:
				a.blocks.goal++
			case block.TypeGnaw:
				a.blocks.gnawKills = append(a.blocks.gnawKills, block.XYB{
					XY:    block.XY{X: col, Y: int(v.Row)},
					Block: b,
				})
			}
		}
	case *op.FieldDestroyColumn:
		a.blocks.removed++
		switch v.Block.Type {
		case block.TypeGoal:
			a.blocks.goal++
		case block.TypeGnaw:
			a.blocks.gnawKills = append(a.blocks.gnawKills, block.XYB{
				XY:    block.XY{X: int(v.Col), Y: int(v.Row)},
				Block: v.Block,
			})
		}
	case *op.FieldStat:
		a.stats.removed += int(v.BlocksRemoved)
		a.stats.softened += int(v.BlocksSoftened)
	case *op.FieldMode:
		if v.ModeNew == field.ModeGameOver || v.ModeNew == field.ModeVictory || v.ModeNew == field.ModeDefeat {
			a.endMode = &v.ModeNew
		}
	case *op.PieceShoot:
		if v.Hit {
			a.shots = append(a.shots, v.BlockType)
		}
	}
}
