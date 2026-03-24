// Copyright (c) 2020-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"math/rand"
	"time"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/logic/anim"
)

const (
	AnimNo      = iota
	AnimDestroy // block destroy animation, the same as for DestroyColumn and DestroyRow is Op=clear
	AnimPop     // Pop-out for Op=clear, Pop-in for
	AnimFall    // Fall from the top (AnimParam holds height)
	AnimSpin    // Spin in place
	AnimMeld
	AnimCurtain
)

type exElem struct {
	block.XY
	elem
	next *exElem
}

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
	for range 6 {
		dx := rnd.Int()%11 - 5
		dy := rnd.Int()%11 - 5
		now := time.Now()
		f.addExBlock(x+dx, y+dy, b,
			anim.NewTransQuad(now, piece.DurationAnimBlockChange, float32(dx), float32(dy), 2*rnd.Float32()-1),
			anim.NewSpinOnce(now, piece.DurationAnimBlockChange),
			anim.NewPopOut(now, piece.DurationAnimBlockChange))
	}
}

func (f *Field) AnimQuake(intensity byte) {
	if !f.Config.Anim {
		return
	}

	f.Anim(anim.NewQuake(time.Now(), intensity))
}
