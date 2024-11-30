// Copyright (c) 2020 by Marko Gaćeša

package core

import (
	"gamatet/game/block"
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/op"
)

func conjureBlock(p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimPop, 0, b))
}

func putBlock(p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimNo, 0, b))
}
