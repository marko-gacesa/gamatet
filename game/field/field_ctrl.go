// Copyright (c) 2020 by Marko Gaćeša

package field

import "gamatet/game/piece"

type PiecePlace struct {
	FieldIdx byte
	CtrlIdx  byte
}

func (f *Field) CtrlInfoPositions() []piece.DisplayPosition {
	var ctrlInfos []piece.DisplayPosition
	for i := 0; i < f.Ctrls(); i++ {
		info := f.Ctrl(byte(i)).InfoPosition
		if info == piece.DisplayPositionOff {
			continue
		}
		ctrlInfos = append(ctrlInfos, info)
	}
	return ctrlInfos
}
