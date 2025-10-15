// Copyright (c) 2020, 2025 by Marko Gaćeša

package field

import (
	"gamatet/game/block"
	"gamatet/game/piece"
	"gamatet/logic/anim"
)

const MaxPieces = 4

type elem struct {
	block.Block
	anim.List
}

type exElem struct {
	block.XY
	elem
	next *exElem
}

type Field struct {
	Idx      int
	w        int
	h        int
	blocks   []elem
	pieces   []*piece.Ctrl
	firstEx  *exElem
	animList anim.List
	paused   bool
	doneCh   chan struct{}
	Config
}

type Config struct {
	PieceCollision bool
	Anim           bool
}

const (
	MinWidth  = 4
	MaxWidth  = 40
	MinHeight = 4
	MaxHeight = 40
)

func Make(dimW, dimH, pieceCount int) (f *Field) {
	if dimW < MinWidth {
		panic("too narrow")
	} else if dimW > MaxWidth {
		panic("too wide")
	}

	if dimH < MinHeight {
		panic("too low")
	} else if dimH > MaxHeight {
		panic("too high")
	}

	if pieceCount < 0 || pieceCount > MaxPieces {
		panic("invalid piece count")
	}

	f = &Field{
		w:      dimW,
		h:      dimH,
		blocks: make([]elem, dimW*dimH),
		pieces: make([]*piece.Ctrl, pieceCount),
		doneCh: make(chan struct{}),
	}

	for i := 0; i < pieceCount; i++ {
		f.pieces[i] = piece.NewCtrl(i)
	}

	return
}

func (f *Field) Ctrls() int {
	return len(f.pieces)
}

func (f *Field) Ctrl(idx byte) *piece.Ctrl {
	return f.pieces[idx]
}

func (f *Field) GetDone() <-chan struct{} {
	return f.doneCh
}

func (f *Field) CloseDone() {
	close(f.doneCh)
}

func (f *Field) StartTimers() {
	for _, ctrl := range f.pieces {
		ctrl.RestartTimer(0)
	}
}

func (f *Field) StopTimers() {
	for _, ctrl := range f.pieces {
		ctrl.StopTimer()
	}
}

func (f *Field) Pause() {
	for _, ctrl := range f.pieces {
		ctrl.StopTimer()
		ctrl.PausedState = ctrl.State
		ctrl.State = piece.StatePause
	}
	f.paused = true
}

func (f *Field) Unpause() {
	for _, ctrl := range f.pieces {
		ctrl.State = ctrl.PausedState
		ctrl.PausedState = piece.StatePause
		ctrl.RestartTimer(0)
	}
	f.paused = false
}

func (f *Field) Anim(a anim.Anim) {
	f.animList.Add(a)
}

func (f *Field) setXY(x, y int, b block.Block) *anim.List {
	idx := y*f.w + x
	f.blocks[idx].Block = b
	f.blocks[idx].List.Clear()
	return &f.blocks[idx].List
}

func (f *Field) getXY(x, y int) (block.Block, *anim.List) {
	idx := y*f.w + x
	return f.blocks[idx].Block, &f.blocks[idx].List
}

func (f *Field) fill(b block.Block) {
	for i := 0; i < len(f.blocks); i++ {
		f.blocks[i].Block = b
		f.blocks[i].Clear()
	}
}

func (f *Field) clear() {
	for i := 0; i < len(f.blocks); i++ {
		f.blocks[i] = elem{}
	}
}

// _getXYPieceIdx returns piece index if on field coordinates there is a piece or -1 if there is none
func (f *Field) _getXYPieceIdx(x, y int) int {
	for i, ctrl := range f.pieces {
		if ctrl.Piece == nil {
			continue
		}

		if !ctrl.Piece.IsEmpty(x-ctrl.X, ctrl.Y-y) {
			return i
		}
	}

	return -1
}

func (f *Field) _isXYEmpty(x, y, colMin, colMax int, liftAll bool, liftPiece int) bool {
	if x < colMin || x > colMax || y < 0 || y >= f.h {
		return false
	}

	if f.blocks[y*f.w+x].Type != block.TypeEmpty {
		return false
	}

	if liftAll {
		return true
	}

	p := f._getXYPieceIdx(x, y)
	return p == -1 || p == liftPiece
}

func (f *Field) _canPlacePiece(px, py, colMin, colMax int, p piece.Piece, liftAll bool, liftPiece int) bool {
	dimX := int(p.DimX())
	dimY := int(p.DimY())
	for j := 0; j < dimY; j++ {
		for i := 0; i < dimX; i++ {
			if p.IsEmpty(i, j) {
				continue
			}

			if !f._isXYEmpty(px+i, py-j, colMin, colMax, liftAll, liftPiece) {
				return false
			}
		}
	}

	return true
}
