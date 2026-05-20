package field

import (
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/logic/random"
)

type Reader interface {
	GetMode() Mode
	IsFinished() bool

	GetWidth() int
	GetHeight() int

	Random(seed uint64) *random.Random

	GetXY(x, y int) (b block.Block)
	GetRow(y int) []block.Block

	Ctrls() int
	CtrlLevel(idx byte) uint
	CtrlPieceCount(idx byte) uint
	CtrlPieceOverridden(ctrlIdx byte, pieceIdx uint) bool
	CtrlStateIsTerminal(ctrlIdx byte) bool
	CtrlWidth() int
	CtrlPlayerIndex(idx byte) byte

	CanMovePiece(dx, dy, pIdx int, liftAll bool) (success bool)
	CanRotatePiece(pIdx int, liftAll bool) (success bool, inverted bool, dx, dy int, rotated piece.Piece)
	CanFlipVPiece(pIdx int, liftAll bool) (success bool, flipped piece.Piece)
	CanFlipHPiece(pIdx int, liftAll bool) (success bool, flipped piece.Piece)
	GetDropHeight(pIdx int, liftAll bool) (height int)
	GetPieceBlockLocations(x, y int, p piece.Piece) (result []block.XYB)
	GetPieceStartPosition(pIdx int, ctrl *piece.Ctrl, p piece.Piece, liftAll bool) (success bool, x, y int)

	GetTopmostEmpty(x int) int
	GetTopmostFull(x int) int
	GetHeightToTopmostEmpty(x, y int) (height int)
	GetHeightToTopmostFull(x, y int) (height int)
	GetHeightToHighestHole(x, y int) (height int)
	GetHeightToLowestHole(x, y int) (height int)
	GetDestroyInfo() (info DestroyInfo)

	RangeBlocks(inspect func(xyb block.XYB) bool)
	FindBlizzardTops() []block.XY
	FindAcidRainTops() []block.XY
	FindMovableColumnSections(col int, filter func(*Field, ColumnSection) bool) []ColumnSection
	FindMovableSections(filter func(*Field, ColumnSection) bool) []ColumnSection

	HasLOS(p0, p1 block.XY) bool
	Neighbors8(pos block.XY, fnOk func(block.Type) bool) Neighbors8
	Neighbors4(pos block.XY, fnOk func(block.Type) bool) Neighbors4
	Path4(start, goal block.XY, fnOk func(block.Type) bool) []block.XY
	Path8(start, goal block.XY, fnOk func(block.Type) bool) []block.XY
	FindNearest8(pos block.XY, r int, fn func(block.XYB, int) bool) (block.XYB, bool)
	FindNearest4(pos block.XY, r int, fn func(block.XYB, int) bool) (block.XYB, bool)

	Blizzard(intensity int) []block.XY
	GetRandomBlock() (block.XYB, bool)

	SpawnLocation(loc SpawnLocation) (block.XY, bool)

	GetBlocksRemoved() int
	GetEffect() (Effect, byte)
	GetLingering(effect Effect) int
	LingeringEffects(fn func(effect Effect, amount int))
	LingeringEffect() (Effect, int)
}
