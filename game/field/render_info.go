// Copyright (c) 2020-2025 by Marko Gaćeša

package field

import (
	"gamatet/game/block"
	"gamatet/game/piece"
	"gamatet/logic/anim"
	"sync"
	"time"
)

type RenderRequest struct {
	FieldIdx   int
	Time       time.Time
	RenderInfo chan<- *RenderInfo
}

type BlockRenderInfo struct {
	block.XYB
	anim.Result
}

type NextPieceRenderInfo struct {
	Type   piece.Type
	Blocks []block.XYB
}

type PieceRenderInfo struct {
	Position piece.DisplayPosition
	State    piece.State

	PieceTextData
	IsLimited bool
	Limits    piece.ColumnLimit

	DirectionCW bool
	NextPieces  [piece.NextBlockCount]NextPieceRenderInfo

	PieceEmpty bool

	Blocks     []block.XYB
	X, Y       int
	DimX, DimY int
	Type       piece.Type
	ActCount   int
	Result     anim.Result
	DrawShadow bool
	Shadow     piece.Shadow
}

type PieceTextData struct {
	Name     string
	Score    string
	PieceNum string
	Level    string
}

type TextData struct {
	BlocksRemoved string
	Latencies     string
}

type RenderInfo struct {
	W, H   int
	Mode   Mode
	Blocks []BlockRenderInfo
	Pieces [MaxPieces]PieceRenderInfo
	Result anim.Result
	TextData
}

var syncPoolRenderInfo = &sync.Pool{
	New: func() any {
		info := &RenderInfo{}
		info.Blocks = make([]BlockRenderInfo, 0, 256)
		for i := range len(info.Pieces) {
			info.Pieces[i].Blocks = make([]block.XYB, 0, 8)
			for j := range piece.NextBlockCount {
				info.Pieces[i].NextPieces[j].Blocks = make([]block.XYB, 0, 8)
			}
		}
		return info
	},
}

func ObtainRenderInfo() *RenderInfo {
	info := syncPoolRenderInfo.Get().(*RenderInfo)
	return info
}

func ReturnRenderInfo(info *RenderInfo) {
	if info == nil {
		return
	}

	syncPoolRenderInfo.Put(info)
}

func (f *Field) FillRenderInfo(info *RenderInfo, now time.Time) {
	w := f.w
	h := f.h

	showNextPieces := f.mode == ModeNormal
	showBlocks := f.mode != ModePause && f.mode != ModeSuspended && f.mode != ModeServerLost

	// reset the RenderInfo

	info.W = w
	info.H = h
	info.Mode = f.mode
	info.Blocks = info.Blocks[:0] // empty it, but keep the memory
	info.Result = f.animList.Process(now)

	info.TextData = TextData{
		BlocksRemoved: f.stats.blocksRemovedStr,
	}

	// process all blocks of the Field

	if showBlocks {
		idx := 0
		for y := range h {
			for x := range w {
				b := &f.blocks[idx]
				idx++

				if b.Type == block.TypeEmpty {
					continue
				}

				info.Blocks = append(info.Blocks, BlockRenderInfo{
					XYB: block.XYB{
						XY:    block.XY{X: x, Y: y},
						Block: b.Block,
					},
					Result: b.List.Process(now),
				})
			}
		}

		// process all external blocks of the Field

		var prev, curr *exElem
		curr = f.firstEx
		for curr != nil {
			next := curr.next

			result := curr.List.Process(now)
			if result.Feature == 0 {
				if prev == nil {
					f.firstEx = next
				} else {
					prev.next = next
				}

				curr.next = nil
				curr = next

				continue
			}

			info.Blocks = append(info.Blocks, BlockRenderInfo{
				XYB: block.XYB{
					XY:    curr.XY,
					Block: curr.Block,
				},
				Result: result,
			})

			prev = curr
			curr = next
		}
	}

	// process each Piece of the Field

	pieceCount := f.Ctrls()

	for i := pieceCount; i < len(info.Pieces); i++ {
		info.Pieces[i].Position = piece.DisplayPositionOff
		info.Pieces[i].State = 0
		info.Pieces[i].IsLimited = false
		info.Pieces[i].PieceEmpty = true
		for j := range piece.NextBlockCount {
			info.Pieces[i].NextPieces[j].Type = piece.TypeNone
			info.Pieces[i].NextPieces[j].Blocks = info.Pieces[i].NextPieces[j].Blocks[:0]
		}
	}

	for pIdx := range pieceCount {
		ctrl := f.pieces[pIdx]

		pinfo := &info.Pieces[pIdx]
		pinfo.PieceEmpty = false

		pinfo.Position = ctrl.InfoPosition
		pinfo.State = ctrl.State

		pinfo.PieceTextData = PieceTextData{
			Name:     ctrl.Name,
			Score:    ctrl.ScoreStr,
			PieceNum: ctrl.PieceCountStr,
			Level:    ctrl.LevelStr,
		}
		pinfo.IsLimited = ctrl.IsColumnLimited
		pinfo.Limits = ctrl.ColumnLimit

		pinfo.Blocks = pinfo.Blocks[:0]
		pinfo.DrawShadow = false
		pinfo.DirectionCW = ctrl.Config.RotationDirectionCW

		if !showBlocks || ctrl.State.IsTerminal() {
			continue
		}

		if showNextPieces {
			for i := range piece.NextBlockCount {
				pinfo.NextPieces[i].Type = ctrl.NextPieces[i].Type
				pinfo.NextPieces[i].Blocks = append(pinfo.NextPieces[i].Blocks, ctrl.NextPieces[i].Blocks...)
			}
		}

		p := ctrl.Piece
		if p == nil {
			continue
		}

		pinfo.X = ctrl.X
		pinfo.Y = ctrl.Y
		pinfo.DimX = int(p.DimX())
		pinfo.DimY = int(p.DimY())
		pinfo.Type = p.Type()
		pinfo.ActCount = int(p.ActivationCount())
		pinfo.Result = ctrl.List.Process(now)

		pinfo.Blocks = append(pinfo.Blocks, ctrl.Blocks...)

		pinfo.DrawShadow = ctrl.IsShadowShown
		pinfo.Shadow = ctrl.Shadow
	}
}
