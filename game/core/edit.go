// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

func conjureBlock(p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimPop, 0, b))
}

func putBlock(p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimNo, 0, b))
}
