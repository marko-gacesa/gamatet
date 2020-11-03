// Copyright (c) 2020 by Marko Gaćeša

package op

import (
	"gamatet/game/block"
	"gamatet/game/field"
	"gamatet/game/piece"
)

func updatePieceShadow(f *field.Field, ctrl *piece.Ctrl) {
	p := ctrl.Piece

	if p == nil {
		ctrl.IsShadowShown = false
		return
	}

	ctrl.IsShadowShown = true
	ctrl.Shadow.ColL = ctrl.X + p.LeftEmptyColumns()
	ctrl.Shadow.ColR = ctrl.X + p.DimX() - p.RightEmptyColumns()
	ctrl.Shadow.Blocks = ctrl.Shadow.Blocks[:0]

	switch p.Type() {
	case piece.TypeStandard:
		blockType := ctrl.Blocks[0].Type
		if blockType == block.TypeLava || blockType == block.TypeAcid {
			setLiquidPieceShadow(f, ctrl)
			return
		}

		setSolidPieceShadow(f, ctrl)

	case piece.TypeShooter:
		x := ctrl.X
		y := ctrl.Y
		height := f.GetHeightToTopmostEmpty(x, y)
		ctrl.Shadow.Blocks = append(ctrl.Shadow.Blocks, block.XYB{
			XY:    block.XY{X: x, Y: y - height},
			Block: p.Get(0, 0),
		})
	}
}

func updateAllPiecesShadow(f *field.Field) {
	n := f.Ctrls()
	for i := 0; i < n; i++ {
		ctrl := f.Ctrl(byte(i))
		updatePieceShadow(f, ctrl)
	}
}

func setSolidPieceShadow(f *field.Field, ctrl *piece.Ctrl) {
	height := f.GetDropHeight(ctrl.Idx, !f.PieceCollision)
	if height == 0 {
		ctrl.IsShadowShown = false
		return
	}

	for i := len(ctrl.Blocks) - 1; i >= 0; i-- {
		x := ctrl.X + ctrl.Blocks[i].X
		y := ctrl.Y + ctrl.Blocks[i].Y
		ctrl.Shadow.Blocks = append(ctrl.Shadow.Blocks, block.XYB{
			XY:    block.XY{X: x, Y: y - height},
			Block: ctrl.Blocks[i].Block,
		})
	}
}

func setLiquidPieceShadow(f *field.Field, ctrl *piece.Ctrl) {
	top := make(map[int]int)
	for i := len(ctrl.Blocks) - 1; i >= 0; i-- {
		x := ctrl.X + ctrl.Blocks[i].X
		y := ctrl.Y + ctrl.Blocks[i].Y
		height := f.GetHeightToTopmostEmpty(x, y)

		h := height - top[x]
		if h > 0 && !isPieceAt(x, y-h, ctrl) {
			ctrl.Shadow.Blocks = append(ctrl.Shadow.Blocks, block.XYB{
				XY:    block.XY{X: x, Y: y - h},
				Block: ctrl.Blocks[i].Block,
			})
		}

		top[x] += 1
	}
}

func isPieceAt(x, y int, ctrl *piece.Ctrl) bool {
	for i := range ctrl.Blocks {
		px := ctrl.Blocks[i].X + ctrl.X
		py := ctrl.Blocks[i].Y + ctrl.Y
		if px == x && py == y {
			return true
		}
	}
	return false
}
