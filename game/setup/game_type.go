// Copyright (c) 2025 by Marko Gaćeša

package setup

type GameType byte

const (
	GameTypeRotatingPolyominoes GameType = iota
	GameTypeHMirroringPolyominoes
	GameTypeVMirroringPolyominoes
)

var GameTypeAll = []GameType{
	GameTypeRotatingPolyominoes,
	GameTypeHMirroringPolyominoes,
	GameTypeVMirroringPolyominoes,
}

var GameTypeNameMap = map[GameType]string{
	GameTypeRotatingPolyominoes:   "Rotating Polyominoes",
	GameTypeHMirroringPolyominoes: "Horizontal Mirroring Polyominoes",
	GameTypeVMirroringPolyominoes: "Vertical Mirroring Polyominoes",
}
