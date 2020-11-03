// Copyright (c) 2020 by Marko Gaćeša

package field

import (
	"gamatet/game/block"
	"gamatet/game/piece"
	"gamatet/logic/anim"
	"time"
)

func (f *Field) clearBlock(x, y int) (b block.Block) {
	idx := y*f.w + x
	b = f.blocks[idx].Block
	f.blocks[idx] = elem{}
	return
}

func (f *Field) ShiftRowsDown(y int) {
	// destroy blocks in the row=y
	for x := 0; x < f.w; x++ {
		b := f.clearBlock(x, y)
		f.animBlockDestroy(x, y, b)
	}

	// add fall animations for all the blocks
	if f.Config.Anim {
		now := time.Now().UTC()
		for i := (y + 1) * f.w; i < len(f.blocks); i++ {
			if f.blocks[i].Block.Type != block.TypeEmpty {
				f.blocks[i].List.Add(anim.NewFall(now, piece.GetFallDuration(1), 1))
			}
		}
	}

	// copy all the blocks above the row=y to the new location
	copy(f.blocks[y*f.w:], f.blocks[(y+1)*f.w:])

	// delete the top row
	idx := (f.h - 1) * f.w
	lim := idx + f.w
	for ; idx < lim; idx++ {
		f.blocks[idx] = elem{}
	}
}

func (f *Field) UndoShiftRowsDown(y int, blocks []block.Block) {
	// copy all the blocks above the row=y to the new location
	copy(f.blocks[(y+1)*f.w:], f.blocks[y*f.w:])

	// restore the row
	idx := y * f.w
	for i := 0; i < f.w; i++ {
		f.blocks[idx+i] = elem{Block: blocks[i]}
	}
}

func (f *Field) ShiftColumnDownByN(x, y, n, height int) {
	f.animBlockDestroy(x, y, f.clearBlock(x, y))

	if n == 0 {
		return
	}

	var duration time.Duration
	if f.Config.Anim {
		duration = piece.GetFallDuration(height)
	}

	// move n blocks above it by height rows down
	d := f.w * height
	idx := (y-height+1)*f.w + x
	for i := 0; i < n; i++ {
		idxSrc := idx + d
		f.blocks[idx] = f.blocks[idxSrc]
		f.blocks[idxSrc] = elem{}

		if f.Config.Anim && f.blocks[idx].Block.Type != block.TypeEmpty {
			f.blocks[idx].List.Add(anim.NewFall(time.Now(), duration, float32(height)))
		}

		idx += f.w
	}
}

func (f *Field) UndoShiftColumnByN(x, y, n, height int, b block.Block) {
	// move n blocks above it by height rows up
	d := f.w * height
	idx := (y+n)*f.w + x
	for i := 0; i < n; i++ {
		idxSrc := idx - d
		f.blocks[idx] = f.blocks[idxSrc]
		f.blocks[idxSrc] = elem{}
		idx -= f.w
	}

	f.setXY(x, y, b)
}

func (f *Field) SetXY(x, y, animType, animParam int, b block.Block) {
	if old, _ := f.getXY(x, y); b.Type == block.TypeEmpty || old.Type != block.TypeEmpty {
		panic("empty block in f.SetXY") // should not happen
	}

	animList := f.setXY(x, y, b)

	if f.Config.Anim {
		switch animType {
		case AnimShot, AnimFall:
			t := time.Now()
			var delay time.Duration

			if height := animParam; height > 0 {
				rowFull := true
				idx := y * f.w
				lim := idx + f.w
				for ; idx < lim; idx++ {
					if f.blocks[idx].Type == block.TypeEmpty {
						rowFull = false
						break
					}
				}

				delay = piece.GetFallDuration(height)

				if rowFull {
					// if this block completes a line, add external bullet animation because the block will be destroyed
					f.animBullet(x, y, height, b)
				} else {
					// animate the newly created piece - falling
					animList.Add(anim.NewFall(t, delay, float32(height)))
				}
			}

			// animate the newly created piece - color transition
			animList.Add(anim.NewColorTrans(t.Add(delay), piece.DurationAnimBlockChange, 0xFF7F00FF, 0xFFFFFFFF))
		case AnimPop:
			animList.Add(anim.NewPopIn(time.Now(), piece.DurationAnimBlockChange))
		case AnimMeld:
			animList.Add(anim.NewMeld(time.Now(), piece.DurationAnimBlockChange))
		}

		//if b.Type == block.TypeRuby {
		//	animList.Add(anim.NewRotateZ(time.Now(), 4*time.Second))
		//}
	}
}

func (f *Field) ClearXY(x, y, animType, animParam int) (b block.Block) {
	b = f.clearBlock(x, y)

	if b.Type == block.TypeEmpty {
		panic("empty block in f.ClearXY") // should not happen
	}

	if f.Config.Anim {
		switch animType {
		case AnimDestroy:
			f.animBlockDestroy(x, y, b)

		/*
			case AnimShot:
				if height := animParam; height > 0 {
					t := time.Now().Add(piece.GetFallDuration(height))
					f.addExBlock(x, y+height, b,
						anim.NewFall(t, piece.GetFallDuration(height), float32(-height)),
						anim.NewPopOut(t, piece.GetFallDuration(height)))
				}

		*/

		case AnimPop:
			f.addExBlock(x, y, b, anim.NewPopOut(time.Now(), piece.DurationAnimBlockChange))

		case AnimFall, AnimShot:
			t := time.Now()

			if height := animParam; height > 0 {
				t.Add(piece.GetFallDuration(height))
				f.animBullet(x, y, height, block.Acid)
			}

			f.addExBlock(x, y, b,
				anim.NewSpin(t, piece.DurationAnimBlockChange),
				anim.NewPopOut(t, piece.DurationAnimBlockChange))
		}
	}

	return
}

func (f *Field) HardnessXY(x, y, delta, animType, animParam int) (blockOld, blockNew block.Block) {
	idx := y*f.w + x

	blockOld = f.blocks[idx].Block

	if -delta > int(blockOld.Hardness) || blockOld.Hardness == block.HardnessMax {
		panic("unsupported action in f.HardnessXY") // should not happen
	}

	f.blocks[idx].Hardness += byte(delta)

	blockOld = f.blocks[idx].Block

	if f.Config.Anim {
		t := time.Now()

		switch animType {
		case AnimShot:
			if height := animParam; height > 0 && delta < 0 {
				t.Add(piece.GetFallDuration(height))
				f.animBullet(x, y, height, block.Acid)
			}
			fallthrough

		case AnimSpin, AnimFall:
			f.blocks[idx].List.Add(anim.NewSpin(t, piece.DurationAnimBlockChange))
		}
	}

	return
}
