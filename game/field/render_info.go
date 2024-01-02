// Copyright (c) 2020 by Marko Gaćeša

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

type PieceRenderInfo struct {
	Empty      bool
	Blocks     []block.XYB
	X, Y       int
	DimX, DimY int
	Type       piece.Type
	Result     anim.Result
	DrawShadow bool
	Shadow     piece.Shadow

	IsLimited bool
	Limits    piece.ColumnLimit

	NextBlocks [piece.NextBlockCount][]block.XYB
}

type RenderInfo struct {
	W, H int

	Blocks   []BlockRenderInfo
	BlockCnt int

	Pieces   [4]PieceRenderInfo
	PieceCnt int

	Result anim.Result
}

var syncPoolRenderInfo = &sync.Pool{
	New: func() interface{} {
		info := &RenderInfo{}
		info.Blocks = make([]BlockRenderInfo, 0, 256)
		for i := 0; i < len(info.Pieces); i++ {
			info.Pieces[i].Blocks = make([]block.XYB, 0, 8)
			for j := 0; j < piece.NextBlockCount; j++ {
				info.Pieces[i].NextBlocks[j] = make([]block.XYB, 0, 8)
			}
		}
		return info
	},
}

func ReturnRenderInfo(info *RenderInfo) {
	if info == nil {
		return
	}

	syncPoolRenderInfo.Put(info)
}

func (f *Field) GetRenderInfo(now time.Time) *RenderInfo {
	w := f.w
	h := f.h

	info := syncPoolRenderInfo.Get().(*RenderInfo)

	// reset the RenderInfo

	info.Blocks = info.Blocks[:0] // empty it, but keep the memory
	info.BlockCnt = 0

	info.W = w
	info.H = h

	// process all blocks of the Field

	idx := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
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
			info.BlockCnt++
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
		info.BlockCnt++

		prev = curr
		curr = next
	}

	// process each Piece of the Field

	for i := len(f.pieces); i < len(info.Pieces); i++ {
		info.Pieces[i].Empty = true
		info.Pieces[i].IsLimited = false
	}

	for pIdx := 0; pIdx < len(f.pieces); pIdx++ {
		ctrl := f.pieces[pIdx]

		pinfo := &info.Pieces[pIdx]

		pinfo.IsLimited = ctrl.IsColumnLimited
		pinfo.Limits = ctrl.ColumnLimit

		p := ctrl.Piece

		if p == nil {
			pinfo.Empty = true
			continue
		}

		pinfo.Empty = false
		info.PieceCnt++

		dw := p.DimX()
		dh := p.DimY()

		pinfo.X = ctrl.X
		pinfo.Y = ctrl.Y
		pinfo.DimX = dw
		pinfo.DimY = dh
		pinfo.Type = ctrl.Piece.Type()
		pinfo.Result = ctrl.List.Process(now)

		pinfo.Blocks = pinfo.Blocks[:0]
		pinfo.Blocks = append(pinfo.Blocks, ctrl.Blocks...)

		pinfo.DrawShadow = ctrl.IsShadowShown
		pinfo.Shadow = ctrl.Shadow

		for i := 0; i < piece.NextBlockCount; i++ {
			pinfo.NextBlocks[i] = pinfo.NextBlocks[i][:0]
			pinfo.NextBlocks[i] = append(pinfo.NextBlocks[i], ctrl.NextBlocks[i]...)
		}
	}

	return info
}
