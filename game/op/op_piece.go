// Copyright (c) 2020-2025 by Marko Gaćeša

package op

import (
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/game/serialize"
	"io"
	"strconv"
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

func NewPieceSet(pIdx int, op Type, x, y int, p piece.Piece, pCount int) *PieceSet {
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
	PieceCount int
}

var _ event.Event = (*PieceSet)(nil)

func (e *PieceSet) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	switch e.Op {
	case TypeSet:
		ctrl.SetXYP(int(e.X), int(e.Y), e.Piece)
		ctrl.PieceCount = e.PieceCount
		ctrl.PieceCountStr = strconv.Itoa(e.PieceCount)
		animateNewPiece(ctrl, f.Config.Anim)
		ctrl.Blocks = piece.GetBlocks(e.Piece, ctrl.Blocks[:0])
		for i := 0; i < piece.NextBlockCount; i++ {
			np := ctrl.Feed.Get(ctrl.PieceCount + i)
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
		ctrl.PieceCountStr = strconv.Itoa(e.PieceCount)
		ctrl.Blocks = piece.GetBlocks(e.Piece, ctrl.Blocks[:0])
		for i := 0; i < piece.NextBlockCount; i++ {
			np := ctrl.Feed.Get(ctrl.PieceCount + i)
			ctrl.NextPieces[i].Type = np.Type()
			ctrl.NextPieces[i].Blocks = piece.GetBlocks(np, ctrl.NextPieces[i].Blocks[:0])
		}
	}
	updatePieceShadow(f, ctrl)
}

func (e *PieceSet) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceSet)
	return ok && e.PieceIdx == e.PieceIdx &&
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

	e.PieceCount, err = serialize.ReadInt(r)
	if err != nil {
		return err
	}

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

func NewPieceActivate(pIdx int) *PieceActivate {
	return &PieceActivate{
		PieceIdx: byte(pIdx),
		Param:    0,
	}
}

func NewPieceRotate(pIdx int, dirCW bool) *PieceActivate {
	param := paramRotCCW
	if dirCW {
		param++
	}
	return &PieceActivate{
		PieceIdx: byte(pIdx),
		Param:    param,
	}
}

const (
	paramRotCCW byte = 0
	paramRotCW  byte = 1
)

type PieceActivate struct {
	PieceIdx byte
	Param    byte
}

var _ event.Event = (*PieceActivate)(nil)

func (e *PieceActivate) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	switch ctrl.Piece.Type() {
	case piece.TypeFlipV:
		ctrl.Piece.Activate()
		animateFlipVPiece(ctrl, f.Config.Anim)
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	case piece.TypeFlipH:
		ctrl.Piece.Activate()
		animateFlipHPiece(ctrl, f.Config.Anim)
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	case piece.TypeRotation:
		var inverted bool
		switch e.Param {
		case paramRotCCW:
			inverted = ctrl.Piece.Activate()
		case paramRotCW:
			inverted = ctrl.Piece.UndoActivate()
		default:
			return
		}
		animateRotatePiece(ctrl, e.Param == paramRotCW, inverted, f.Config.Anim)
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	case piece.TypeShooter:
		ctrl.Piece.Activate()
		if ctrl.Piece.ActivationCount() == 1 {
			animateBlinkPiece(ctrl, f.Config.Anim)
		}
	}
}

func (e *PieceActivate) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	switch ctrl.Piece.Type() {
	case piece.TypeFlipV, piece.TypeFlipH:
		ctrl.Piece.UndoActivate()
		ctrl.List.Clear()
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	case piece.TypeRotation:
		switch e.Param {
		case paramRotCCW:
			ctrl.Piece.UndoActivate()
		case paramRotCW:
			ctrl.Piece.Activate()
		default:
			return
		}
		ctrl.List.Clear()
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	case piece.TypeShooter:
		ctrl.Piece.UndoActivate()
		ctrl.List.Clear()
	}
}

func (e *PieceActivate) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceActivate)
	return ok && e.PieceIdx == q.PieceIdx
}

func (e *PieceActivate) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx, e.Param}); err != nil {
		return err
	}
	return nil
}

func (e *PieceActivate) Read(r io.Reader) error {
	var buffer [2]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.Param = buffer[1]

	return nil
}

func (e *PieceActivate) TypeID() event.Code { return codePieceActivate }

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
