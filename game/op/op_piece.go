// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package op

import (
	"encoding/binary"
	"io"
	"strconv"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/serialize"
)

func NewPieceState(pIdx int, oldState, newState piece.State, oldParam, newParam int) *PieceState {
	return &PieceState{
		PieceIdx: byte(pIdx),
		OldState: oldState,
		NewState: newState,
		OldParam: oldParam,
		NewParam: newParam,
	}
}

type PieceState struct {
	PieceIdx           byte
	OldState, NewState piece.State
	OldParam, NewParam int
}

var _ event.Event = (*PieceState)(nil)

func (e *PieceState) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.State = e.NewState
	ctrl.RestartTimer(e.NewParam)

	if e.NewState == piece.StateSlide && !ctrl.SlideDisabled {
		animateSlidePiece(ctrl, f.Config.Anim)
	}
}

func (e *PieceState) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.State = e.OldState
	ctrl.RestartTimer(e.OldParam)
}

func (e *PieceState) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceState)
	return ok && e.PieceIdx == q.PieceIdx &&
		e.OldState == q.OldState && e.NewState == q.NewState &&
		e.OldParam == q.OldParam && e.NewParam == q.NewParam
}

func (e *PieceState) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx, byte(e.OldState), byte(e.NewState)}); err != nil {
		return err
	}
	if err := serialize.Write32(w, uint32(e.OldParam)); err != nil {
		return err
	}
	if err := serialize.Write32(w, uint32(e.NewParam)); err != nil {
		return err
	}
	return nil
}

func (e *PieceState) Read(r io.Reader) error {
	var buffer [3]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.OldState = piece.State(buffer[1])
	e.NewState = piece.State(buffer[2])

	var err error

	e.OldParam, err = serialize.ReadInt(r)
	if err != nil {
		return err
	}

	e.NewParam, err = serialize.ReadInt(r)
	if err != nil {
		return err
	}

	return nil
}

func (e *PieceState) TypeID() event.Code { return codePieceState }

func NewPieceSet(pIdx int, op Type, x, y int, p piece.Piece, pCount uint) *PieceSet {
	return &PieceSet{
		PieceIdx:   byte(pIdx),
		Op:         op,
		X:          int8(x),
		Y:          int8(y),
		Piece:      p,
		PieceCount: pCount,
	}
}

type PieceSet struct {
	PieceIdx   byte
	Op         Type
	X, Y       int8
	Piece      piece.Piece
	PieceCount uint
}

var _ event.Event = (*PieceSet)(nil)

func (e *PieceSet) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	switch e.Op {
	case TypeSet:
		ctrl.SetXYP(int(e.X), int(e.Y), e.Piece)
		ctrl.PieceCount = e.PieceCount
		ctrl.PieceCountStr = strconv.Itoa(int(e.PieceCount))
		animateNewPiece(ctrl, f.Config.Anim)
		ctrl.Blocks = piece.GetBlocks(e.Piece, ctrl.Blocks[:0])
		for i := range piece.NextBlockCount {
			np := ctrl.Feed.Get(ctrl.PieceCount+i, ctrl.PlayerIndex)
			ctrl.NextPieces[i].Type = np.Type()
			ctrl.NextPieces[i].Blocks = piece.GetBlocks(np, ctrl.NextPieces[i].Blocks[:0])
		}
	case TypeClear:
		ctrl.SetXYP(0, 0, nil)
		ctrl.Blocks = ctrl.Blocks[:0]
	}
	updatePieceShadow(f, ctrl)
}

func (e *PieceSet) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	switch e.Op {
	case TypeSet:
		ctrl.SetXYP(0, 0, nil)
		ctrl.Blocks = ctrl.Blocks[:0]
	case TypeClear:
		ctrl.SetXYP(int(e.X), int(e.Y), e.Piece)
		ctrl.PieceCount = e.PieceCount
		ctrl.PieceCountStr = strconv.Itoa(int(e.PieceCount))
		ctrl.Blocks = piece.GetBlocks(e.Piece, ctrl.Blocks[:0])
		for i := range piece.NextBlockCount {
			np := ctrl.Feed.Get(ctrl.PieceCount+i, ctrl.PlayerIndex)
			ctrl.NextPieces[i].Type = np.Type()
			ctrl.NextPieces[i].Blocks = piece.GetBlocks(np, ctrl.NextPieces[i].Blocks[:0])
		}
	}
	updatePieceShadow(f, ctrl)
}

func (e *PieceSet) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceSet)
	return ok && e.PieceIdx == q.PieceIdx &&
		e.Op == q.Op && e.X == q.X && e.Y == q.Y && e.Piece.Equals(q.Piece) && e.PieceCount != q.PieceCount
}

func (e *PieceSet) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx, byte(e.Op), byte(e.X), byte(e.Y)}); err != nil {
		return err
	}
	if err := piece.Write(w, e.Piece); err != nil {
		return err
	}
	if err := serialize.Write32(w, uint32(e.PieceCount)); err != nil {
		return err
	}
	return nil
}

func (e *PieceSet) Read(r io.Reader) error {
	var buffer [4]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.Op = Type(buffer[1])
	e.X = int8(buffer[2])
	e.Y = int8(buffer[3])

	var err error
	e.Piece, err = piece.Read(r)
	if err != nil {
		return err
	}

	pieceCount, err := serialize.ReadInt(r)
	if err != nil {
		return err
	}

	e.PieceCount = uint(pieceCount)

	return nil
}

func (e *PieceSet) TypeID() event.Code { return codePieceSet }

func NewPieceMove(pIdx int, dx, dy int) *PieceMove {
	return &PieceMove{
		PieceIdx: byte(pIdx),
		DX:       int8(dx),
		DY:       int8(dy),
	}
}

type PieceMove struct {
	PieceIdx byte
	DX, DY   int8
}

var _ event.Event = (*PieceMove)(nil)

func (e *PieceMove) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	dx := int(e.DX)
	dy := int(e.DY)

	ctrl.X += dx
	ctrl.Y += dy
	animateMovePiece(ctrl, dx, dy, f.Config.Anim)

	updatePieceShadow(f, ctrl)
}

func (e *PieceMove) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	dx := int(e.DX)
	dy := int(e.DY)

	ctrl.X -= dx
	ctrl.Y -= dy

	ctrl.List.Clear()

	updatePieceShadow(f, ctrl)
}

func (e *PieceMove) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceMove)
	return ok && e.PieceIdx == q.PieceIdx && e.DX == q.DX && e.DY == q.DY
}

func (e *PieceMove) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx, byte(e.DX), byte(e.DY)}); err != nil {
		return err
	}
	return nil
}

func (e *PieceMove) Read(r io.Reader) error {
	var buffer [3]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.DX = int8(buffer[1])
	e.DY = int8(buffer[2])

	return nil
}

func (e *PieceMove) TypeID() event.Code { return codePieceMove }

func NewPieceRotate(pIdx int, dirCW bool) *PieceRotate {
	return &PieceRotate{
		PieceIdx: byte(pIdx),
		DirCW:    dirCW,
	}
}

type PieceRotate struct {
	PieceIdx byte
	DirCW    bool
}

var _ event.Event = (*PieceRotate)(nil)

func (e *PieceRotate) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	if ctrl.Piece.Type() != piece.TypeRotation {
		return
	}

	var inverted bool
	if e.DirCW {
		inverted = ctrl.Piece.UndoActivate()
	} else {
		inverted = ctrl.Piece.Activate()
	}

	animateRotatePiece(ctrl, e.DirCW, inverted, f.Config.Anim)
	ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
	updatePieceShadow(f, ctrl)
}

func (e *PieceRotate) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	if ctrl.Piece.Type() != piece.TypeRotation {
		return
	}

	if e.DirCW {
		ctrl.Piece.Activate()
	} else {
		ctrl.Piece.UndoActivate()
	}

	ctrl.List.Clear()
	ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
	updatePieceShadow(f, ctrl)
}

func (e *PieceRotate) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceRotate)
	return ok && e.PieceIdx == q.PieceIdx && e.DirCW == q.DirCW
}

func (e *PieceRotate) Write(w io.Writer) error {
	var dirCW byte
	if e.DirCW {
		dirCW = 1
	}

	if _, err := w.Write([]byte{e.PieceIdx, dirCW}); err != nil {
		return err
	}
	return nil
}

func (e *PieceRotate) Read(r io.Reader) error {
	var buffer [2]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.DirCW = buffer[1] != 0

	return nil
}

func (e *PieceRotate) TypeID() event.Code { return codePieceRotate }

func NewPieceFlip(pIdx int) *PieceFlip {
	return &PieceFlip{
		PieceIdx: byte(pIdx),
	}
}

type PieceFlip struct {
	PieceIdx byte
}

var _ event.Event = (*PieceFlip)(nil)

func (e *PieceFlip) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	if pt := ctrl.Piece.Type(); pt != piece.TypeFlipV && pt != piece.TypeFlipH {
		return
	}

	ctrl.Piece.Activate()

	if ctrl.Piece.Type() == piece.TypeFlipV {
		animateFlipVPiece(ctrl, f.Config.Anim)
	} else {
		animateFlipHPiece(ctrl, f.Config.Anim)
	}

	ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
	updatePieceShadow(f, ctrl)
}

func (e *PieceFlip) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	if pt := ctrl.Piece.Type(); pt != piece.TypeFlipV && pt != piece.TypeFlipH {
		return
	}

	ctrl.Piece.UndoActivate()
	ctrl.List.Clear()
	ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
	updatePieceShadow(f, ctrl)
}

func (e *PieceFlip) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceFlip)
	return ok && e.PieceIdx == q.PieceIdx
}

func (e *PieceFlip) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx}); err != nil {
		return err
	}
	return nil
}

func (e *PieceFlip) Read(r io.Reader) error {
	var buffer [1]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]

	return nil
}

func (e *PieceFlip) TypeID() event.Code { return codePieceFlip }

func NewPieceShoot(pIdx int, hit bool, t block.Type) *PieceShoot {
	return &PieceShoot{
		PieceIdx:  byte(pIdx),
		Hit:       hit,
		BlockType: t,
	}
}

type PieceShoot struct {
	PieceIdx  byte
	Hit       bool
	BlockType block.Type
}

var _ event.Event = (*PieceShoot)(nil)

func (e *PieceShoot) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	if ctrl.Piece.Type() != piece.TypeShooter {
		return
	}

	ctrl.Piece.Activate()
	if ctrl.Piece.ActivationCount() == 1 {
		animateBlinkPiece(ctrl, f.Config.Anim)
	}
}

func (e *PieceShoot) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	if ctrl.Piece.Type() != piece.TypeShooter {
		return
	}

	ctrl.Piece.UndoActivate()
	ctrl.List.Clear()
}

func (e *PieceShoot) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceShoot)
	return ok && e.PieceIdx == q.PieceIdx && e.Hit == q.Hit && e.BlockType == q.BlockType
}

func (e *PieceShoot) Write(w io.Writer) error {
	var hit byte
	if e.Hit {
		hit = 1
	}

	if _, err := w.Write([]byte{e.PieceIdx, hit, byte(e.BlockType)}); err != nil {
		return err
	}

	return nil
}

func (e *PieceShoot) Read(r io.Reader) error {
	var buffer [3]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.Hit = buffer[1] != 0
	e.BlockType = block.Type(buffer[2])

	return nil
}

func (e *PieceShoot) TypeID() event.Code { return codePieceShoot }

func NewPieceFall(pIdx int, height int) *PieceFall {
	return &PieceFall{
		PieceIdx: byte(pIdx),
		Height:   byte(height),
	}
}

type PieceFall struct {
	PieceIdx byte
	Height   byte
}

var _ event.Event = (*PieceFall)(nil)

func (e *PieceFall) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	height := int(e.Height)
	ctrl.Y -= height
	animateDropPiece(ctrl, height, f.Config.Anim)
	updatePieceShadow(f, ctrl)
}

func (e *PieceFall) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	height := int(e.Height)
	ctrl.Y += height
	ctrl.List.Clear()
	updatePieceShadow(f, ctrl)
}

func (e *PieceFall) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceFall)
	return ok && e.PieceIdx == q.PieceIdx && e.Height == q.Height
}

func (e *PieceFall) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx, e.Height}); err != nil {
		return err
	}
	return nil
}

func (e *PieceFall) Read(r io.Reader) error {
	var buffer [2]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.Height = buffer[1]

	return nil
}

func (e *PieceFall) TypeID() event.Code { return codePieceFall }

func NewPieceScore(pIdx int, delta int) *PieceScore {
	return &PieceScore{
		PieceIdx: byte(pIdx),
		Delta:    delta,
	}
}

type PieceScore struct {
	PieceIdx byte
	Delta    int
}

var _ event.Event = (*PieceScore)(nil)

func (e *PieceScore) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.ModifyScore(e.Delta)
}

func (e *PieceScore) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.ModifyScore(-e.Delta)
}

func (e *PieceScore) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceScore)
	return ok && e.PieceIdx == q.PieceIdx && e.Delta == q.Delta
}

func (e *PieceScore) Write(w io.Writer) error {
	var buffer [5]byte

	buffer[0] = e.PieceIdx
	binary.LittleEndian.PutUint32(buffer[1:5], uint32(e.Delta))

	if _, err := w.Write(buffer[:]); err != nil {
		return err
	}

	return nil
}

func (e *PieceScore) Read(r io.Reader) error {
	var buffer [5]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.Delta = int(binary.LittleEndian.Uint32(buffer[1:5]))

	return nil
}

func (e *PieceScore) TypeID() event.Code { return codePieceScore }

func NewPieceLevelBoost(pIdx int, boost bool) *PieceLevelBoost {
	return &PieceLevelBoost{
		PieceIdx: byte(pIdx),
		Boost:    boost,
	}
}

type PieceLevelBoost struct {
	PieceIdx byte
	Boost    bool
}

var _ event.Event = (*PieceLevelBoost)(nil)

func (e *PieceLevelBoost) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.LevelBoost = e.Boost
}

func (e *PieceLevelBoost) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.LevelBoost = e.Boost
}

func (e *PieceLevelBoost) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceLevelBoost)
	return ok && e.PieceIdx == q.PieceIdx && e.Boost == q.Boost
}

func (e *PieceLevelBoost) Write(w io.Writer) error {
	var b byte
	if e.Boost {
		b++
	}
	if _, err := w.Write([]byte{e.PieceIdx, b}); err != nil {
		return err
	}
	return nil
}

func (e *PieceLevelBoost) Read(r io.Reader) error {
	var buffer [2]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.Boost = buffer[1] != 0

	return nil
}

func (e *PieceLevelBoost) TypeID() event.Code { return codePieceLevelBoost }

func NewPieceOverride(pIdx byte, p piece.Piece, pCount uint) *PieceOverride {
	return &PieceOverride{
		PieceIdx:   pIdx,
		Piece:      p,
		PieceCount: pCount,
	}
}

type PieceOverride struct {
	PieceIdx   byte
	Piece      piece.Piece
	PieceCount uint
}

var _ event.Event = (*PieceOverride)(nil)

func (e *PieceOverride) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Feed.Override(e.PieceCount, e.Piece)
	for i := range piece.NextBlockCount {
		np := ctrl.Feed.Get(ctrl.PieceCount+i, ctrl.PlayerIndex)
		ctrl.NextPieces[i].Type = np.Type()
		ctrl.NextPieces[i].Blocks = piece.GetBlocks(np, ctrl.NextPieces[i].Blocks[:0])
	}
}

func (e *PieceOverride) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Feed.OverrideClear(e.PieceCount)
	for i := range piece.NextBlockCount {
		np := ctrl.Feed.Get(ctrl.PieceCount+i, ctrl.PlayerIndex)
		ctrl.NextPieces[i].Type = np.Type()
		ctrl.NextPieces[i].Blocks = piece.GetBlocks(np, ctrl.NextPieces[i].Blocks[:0])
	}
}

func (e *PieceOverride) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceOverride)
	return ok && e.PieceIdx == q.PieceIdx &&
		e.Piece.Equals(q.Piece) && e.PieceCount != q.PieceCount
}

func (e *PieceOverride) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx}); err != nil {
		return err
	}
	if err := piece.Write(w, e.Piece); err != nil {
		return err
	}
	if err := serialize.Write32(w, uint32(e.PieceCount)); err != nil {
		return err
	}
	return nil
}

func (e *PieceOverride) Read(r io.Reader) error {
	var buffer [1]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]

	var err error
	e.Piece, err = piece.Read(r)
	if err != nil {
		return err
	}

	pieceCount, err := serialize.ReadInt(r)
	if err != nil {
		return err
	}

	e.PieceCount = uint(pieceCount)

	return nil
}

func (e *PieceOverride) TypeID() event.Code { return codePieceOverride }

func NewPieceSpeedUp(pIdx byte, delta int8) *PieceSpeedUp {
	return &PieceSpeedUp{
		PieceIdx: pIdx,
		Delta:    delta,
	}
}

type PieceSpeedUp struct {
	PieceIdx byte
	Delta    int8
}

var _ event.Event = (*PieceSpeedUp)(nil)

func (e *PieceSpeedUp) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.SetLevel(uint(int(ctrl.Level) + int(e.Delta)))
	f.UpdateBlocksRemoved(0)
}

func (e *PieceSpeedUp) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.SetLevel(uint(int(ctrl.Level) - int(e.Delta)))
	f.UpdateBlocksRemoved(0)
}

func (e *PieceSpeedUp) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceSpeedUp)
	return ok && e.PieceIdx == q.PieceIdx
}

func (e *PieceSpeedUp) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx, byte(e.Delta)}); err != nil {
		return err
	}
	return nil
}

func (e *PieceSpeedUp) Read(r io.Reader) error {
	var buffer [2]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.Delta = int8(buffer[1])

	return nil
}

func (e *PieceSpeedUp) TypeID() event.Code { return codePieceSpeedUp }
