// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
)

var _ Sweeper = (*Row)(nil)

func NewRow(f *field.Field) *Row {
	b := newBase(f)
	return &Row{base: *b}
}

type Row struct {
	base
	countRemoved  int
	countSoftened int
}

func (s *Row) Start(analyzer *Analyzer) (started bool) {
	if analyzer.blocks.added > 0 {
		started = s.base.Start(analyzer)
		if started {
			s.countRemoved = 0
			s.countSoftened = 0
		}
	}
	return
}

// Sweep removes blocks from the field. Removed blocks are from full rows.
// It returns number of blocks that are removed as the result of the function call.
func (s *Row) Sweep(p event.Pusher) {
	f := s.field

	result := f.GetDestroyInfo()
	if result.RowCount == 0 {
		s.finish(p)
		return
	}

	// Examine if this is a trivial case: Only one row is destroyed,
	// all columns are the same and there are no issues with block hardness or immovable blocks.
	row := result.GetSimpleRow()
	if row > -1 {
		blocks := f.GetRow(row)
		p.Push(op.NewFieldDestroyRow(row, blocks))
		s.countRemoved += f.GetWidth()
		s.finish(p)
		return
	}

	// First, must process hardness decrement
	for _, loc := range result.HardDec {
		p.Push(op.NewFieldBlockHardness(loc.X, loc.Y, -1, field.AnimSpin, 0))
	}

	s.countSoftened += len(result.HardDec)

	var maxHeight int

	// Process each of the columns, from the bottom row up to the top (rows are returned in that order).
	for col, columnInfo := range result.Columns {
		for _, blockInfo := range columnInfo.Rows {
			b := f.GetXY(col, blockInfo.Row)
			s.countRemoved++
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
		s.finish(p)
		return
	}

	s.reschedule(piece.GetFallDuration(maxHeight) + piece.DurationFullLine)
}

func (s *Row) finish(p event.Pusher) {
	s.endIteration()

	if s.countRemoved == 0 && s.countSoftened == 0 {
		return
	}

	p.Push(op.NewFieldStat(int16(s.countRemoved), int16(s.countSoftened)))

	s.countRemoved = 0
	s.countSoftened = 0
}
