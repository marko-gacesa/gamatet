// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package setup

type GameType byte

const (
	GameTypeFallingPolyominoes = iota
)

type FieldInit byte // max value 31

func (q FieldInit) String() string {
	return FieldInitStringMap[q]
}

const (
	FieldInitEmpty            FieldInit = 0
	FieldInitLowSparseBlocks  FieldInit = 10
	FieldInitLowDenseBlocks   FieldInit = 11
	FieldInitHighSparseBlocks FieldInit = 12
	FieldInitHighDenseBlocks  FieldInit = 13
	FieldInitFunnel           FieldInit = 20
	FieldInitTriangle         FieldInit = 21
)

var FieldInitStringMap = map[FieldInit]string{
	FieldInitEmpty:            "empty",
	FieldInitLowSparseBlocks:  "low_sparse_blocks",
	FieldInitLowDenseBlocks:   "low_dense_blocks",
	FieldInitHighSparseBlocks: "high_sparse_blocks",
	FieldInitHighDenseBlocks:  "high_dense_blocks",
	FieldInitFunnel:           "funnel",
	FieldInitTriangle:         "triangle",
}

var FieldInits = []FieldInit{
	FieldInitEmpty,
	FieldInitLowSparseBlocks,
	FieldInitLowDenseBlocks,
	FieldInitHighSparseBlocks,
	FieldInitHighDenseBlocks,
	FieldInitFunnel,
	FieldInitTriangle,
}
