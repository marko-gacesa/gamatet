// Copyright (c) 2020 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/op"
	"gamatet/game/piece"
	"time"
)

func NewFullRowSweeper(f *field.Field) Sweeper {
	l := &fullRowSweeper{}
	l.field = f
	l.timer = time.NewTimer(time.Second)
	l.timer.Stop()
	return l
}

type fullRowSweeper struct {
	field  *field.Field
	timer  *time.Timer
	active bool
}

func (l *fullRowSweeper) Field() *field.Field {
	return l.field
}

func (l *fullRowSweeper) Timer() <-chan time.Time {
	return l.timer.C
}

func (l *fullRowSweeper) Start() {
	if l.active {
		// timer is already active or nothing to do
		return
	}

	l.active = true
	l.timer.Reset(time.Microsecond)
}

func (l *fullRowSweeper) Pause() {
	if !l.active {
		return
	}

	l.timer.Stop()
	select {
	default:
	case <-l.timer.C:
	}
}

func (l *fullRowSweeper) Unpause() {
	if !l.active {
		return
	}

	l.timer.Reset(time.Millisecond)
}

// Sweep removes blocks from the field. Removed blocks are from full rows.
// It returns number of blocks that are removed as the result of the function call.
func (l *fullRowSweeper) Sweep(p event.Pusher) {
	f := l.field

	result := f.GetDestroyInfo()
	if result.RowCount == 0 {
		l.endIteration()
		return
	}

	// Examine if this is a trivial case: Only one row is destroyed,
	// all columns are the same and there are no issues with block hardness or immovable blocks.
	row := result.GetSimpleRow()
	if row > -1 {
		blocks := f.GetRow(row)
		p.Push(op.NewFieldDestroyRow(row, blocks))
		l.endIteration()
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
		l.endIteration()
		return
	}

	l.timer.Reset(piece.GetFallDuration(maxHeight) + piece.DurationFullLine)
}

func (l *fullRowSweeper) endIteration() {
	l.active = false
}
