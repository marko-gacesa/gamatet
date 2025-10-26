// Copyright (c) 2020-2025 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/op"
)

type Analyzer struct {
	Field *field.Field

	blocks delta
	stats  delta

	endMode *field.Mode
}

type delta struct {
	added    int
	removed  int
	hardened int
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
	case *op.FieldDestroyColumn:
		a.blocks.removed++
	case *op.FieldStat:
		a.stats.removed += int(v.BlocksRemoved)
		a.stats.softened += int(v.BlocksSoftened)
	case *op.FieldMode:
		if v.ModeNew == field.ModeGameOver || v.ModeNew == field.ModeVictory || v.ModeNew == field.ModeDefeat {
			a.endMode = &v.ModeNew
		}
	}
}
