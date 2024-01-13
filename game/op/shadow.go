// Copyright (c) 2020-2024 by Marko Gaćeša

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

	ctrl.Shadow.ColL = ctrl.X + p.LeftEmptyColumns()
	ctrl.Shadow.ColR = ctrl.X + p.DimX() - p.RightEmptyColumns()
	ctrl.Shadow.Blocks = ctrl.Shadow.Blocks[:0]

	switch p.Type() {
	case piece.TypeStandard:
		switch ctrl.Blocks[0].Type {
		case block.TypeLava, block.TypeAcid:
			setLiquidPieceShadow(f, ctrl)
		case block.TypeWave:
			setWavePieceShadow(f, ctrl)
		default:
			setSolidPieceShadow(f, ctrl)
		}
	case piece.TypeShooter:
		setShooterPieceShadow(f, ctrl)
	}

	ctrl.IsShadowShown = len(ctrl.Blocks) > 0
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
		return
	}

	for i := len(ctrl.Blocks) - 1; i >= 0; i-- {
		x := ctrl.X + ctrl.Blocks[i].X
		y := ctrl.Y + ctrl.Blocks[i].Y

		yh := y - height

		if isPieceAt(x, yh, ctrl) {
			continue
		}

		ctrl.Shadow.Blocks = append(ctrl.Shadow.Blocks, block.XYB{
			XY:    block.XY{X: x, Y: yh},
			Block: ctrl.Blocks[i].Block,
		})
	}
}

func setLiquidPieceShadow(f *field.Field, ctrl *piece.Ctrl) {
	t := field.NewTemp(f)
	defer t.Revert()

	for i := len(ctrl.Blocks) - 1; i >= 0; i-- {
		x := ctrl.X + ctrl.Blocks[i].X
		y := ctrl.Y + ctrl.Blocks[i].Y
		height := f.GetHeightToTopmostEmpty(x, y)
		if height == 0 {
			continue
		}

		yh := y - height

		if isPieceAt(x, yh, ctrl) {
			continue
		}

		ctrl.Shadow.Blocks = append(ctrl.Shadow.Blocks, block.XYB{
			XY:    block.XY{X: x, Y: yh},
			Block: ctrl.Blocks[i].Block,
		})

		t.Set(x, yh, block.Rock)
	}
}

func setWavePieceShadow(f *field.Field, ctrl *piece.Ctrl) {
	t := field.NewTemp(f)
	defer t.Revert()

	for i := len(ctrl.Blocks) - 1; i >= 0; i-- {
		x := ctrl.X + ctrl.Blocks[i].X
		y := ctrl.Y + ctrl.Blocks[i].Y
		height := f.GetHeightToTopmostHole(x, y)
		if height == 0 {
			continue
		}

		yh := y - height

		ctrl.Shadow.Blocks = append(ctrl.Shadow.Blocks, block.XYB{
			XY:    block.XY{X: x, Y: yh},
			Block: ctrl.Blocks[i].Block,
		})

		t.Set(x, yh, block.Rock)
	}
}

func setShooterPieceShadow(f *field.Field, ctrl *piece.Ctrl) {
	x := ctrl.X
	y := ctrl.Y
	var height int

	switch ctrl.Blocks[0].Type {
	case block.TypeWave:
		height = f.GetHeightToTopmostHole(x, y)
	default:
		height = f.GetHeightToTopmostEmpty(x, y)
	}

	if height == 0 {
		return
	}

	ctrl.Shadow.Blocks = append(ctrl.Shadow.Blocks, block.XYB{
		XY:    block.XY{X: x, Y: y - height},
		Block: ctrl.Piece.Get(0, 0),
	})
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
