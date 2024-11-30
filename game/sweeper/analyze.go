// Copyright (c) 2020-2024 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/op"
)

type Analyzer struct {
	Field   *field.Field
	added   int
	removed int
}

func (a *Analyzer) Analyze(e event.Event) {
	switch v := e.(type) {
	case *op.FieldBlockSet:
		switch v.Op {
		case op.TypeSet:
			a.added++
		case op.TypeClear:
			a.removed++
		}
	case *op.FieldDestroyRow:
		a.removed += a.Field.GetWidth()
	case *op.FieldDestroyColumn:
		a.removed++
	}
}
