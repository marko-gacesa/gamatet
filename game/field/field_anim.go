// Copyright (c) 2020-2024 by Marko Gaćeša

package field

import (
	"gamatet/game/block"
	"gamatet/game/piece"
	"gamatet/logic/anim"
	"math/rand"
	"time"
)

const (
	AnimNo      = iota
	AnimDestroy // block destroy animation, the same as for DestroyColumn and DestroyRow is Op=clear
	AnimShot    // Shooter block animation (AnimParam holds height)
	AnimPop     // Pop-out for Op=clear, Pop-in for
	AnimFall    // Fall from the top
	AnimSpin    // Spin in place
	AnimRotateZ
	AnimMeld
)

func (f *Field) addExBlock(x, y int, b block.Block, anims ...anim.Anim) {
	list := anim.List{}
	list.AddAll(anims...)

	f.firstEx = &exElem{
		XY: block.XY{X: x, Y: y},
		elem: elem{
			Block: b,
			List:  list,
		},
		next: f.firstEx,
	}
}

func (f *Field) animBlockDestroy(x, y int, b block.Block) {
	if !f.Config.Anim || !b.Type.SupportsExBlock() {
		return
	}

	var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		dx := rnd.Int()%11 - 5
		dy := rnd.Int()%11 - 5
		now := time.Now()
		f.addExBlock(x+dx, y+dy, b,
			anim.NewTransQuad(now, piece.DurationAnimBlockChange, float32(dx), float32(dy), 2*rnd.Float32()-1),
			anim.NewSpin(now, piece.DurationAnimBlockChange),
			anim.NewPopOut(now, piece.DurationAnimBlockChange))
	}
}

func (f *Field) animBullet(x, y, height int, b block.Block) {
	if !f.Config.Anim || height == 0 {
		return
	}

	now := time.Now()
	duration := piece.GetFallDuration(height)

	f.addExBlock(x, y, b, anim.NewFall(now, duration, float32(height)))
}
