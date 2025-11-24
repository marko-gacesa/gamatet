// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"github.com/marko-gacesa/gamatet/game/block"
)

type DestroyBlockInfo struct {
	Row    int
	Height int
	N      int
	block.Block
}

type DestroyColumnInfo struct {
	Rows    []DestroyBlockInfo
	HasImm  bool
	HasHard bool
}

type DestroyInfo struct {
	RowCount int
	Columns  []DestroyColumnInfo
	HardDec  []block.XY
}

func (info DestroyInfo) GetSimpleRow() int {
	if info.RowCount != 1 || len(info.Columns) == 0 || len(info.HardDec) > 0 {
		return -1
	}

	var row, n int

	for i := 0; i < len(info.Columns); i++ {
		if len(info.Columns[i].Rows) != 1 {
			return -1
		}

		q := info.Columns[i].Rows[0]

		if q.Height != 1 || q.Type != block.TypeRock { // simple row is made of only TypeRock blocks
			return -1
		}

		if i == 0 {
			row = q.Row
			n = q.N
		} else if q.Row != row || q.N != n {
			return -1
		}
	}

	return row
}

func (info DestroyInfo) HasHardOrImm() bool {
	if info.RowCount == 0 {
		return false
	}
	if len(info.HardDec) > 0 {
		return true
	}
	for _, col := range info.Columns {
		if col.HasImm || col.HasHard {
			return true
		}
	}
	return false
}
