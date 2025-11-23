// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package op

import (
	"bytes"
	"encoding/binary"
	"io"
	"slices"

	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
)

// FieldStop is a signal that not more events will be fired
type FieldStop struct{}

var _ event.Event = FieldStop{}

func (e FieldStop) Do(f *field.Field) { f.CloseDone() }
func (e FieldStop) Undo(*field.Field) { /* can't undo */ }

func (e FieldStop) Equals(ev event.Event) bool {
	_, ok := ev.(FieldStop)
	return ok
}

func (e FieldStop) Read(io.Reader) error  { return nil }
func (e FieldStop) Write(io.Writer) error { return nil }

func (e FieldStop) TypeID() event.Code { return codeFieldStop }

func NewFieldMode(f *field.Field, modeNew field.Mode, stopPieces bool) *FieldMode {
	n := byte(f.Ctrls())
	ctrlStates := make([]byte, 2*n)
	for i := range n {
		state := f.Ctrl(i).State
		ctrlStates[i*2] = i
		ctrlStates[i*2+1] = byte(state)
	}

	return &FieldMode{
		ModeNew:    modeNew,
		ModeOld:    f.GetMode(),
		CtrlStates: ctrlStates,
		StopPieces: stopPieces,
	}
}

type FieldMode struct {
	ModeNew    field.Mode
	ModeOld    field.Mode
	CtrlStates []byte
	StopPieces bool
}

var _ event.Event = (*FieldMode)(nil)

func (e *FieldMode) Do(f *field.Field) {
	f.SetMode(e.ModeNew)
	if !e.StopPieces {
		return
	}
	for i := 0; i < len(e.CtrlStates); i += 2 {
		ctrlIdx := e.CtrlStates[i]
		ctrl := f.Ctrl(ctrlIdx)
		ctrl.State = piece.StateStop
		ctrl.RestartTimer(0)
	}
}

func (e *FieldMode) Undo(f *field.Field) {
	f.SetMode(e.ModeOld)
	if !e.StopPieces {
		return
	}
	for i := 0; i < len(e.CtrlStates); i += 2 {
		ctrlIdx := e.CtrlStates[i]
		ctrlState := e.CtrlStates[i+1]
		ctrl := f.Ctrl(ctrlIdx)
		ctrl.State = piece.State(ctrlState)
		ctrl.RestartTimer(0)
	}
}

func (e *FieldMode) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldMode)
	return ok && q.ModeNew == e.ModeNew && q.ModeOld == e.ModeOld &&
		bytes.Equal(q.CtrlStates, e.CtrlStates) && q.StopPieces == e.StopPieces
}

func (e *FieldMode) Read(r io.Reader) error {
	var bufferMode [2]byte
	if _, err := io.ReadFull(r, bufferMode[:]); err != nil {
		return err
	}

	e.ModeNew = field.Mode(bufferMode[0])
	e.ModeOld = field.Mode(bufferMode[1])

	var bufferStates [1]byte
	if _, err := io.ReadFull(r, bufferStates[:]); err != nil {
		return err
	}

	e.CtrlStates = make([]byte, bufferStates[0])
	if _, err := io.ReadFull(r, e.CtrlStates); err != nil {
		return err
	}

	var bufferStateNew [1]byte
	if _, err := io.ReadFull(r, bufferStateNew[:]); err != nil {
		return err
	}

	e.StopPieces = bufferStateNew[0] != 0

	return nil
}

func (e *FieldMode) Write(w io.Writer) error {
	if _, err := w.Write([]byte{byte(e.ModeNew), byte(e.ModeOld)}); err != nil {
		return err
	}

	if _, err := w.Write([]byte{byte(len(e.CtrlStates))}); err != nil {
		return err
	}
	if _, err := w.Write(e.CtrlStates); err != nil {
		return err
	}

	var stopPieces byte
	if e.StopPieces {
		stopPieces = 1
	}

	if _, err := w.Write([]byte{stopPieces}); err != nil {
		return err
	}

	return nil
}

func (e *FieldMode) TypeID() event.Code { return codeFieldMode }

func NewFieldDestroyRow(row int, blocks []block.Block) *FieldDestroyRow {
	return &FieldDestroyRow{
		Row:    byte(row),
		Blocks: blocks,
	}
}

type FieldDestroyRow struct {
	Row    byte
	Blocks []block.Block
}

var _ event.Event = (*FieldDestroyRow)(nil)

func (e *FieldDestroyRow) Do(f *field.Field) {
	f.ShiftRowsDown(int(e.Row))
	updateAllPiecesShadow(f)
}

func (e *FieldDestroyRow) Undo(f *field.Field) {
	f.UndoShiftRowsDown(int(e.Row), e.Blocks)
	updateAllPiecesShadow(f)
}

func (e *FieldDestroyRow) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldDestroyRow)
	return ok && e.Row == q.Row && slices.Equal(e.Blocks, q.Blocks)
}

func (e *FieldDestroyRow) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.Row, byte(len(e.Blocks))}); err != nil {
		return err
	}
	for i := range e.Blocks {
		if err := e.Blocks[i].Write(w); err != nil {
			return err
		}
	}
	return nil
}

func (e *FieldDestroyRow) Read(r io.Reader) error {
	var buffer [2]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.Row = buffer[0]
	e.Blocks = make([]block.Block, buffer[1])

	for i := 0; i < len(e.Blocks); i++ {
		if err := e.Blocks[i].Read(r); err != nil {
			return err
		}
	}

	return nil
}

func (e *FieldDestroyRow) TypeID() event.Code { return codeFieldDestroyRow }

func NewFieldDestroyColumn(col, row, n, height int, b block.Block) *FieldDestroyColumn {
	return &FieldDestroyColumn{
		Col:    byte(col),
		Row:    byte(row),
		N:      byte(n),
		Height: byte(height),
		Block:  b,
	}
}

type FieldDestroyColumn struct {
	Col, Row  byte
	N, Height byte
	Block     block.Block
}

var _ event.Event = (*FieldDestroyColumn)(nil)

func (e *FieldDestroyColumn) Do(f *field.Field) {
	f.ShiftColumnDownByN(int(e.Col), int(e.Row), int(e.N), int(e.Height))
	updateAllPiecesShadow(f)
}

func (e *FieldDestroyColumn) Undo(f *field.Field) {
	f.UndoShiftColumnByN(int(e.Col), int(e.Row), int(e.N), int(e.Height), e.Block)
	updateAllPiecesShadow(f)
}

func (e *FieldDestroyColumn) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldDestroyColumn)
	return ok && e.Col == q.Col && e.Row == q.Row &&
		e.N == q.N && e.Height == q.Height && e.Block == q.Block
}

func (e *FieldDestroyColumn) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.Col, e.Row, e.N, e.Height}); err != nil {
		return err
	}
	if err := e.Block.Write(w); err != nil {
		return err
	}
	return nil
}

func (e *FieldDestroyColumn) Read(r io.Reader) error {
	var buffer [4]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.Col = buffer[0]
	e.Row = buffer[1]
	e.N = buffer[2]
	e.Height = buffer[3]

	if err := e.Block.Read(r); err != nil {
		return err
	}

	return nil
}

func (e *FieldDestroyColumn) TypeID() event.Code { return codeFieldDestroyColumn }

func NewFieldBlockSet(col, row int, op Type, animType, animParam int, b block.Block) *FieldBlockSet {
	return &FieldBlockSet{
		Col:       byte(col),
		Row:       byte(row),
		Op:        op,
		AnimType:  byte(animType),
		AnimParam: byte(animParam),
		Block:     b,
	}
}

type FieldBlockSet struct {
	Col, Row  byte
	Op        Type // 0=clear (the Block contains the block to be cleared), 1=set (the Block contains the block to be added)
	AnimType  byte
	AnimParam byte
	Block     block.Block
}

var _ event.Event = (*FieldBlockSet)(nil)

func (e *FieldBlockSet) Do(f *field.Field) {
	switch e.Op {
	case TypeSet:
		f.SetXY(int(e.Col), int(e.Row), int(e.AnimType), int(e.AnimParam), e.Block)
	case TypeClear:
		_ = f.ClearXY(int(e.Col), int(e.Row), int(e.AnimType), int(e.AnimParam))
	}
	updateAllPiecesShadow(f)
}

func (e *FieldBlockSet) Undo(f *field.Field) {
	switch e.Op {
	case TypeSet:
		_ = f.ClearXY(int(e.Col), int(e.Row), field.AnimNo, 0)
	case TypeClear:
		f.SetXY(int(e.Col), int(e.Row), field.AnimNo, 0, e.Block)
	}
	updateAllPiecesShadow(f)
}

func (e *FieldBlockSet) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldBlockSet)
	return ok && e.Col == q.Col && e.Row == q.Row && e.Op == q.Op && e.Block == q.Block
}

func (e *FieldBlockSet) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.Col, e.Row, byte(e.Op), e.AnimType, e.AnimParam}); err != nil {
		return err
	}
	if err := e.Block.Write(w); err != nil {
		return err
	}
	return nil
}

func (e *FieldBlockSet) Read(r io.Reader) error {
	var buffer [5]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.Col = buffer[0]
	e.Row = buffer[1]
	e.Op = Type(buffer[2])
	e.AnimType = buffer[3]
	e.AnimParam = buffer[4]

	if err := e.Block.Read(r); err != nil {
		return err
	}

	return nil
}

func (e *FieldBlockSet) TypeID() event.Code { return codeFieldBlockSet }

func NewFieldBlockHardness(col, row, hardness, animType, animParam int) *FieldBlockHardness {
	return &FieldBlockHardness{
		Col:       byte(col),
		Row:       byte(row),
		Hardness:  int8(hardness),
		AnimType:  byte(animType),
		AnimParam: byte(animParam),
	}
}

type FieldBlockHardness struct {
	Col, Row  byte
	Hardness  int8
	AnimType  byte
	AnimParam byte
}

var _ event.Event = (*FieldBlockHardness)(nil)

func (e *FieldBlockHardness) Do(f *field.Field) {
	f.HardnessXY(int(e.Col), int(e.Row), int(e.Hardness), int(e.AnimType), int(e.AnimParam))
}

func (e *FieldBlockHardness) Undo(f *field.Field) {
	f.HardnessXY(int(e.Col), int(e.Row), -int(e.Hardness), field.AnimNo, 0)
}

func (e *FieldBlockHardness) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldBlockHardness)
	return ok && e.Col == q.Col && e.Row == q.Row && e.Hardness == q.Hardness
}

func (e *FieldBlockHardness) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.Col, e.Row, byte(e.Hardness), e.AnimType, e.AnimParam}); err != nil {
		return err
	}
	return nil
}

func (e *FieldBlockHardness) Read(r io.Reader) error {
	var buffer [5]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.Col = buffer[0]
	e.Row = buffer[1]
	e.Hardness = int8(buffer[2])
	e.AnimType = buffer[3]
	e.AnimParam = buffer[4]

	return nil
}

func (e *FieldBlockHardness) TypeID() event.Code { return codeFieldBlockHardness }

func NewFieldBlockTransform(col, row int, oldBlock, newBlock block.Block, animType, animParam int) *FieldBlockTransform {
	return &FieldBlockTransform{
		Col:       byte(col),
		Row:       byte(row),
		OldBlock:  oldBlock,
		NewBlock:  newBlock,
		AnimType:  byte(animType),
		AnimParam: byte(animParam),
	}
}

type FieldBlockTransform struct {
	Col, Row  byte
	OldBlock  block.Block
	NewBlock  block.Block
	AnimType  byte
	AnimParam byte
}

var _ event.Event = (*FieldBlockTransform)(nil)

func (e *FieldBlockTransform) Do(f *field.Field) {
	f.TransformXY(int(e.Col), int(e.Row), int(e.AnimType), int(e.AnimParam), e.OldBlock, e.NewBlock)
}

func (e *FieldBlockTransform) Undo(f *field.Field) {
	f.TransformXY(int(e.Col), int(e.Row), field.AnimNo, 0, e.NewBlock, e.OldBlock)
}

func (e *FieldBlockTransform) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldBlockTransform)
	return ok && e.Col == q.Col && e.Row == q.Row && e.OldBlock == q.OldBlock && e.NewBlock == q.NewBlock
}

func (e *FieldBlockTransform) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.Col, e.Row, e.AnimType, e.AnimParam}); err != nil {
		return err
	}
	if err := e.OldBlock.Write(w); err != nil {
		return err
	}
	if err := e.NewBlock.Write(w); err != nil {
		return err
	}
	return nil
}

func (e *FieldBlockTransform) Read(r io.Reader) error {
	var buffer [4]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.Col = buffer[0]
	e.Row = buffer[1]
	e.AnimType = buffer[2]
	e.AnimParam = buffer[3]

	if err := e.OldBlock.Read(r); err != nil {
		return err
	}
	if err := e.NewBlock.Read(r); err != nil {
		return err
	}

	return nil
}

func (e *FieldBlockTransform) TypeID() event.Code { return codeFieldBlockTransform }

func NewFieldExBlock(col, row int, animType, animParam int, b block.Block) *FieldExBlock {
	return &FieldExBlock{
		Col:       byte(col),
		Row:       byte(row),
		AnimType:  byte(animType),
		AnimParam: byte(animParam),
		Block:     b,
	}
}

type FieldExBlock struct {
	Col, Row  byte
	AnimType  byte
	AnimParam byte
	Block     block.Block
}

var _ event.Event = (*FieldExBlock)(nil)

func (e *FieldExBlock) Do(f *field.Field) {
	f.AddExXY(int(e.Col), int(e.Row), int(e.AnimType), int(e.AnimParam), e.Block)
}

func (e *FieldExBlock) Undo(f *field.Field) {}

func (e *FieldExBlock) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldExBlock)
	return ok && e.Col == q.Col && e.Row == q.Row && e.Block == q.Block
}

func (e *FieldExBlock) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.Col, e.Row, e.AnimType, e.AnimParam}); err != nil {
		return err
	}
	if err := e.Block.Write(w); err != nil {
		return err
	}
	return nil
}

func (e *FieldExBlock) Read(r io.Reader) error {
	var buffer [4]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.Col = buffer[0]
	e.Row = buffer[1]
	e.AnimType = buffer[2]
	e.AnimParam = buffer[3]

	if err := e.Block.Read(r); err != nil {
		return err
	}

	return nil
}

func (e *FieldExBlock) TypeID() event.Code { return codeFieldExBlock }

func NewFieldStat(removed, softened int16) *FieldStat {
	return &FieldStat{
		BlocksRemoved:  removed,
		BlocksSoftened: softened,
	}
}

type FieldStat struct {
	BlocksRemoved  int16
	BlocksSoftened int16
}

var _ event.Event = (*FieldStat)(nil)

func (e *FieldStat) Do(f *field.Field) {
	f.UpdateBlocksRemoved(int(e.BlocksRemoved))
}

func (e *FieldStat) Undo(f *field.Field) {
	f.UpdateBlocksRemoved(int(-e.BlocksRemoved))
}

func (e *FieldStat) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldStat)
	return ok && e.BlocksRemoved == q.BlocksRemoved && e.BlocksSoftened == q.BlocksSoftened
}

func (e *FieldStat) Read(r io.Reader) error {
	var buffer [4]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.BlocksRemoved = int16(binary.LittleEndian.Uint16(buffer[:2]))
	e.BlocksSoftened = int16(binary.LittleEndian.Uint16(buffer[2:4]))

	return nil
}

func (e *FieldStat) Write(w io.Writer) error {
	var buffer [4]byte

	binary.LittleEndian.PutUint16(buffer[:2], uint16(e.BlocksRemoved))
	binary.LittleEndian.PutUint16(buffer[2:4], uint16(e.BlocksSoftened))

	if _, err := w.Write(buffer[:]); err != nil {
		return err
	}

	return nil
}

func (e *FieldStat) TypeID() event.Code {
	return codeFieldStat
}

func NewFieldQuake(intensity byte) *FieldQuake {
	return &FieldQuake{
		Intensity: intensity,
	}
}

type FieldQuake struct {
	Intensity byte
}

var _ event.Event = (*FieldQuake)(nil)

func (e *FieldQuake) Do(f *field.Field) {
	f.AnimQuake(e.Intensity)
}

func (e *FieldQuake) Undo(*field.Field) {}

func (e *FieldQuake) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldQuake)
	return ok && e.Intensity == q.Intensity
}

func (e *FieldQuake) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.Intensity}); err != nil {
		return err
	}
	return nil
}

func (e *FieldQuake) Read(r io.Reader) error {
	var buffer [1]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.Intensity = buffer[0]
	return nil
}

func (e *FieldQuake) TypeID() event.Code { return codeFieldQuake }
