// Copyright (c) 2020 by Marko Gaćeša

package op

import (
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/util"
	"io"
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

	if e.NewState == piece.StateSlide && ctrl.SlideEnabled {
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
	if err := util.Write32(w, uint32(e.OldParam)); err != nil {
		return err
	}
	if err := util.Write32(w, uint32(e.NewParam)); err != nil {
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

	e.OldParam, err = util.ReadInt(r)
	if err != nil {
		return err
	}

	e.NewParam, err = util.ReadInt(r)
	if err != nil {
		return err
	}

	return nil
}

func NewPieceSet(pIdx int, op OpType, x, y int, p piece.Piece, pCount int) *PieceSet {
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
	Op         OpType
	X, Y       int8
	Piece      piece.Piece
	PieceCount int
}

var _ event.Event = (*PieceSet)(nil)

func (e *PieceSet) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	switch e.Op {
	case OpSet:
		ctrl.SetXYP(int(e.X), int(e.Y), e.Piece)
		ctrl.PieceCount = e.PieceCount
		animateNewPiece(ctrl, f.Config.Anim)
		ctrl.Blocks = piece.GetBlocks(e.Piece, ctrl.Blocks[:0])
		for i := 0; i < piece.NextBlockCount; i++ {
			ctrl.NextBlocks[i] = piece.GetBlocks(ctrl.Feed.Get(ctrl.PieceCount+i), ctrl.NextBlocks[i][:0])
		}
	case OpClear:
		ctrl.SetXYP(0, 0, nil)
		ctrl.Blocks = ctrl.Blocks[:0]
	}
	updatePieceShadow(f, ctrl)
}

func (e *PieceSet) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	switch e.Op {
	case OpSet:
		ctrl.SetXYP(0, 0, nil)
		ctrl.Blocks = ctrl.Blocks[:0]
	case OpClear:
		ctrl.SetXYP(int(e.X), int(e.Y), e.Piece)
		ctrl.PieceCount = e.PieceCount
		ctrl.Blocks = piece.GetBlocks(e.Piece, ctrl.Blocks[:0])
		for i := 0; i < piece.NextBlockCount; i++ {
			ctrl.NextBlocks[i] = piece.GetBlocks(ctrl.Feed.Get(ctrl.PieceCount+i), ctrl.NextBlocks[i][:0])
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
	if err := util.Write32(w, uint32(e.PieceCount)); err != nil {
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
	e.Op = OpType(buffer[1])
	e.X = int8(buffer[2])
	e.Y = int8(buffer[3])

	var err error
	e.Piece, err = piece.Read(r)
	if err != nil {
		return err
	}

	e.PieceCount, err = util.ReadInt(r)
	if err != nil {
		return err
	}

	return nil
}

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

func NewPieceRotate(pIdx int, cw bool) *PieceTransform {
	var rotate int8
	if cw {
		rotate = 1
	} else {
		rotate = -1
	}
	return &PieceTransform{
		PieceIdx: byte(pIdx),
		RotateCW: rotate,
		Activate: 0,
	}
}

func NewPieceActivate(pIdx int, amount int) *PieceTransform {
	return &PieceTransform{
		PieceIdx: byte(pIdx),
		RotateCW: 0,
		Activate: int8(amount),
	}
}

type PieceTransform struct {
	PieceIdx byte
	RotateCW int8 // 1=cw, -1=ccw, 0=nothing
	Activate int8 // 1=activate
}

var _ event.Event = (*PieceTransform)(nil)

func (e *PieceTransform) Do(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	if e.RotateCW > 0 {
		inverted := ctrl.Piece.RotateCW()
		animateRotatePiece(ctrl, true, inverted, f.Config.Anim)
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	} else if e.RotateCW < 0 {
		inverted := ctrl.Piece.RotateCCW()
		animateRotatePiece(ctrl, false, inverted, f.Config.Anim)
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	}

	if e.Activate != 0 {
		activations := ctrl.Piece.GetActivationCount() - int(e.Activate)
		ctrl.Piece.SetActivationCount(activations)
		if e.Activate > 0 && activations == 1 {
			animateBlinkPiece(ctrl, f.Config.Anim)
		}
	}
}

func (e *PieceTransform) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	ctrl.Piece = ctrl.Piece.Clone()

	if e.RotateCW > 0 {
		ctrl.Piece.RotateCCW()
		ctrl.List.Clear()
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	} else if e.RotateCW < 0 {
		ctrl.Piece.RotateCW()
		ctrl.List.Clear()
		ctrl.Blocks = piece.GetBlocks(ctrl.Piece, ctrl.Blocks[:0])
		updatePieceShadow(f, ctrl)
	}

	if e.Activate != 0 {
		activations := ctrl.Piece.GetActivationCount() + int(e.Activate)
		ctrl.Piece.SetActivationCount(activations)
		ctrl.List.Clear()
	}
}

func (e *PieceTransform) Equals(ev event.Event) bool {
	q, ok := ev.(*PieceTransform)
	return ok && e.PieceIdx == q.PieceIdx && e.RotateCW == q.RotateCW && e.Activate == q.Activate
}

func (e *PieceTransform) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.PieceIdx, byte(e.RotateCW), byte(e.Activate)}); err != nil {
		return err
	}
	return nil
}

func (e *PieceTransform) Read(r io.Reader) error {
	var buffer [3]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.PieceIdx = buffer[0]
	e.RotateCW = int8(buffer[1])
	e.Activate = int8(buffer[2])

	return nil
}

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
}

func (e *PieceFall) Undo(f *field.Field) {
	ctrl := f.Ctrl(e.PieceIdx)
	height := int(e.Height)
	ctrl.Y += height
	ctrl.List.Clear()
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
