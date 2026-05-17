// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import "github.com/marko-gacesa/gamatet/game/piece"

func (f *Field) Ctrls() int {
	return len(f.pieces)
}

func (f *Field) Ctrl(idx byte) *piece.Ctrl {
	return f.pieces[idx]
}

func (f *Field) CtrlLevel(idx byte) uint {
	return f.Ctrl(idx).Level
}

func (f *Field) CtrlPieceCount(idx byte) uint {
	return f.Ctrl(idx).PieceCount
}

func (f *Field) CtrlPieceOverridden(ctrlIdx byte, pieceIdx uint) bool {
	return f.Ctrl(ctrlIdx).Feed.Overridden(pieceIdx)
}

func (f *Field) CtrlStateIsTerminal(ctrlIdx byte) bool {
	return f.Ctrl(ctrlIdx).State.IsTerminal()
}
