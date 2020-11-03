// Copyright (c) 2020 by Marko Gaćeša

package op

import "gamatet/game/event"

type Analyzer struct {
	HasAdded   bool
	HasRemoved bool
}

func (a *Analyzer) Analyze(e event.Event) {
	switch v := e.(type) {
	case *FieldBlockSet:
		a.HasAdded = a.HasAdded || v.Op == OpSet
		a.HasRemoved = a.HasRemoved || v.Op == OpClear
	case *FieldDestroyRow:
		a.HasRemoved = true
	case *FieldDestroyColumn:
		a.HasRemoved = true
	}
}

func (a *Analyzer) Reset() {
	a.HasAdded = false
	a.HasRemoved = false
}
