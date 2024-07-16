// Copyright (c) 2020-2024 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/op"
	"gamatet/game/piece"
)

var _ Sweeper = (*Row)(nil)

func NewRow(f *field.Field) *Row {
	b := newBase(f)
	return &Row{base: *b}
}

type Row struct{ base }

func (s *Row) Start(analyzer op.Analyzer) {
	if analyzer.HasAdded {
		s.base.Start(analyzer)
	}
}

// Sweep removes blocks from the field. Removed blocks are from full rows.
// It returns number of blocks that are removed as the result of the function call.
func (s *Row) Sweep(p event.Pusher) {
	f := s.field

	result := f.GetDestroyInfo()
	if result.RowCount == 0 {
		s.endIteration()
		return
	}

	// Examine if this is a trivial case: Only one row is destroyed,
	// all columns are the same and there are no issues with block hardness or immovable blocks.
	row := result.GetSimpleRow()
	if row > -1 {
		blocks := f.GetRow(row)
		p.Push(op.NewFieldDestroyRow(row, blocks))
		s.endIteration()
		return
	}

	// First, must process hardness decrement
	for _, loc := range result.HardDec {
		p.Push(op.NewFieldBlockHardness(loc.X, loc.Y, -1, field.AnimSpin, 0))
	}

	var maxHeight int

	// Process each of the columns, from the bottom row up to the top (rows are returned in that order).
	for col, columnInfo := range result.Columns {
		for _, blockInfo := range columnInfo.Rows {
			b := f.GetXY(col, blockInfo.Row)
			p.Push(op.NewFieldDestroyColumn(col, blockInfo.Row, blockInfo.N, blockInfo.Height, b))
			if blockInfo.Height > maxHeight {
				maxHeight = blockInfo.Height
			}
		}
	}

	// A weird situation can happen after blocks destruction when field contains immovable or hard blocks:
	// A full row can form as a result of block destruction. That why we do another run of detection.
	//
	// Example with Hardness>0:              Example with Immovable block:
	// |            |    |            |      |   [0][0][0]|    |            |
	// |[0][0][0]   | => |            |      |[I]         | => |[I][0][0][0]|
	// |[0][0][0][2]|    |[0][0][0][1]|      |[0][0][0][0]|    |            |

	if !result.HasHardOrImm() {
		s.endIteration()
		return
	}

	s.reschedule(piece.GetFallDuration(maxHeight) + piece.DurationFullLine)
}
