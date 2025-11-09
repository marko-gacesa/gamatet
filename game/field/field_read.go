// Copyright (c) 2020-2025 by Marko Gaćeša

package field

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/piece"
	"strings"
)

func (f *Field) GetWidth() int {
	return f.w
}

func (f *Field) GetHeight() int {
	return f.h
}

func (f *Field) GetXY(x, y int) (b block.Block) {
	b, _ = f.getXY(x, y)
	return
}

func (f *Field) GetRow(y int) []block.Block {
	blocks := make([]block.Block, f.w)

	idx := f.w * y
	for x := 0; x < f.w; x++ {
		blocks[x] = f.blocks[idx+x].Block
	}

	return blocks
}

func (f *Field) _getColumnLimits(ctrl *piece.Ctrl) (colMin, colMax int) {
	if ctrl.IsColumnLimited {
		colMin = ctrl.ColumnLimit.Min
		colMax = ctrl.ColumnLimit.Max
		return
	}

	colMin = 0
	colMax = f.w - 1
	return
}

func (f *Field) CanMovePiece(dx, dy, pIdx int, liftAll bool) (success bool) {
	c := f.pieces[pIdx]
	if c.Piece == nil {
		return
	}

	colMin, colMax := f._getColumnLimits(c)
	success = f._canPlacePiece(c.X+dx, c.Y+dy, colMin, colMax, c.Piece, liftAll, pIdx)
	return
}

func (f *Field) CanRotatePiece(pIdx int, liftAll bool) (success bool, inverted bool, dx, dy int, rotated piece.Piece) {
	c := f.pieces[pIdx]
	if c.Piece == nil || c.Piece.Type() != piece.TypeRotation {
		return
	}

	rotated = c.Piece.Clone()

	if c.RotationDirectionCW {
		inverted = rotated.UndoActivate()
	} else {
		inverted = rotated.Activate()
	}

	colMin, colMax := f._getColumnLimits(c)

	success = f._canPlacePiece(c.X, c.Y, colMin, colMax, rotated, liftAll, pIdx)
	if success {
		return
	}

	pieceWallKick := min(c.Piece.WallKick(), c.Config.WallKick)

	// Wall kick left/right

	for wallKick := 1; wallKick <= int(pieceWallKick); wallKick++ {
		success = f._canPlacePiece(c.X+wallKick, c.Y, colMin, colMax, rotated, liftAll, pIdx)
		if success {
			dx = wallKick
			return
		}

		success = f._canPlacePiece(c.X-wallKick, c.Y, colMin, colMax, rotated, liftAll, pIdx)
		if success {
			dx = -wallKick
			return
		}
	}

	// Ceiling kick down

	if c.Y >= f.h {
		for wallKick := 1; wallKick <= int(pieceWallKick); wallKick++ {
			success = f._canPlacePiece(c.X, c.Y-wallKick, colMin, colMax, rotated, liftAll, pIdx)
			if success {
				dy = -wallKick
				return
			}
		}
	}

	return
}

func (f *Field) CanFlipVPiece(pIdx int, liftAll bool) (success bool, flipped piece.Piece) {
	c := f.pieces[pIdx]
	if c.Piece == nil || c.Piece.Type() != piece.TypeFlipV {
		return
	}

	flipped = c.Piece.Clone()
	flipped.Activate()

	colMin, colMax := f._getColumnLimits(c)
	success = f._canPlacePiece(c.X, c.Y, colMin, colMax, flipped, liftAll, pIdx)
	return
}

func (f *Field) CanFlipHPiece(pIdx int, liftAll bool) (success bool, flipped piece.Piece) {
	c := f.pieces[pIdx]
	if c.Piece == nil || c.Piece.Type() != piece.TypeFlipH {
		return
	}

	flipped = c.Piece.Clone()
	flipped.Activate()

	colMin, colMax := f._getColumnLimits(c)
	success = f._canPlacePiece(c.X, c.Y, colMin, colMax, flipped, liftAll, pIdx)
	return
}

func (f *Field) GetDropHeight(pIdx int, liftAll bool) (height int) {
	c := f.pieces[pIdx]
	if c.Piece == nil {
		return
	}

	colMin, colMax := f._getColumnLimits(c)

	for {
		success := f._canPlacePiece(c.X, c.Y-height-1, colMin, colMax, c.Piece, liftAll, pIdx)
		if success {
			height++
			continue
		}

		break
	}

	return
}

// GetPieceBlockLocations returns all piece blocks on their actual location in the fields.
// Blocks are returned from the lowest row, up to the topmost row.
// The function is used to meld the piece blocks into the field.
func (f *Field) GetPieceBlockLocations(x, y int, p piece.Piece) (result []block.XYB) {
	result = make([]block.XYB, 0, p.BlockCount())

	dimX := int(p.DimX())
	dimY := int(p.DimY())
	for j := range dimY {
		for i := range dimX {
			pBlock := p.Get(i, j)
			if pBlock.Type == block.TypeEmpty {
				continue
			}

			idx := (y-j)*f.w + (x + i)
			fBlock := f.blocks[idx]

			if fBlock.Type != block.TypeEmpty {
				continue
			}

			result = append(result, block.XYB{
				XY:    block.XY{X: x + i, Y: y - j},
				Block: pBlock,
			})
		}
	}

	return
}

func (f *Field) GetPieceStartPosition(pIdx int, ctrl *piece.Ctrl, p piece.Piece, liftAll bool) (success bool, x, y int) {
	dimX := int(p.DimX())

	colMin, colMax := f._getColumnLimits(ctrl)

	if ctrl.IsColumnLimited {
		// example:
		// 01234567890123456789 column
		// ------OOOOOO-------- allowed (min=6, max=11)
		// .......XXX.......... piece
		// (11+6-3+1)/2 = (17-3+1)/2 = 15/2 = 7
		// 01234567890123456789 column
		// ------OOOOOO-------- allowed (min=6, max=11)
		// ........XX.......... piece
		// (11+6-2+1)/2 = (17-2+1)/2 = 16/2 = 8
		x = (colMax + colMin - dimX + 1) / 2
	} else {
		// example:
		// 01234567890123456789 column
		// 00000111112222233333 player
		// ...........XXX...... piece (x=11)
		// f.w=20, pIdx=2, pieceCount=4, dim=3 -> x = 11
		pieceCount := len(f.pieces)
		x = (f.w*pIdx)/pieceCount + (f.w/pieceCount-dimX)/2
	}

	y = f.h - 1
	y += int(p.TopEmptyRows())

	success = f._canPlacePiece(x, y, colMin, colMax, p, liftAll, pIdx)
	return
}

func (f *Field) GetTopmostEmpty(x int) int {
	w := f.w
	y := f.h - 1
	for idx := y*w + x; idx >= 0; idx -= w {
		if f.blocks[idx].Type != block.TypeEmpty {
			return y + 1
		}
		y--
	}
	return 0
}

// GetHeightToTopmostEmpty returns height from which the block at location (x, y) can fall
// until it hits something - either the bottom or another block. So, if height>0, the location
// (x, y - height) always will be empty. If height==0 than falling is not possible.
func (f *Field) GetHeightToTopmostEmpty(x, y int) (height int) {
	w := f.w
	for idx := (y-1)*w + x; idx >= 0; idx -= w {
		if f.blocks[idx].Type != block.TypeEmpty {
			break
		}
		height++
	}
	return
}

// GetHeightToTopmostFull returns height from which the block at location (x, y) can fall
// until it gets to the location of another block. So, if height>0, the location
// (x, y - height) will never be empty. If there are no blocks under, it returns height=0.
func (f *Field) GetHeightToTopmostFull(x, y int) (height int) {
	w := f.w
	for idx := (y-1)*w + x; idx >= 0; idx -= w {
		height++
		if f.blocks[idx].Type != block.TypeEmpty {
			return height
		}
	}

	return 0
}

// GetHeightToHighestHole returns height from which the block at location (x, y) can fall,
// tunnel through some of the existing blocks until it finds the first hole (empty block).
// So, if height>0, the location (x, y - height) will be empty. If height==0 than falling is not possible.
func (f *Field) GetHeightToHighestHole(x, y int) (height int) {
	height = f.GetHeightToTopmostFull(x, y)
	if height == 0 {
		return
	}

	height++

	w := f.w
	for idx := (y-height)*w + x; idx >= 0; idx -= w {
		if f.blocks[idx].Type == block.TypeEmpty {
			return height
		}
		height++
	}

	return 0
}

// GetHeightToLowestHole returns height from which the block at location (x, y) can fall,
// tunnel through all existing blocks until it finds the last hole (empty block).
// So, if height>0, the location (x, y - height) will be empty. If height==0 than falling is not possible.
func (f *Field) GetHeightToLowestHole(x, y int) (height int) {
	height = f.GetHeightToTopmostFull(x, y)
	if height == 0 {
		return
	}

	w := f.w

	yCurr := 0
	idxTop := (y - height) * w

	for idx := x; idx < idxTop; idx += w {
		if f.blocks[idx].Type == block.TypeEmpty {
			return y - yCurr
		}
		yCurr++
	}

	return 0
}

// GetDestroyInfo finds all blocks that should be destroyed (because they are in a full row).
// For each block it returns the block's Column, Row and Type.
// Also, it finds N, the number of blocks (or empty places) above the block that should fall
// and the Height from which these N blocks should fall.
// The algorithm recognizes block Hardness and Immovable block types.
func (f *Field) GetDestroyInfo() (info DestroyInfo) {
	w := f.w
	h := f.h

	blockIdx := 0 // blockIdx is block index and is used to avoid frequent row*w+col calculations

	colIsOpen := make([]bool, w)
	colCurrH := make([]int, w)

	// Proceed from the bottom row up to the top row: Find all rows with all blocks non-empty.
	for row := range h {
		isFull := true
		onlyHard := true
		for col := range w {
			b := f.blocks[blockIdx]

			if b.Type == block.TypeEmpty {
				// Found an empty block => the row is not full, move to the next row.
				blockIdx += w - col
				isFull = false
				break
			}

			if b.Hardness != block.HardnessMax {
				onlyHard = false
			}

			blockIdx++
		}

		// Found a full row
		if isFull && !onlyHard {
			blockIdx -= w

			info.RowCount++
			if len(info.Columns) == 0 {
				info.Columns = make([]DestroyColumnInfo, f.w)
			}

			for col := range w {
				b := f.blocks[blockIdx]
				blockIdx++

				if b.Hardness > 0 {
					info.Columns[col].HasHard = true

					if b.Hardness != block.HardnessMax {
						info.HardDec = append(info.HardDec, block.XY{X: col, Y: row})
					}

					colIsOpen[col] = false
					colCurrH[col] = 0

					continue
				}

				info.Columns[col].Rows = append(info.Columns[col].Rows, DestroyBlockInfo{
					Row:    row,
					Height: colCurrH[col] + 1,
					N:      0,
					Block:  b.Block,
				})
				colIsOpen[col] = true
				colCurrH[col]++
			}

			continue
		}

		// Found a non-full row

		if len(info.Columns) == 0 {
			continue
		}

		blockIdx -= w
		for col := range w {
			b := f.blocks[blockIdx]
			blockIdx++

			if b.Type.IsImmovable() {
				colIsOpen[col] = false
				colCurrH[col] = 0

				info.Columns[col].HasImm = true
				continue
			}

			if colIsOpen[col] {
				l := len(info.Columns[col].Rows)
				info.Columns[col].Rows[l-1].N++
			}
		}
	}

	return
}

func (f *Field) String() string {
	sb := strings.Builder{}
	for y := f.h - 1; y >= 0; y-- {
		for x := 0; x < f.w; x++ {
			pIdx := f._getXYPieceIdx(x, y)
			if pIdx >= 0 {
				sb.WriteByte('<')
				sb.WriteByte('A' + byte(pIdx))
				sb.WriteByte('>')
				continue
			}

			b, _ := f.getXY(x, y)
			if b.Type == block.TypeEmpty {
				sb.WriteString(" . ")
				continue
			}

			sb.WriteByte('[')
			sb.WriteByte('0' + b.Hardness)
			sb.WriteByte(']')
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
