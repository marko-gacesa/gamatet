// Copyright (c) 2020-2024 by Marko Gaćeša

package op

import (
	"bytes"
	"gamatet/game/block"
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/piece"
	"io"
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

type FieldPause struct{}

var _ event.Event = FieldPause{}

func (e FieldPause) Do(f *field.Field)   { f.Pause() }
func (e FieldPause) Undo(f *field.Field) { f.Unpause() }

func (e FieldPause) Equals(ev event.Event) bool {
	_, ok := ev.(FieldPause)
	return ok
}

func (e FieldPause) Read(io.Reader) error  { return nil }
func (e FieldPause) Write(io.Writer) error { return nil }

type FieldUnpause struct{}

var _ event.Event = FieldUnpause{}

func (e FieldUnpause) Do(f *field.Field)   { f.Unpause() }
func (e FieldUnpause) Undo(f *field.Field) { f.Pause() }

func (e FieldUnpause) Equals(ev event.Event) bool {
	_, ok := ev.(*FieldUnpause)
	return ok
}

func (e FieldUnpause) Read(io.Reader) error  { return nil }
func (e FieldUnpause) Write(io.Writer) error { return nil }

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
	return ok && e.Row == q.Row && block.SliceEqual(e.Blocks, q.Blocks)
}

func (e *FieldDestroyRow) Write(w io.Writer) error {
	if _, err := w.Write([]byte{e.Row, byte(len(e.Blocks))}); err != nil {
		return err
	}
	for i := 0; i < len(e.Blocks); i++ {
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

func NewFieldBlockSet(col, row int, op OpType, animType, animParam int, b block.Block) *FieldBlockSet {
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
	Op        OpType // 0=clear (the Block contains the block to be cleared), 1=set (the Block contains the block to be added)
	AnimType  byte
	AnimParam byte
	Block     block.Block
}

var _ event.Event = (*FieldBlockSet)(nil)

func (e *FieldBlockSet) Do(f *field.Field) {
	switch e.Op {
	case OpSet:
		f.SetXY(int(e.Col), int(e.Row), int(e.AnimType), int(e.AnimParam), e.Block)
	case OpClear:
		_ = f.ClearXY(int(e.Col), int(e.Row), int(e.AnimType), int(e.AnimParam))
	}
	updateAllPiecesShadow(f)
}

func (e *FieldBlockSet) Undo(f *field.Field) {
	switch e.Op {
	case OpSet:
		_ = f.ClearXY(int(e.Col), int(e.Row), field.AnimNo, 0)
	case OpClear:
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
	e.Op = OpType(buffer[2])
	e.AnimType = buffer[3]
	e.AnimParam = buffer[4]

	if err := e.Block.Read(r); err != nil {
		return err
	}

	return nil
}

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

func NewFieldLost(f *field.Field) *FieldLost {
	n := byte(f.Ctrls())
	ctrlStates := make([]byte, 0, 2*n)
	for i := byte(0); i < n; i += 2 {
		state := f.Ctrl(i).State
		if state == piece.StateLost {
			continue
		}
		ctrlStates = append(ctrlStates, i)
		ctrlStates = append(ctrlStates, byte(state))
	}

	return &FieldLost{
		CtrlStates: ctrlStates,
	}
}

type FieldLost struct {
	CtrlStates []byte
}

var _ event.Event = (*FieldLost)(nil)

func (e *FieldLost) Do(f *field.Field) {
	for i := 0; i < len(e.CtrlStates); i += 2 {
		ctrl := f.Ctrl(e.CtrlStates[i])
		ctrl.State = piece.StateLost
		ctrl.RestartTimer(0)
	}
}

func (e *FieldLost) Undo(f *field.Field) {
	for i := 0; i < len(e.CtrlStates); i += 2 {
		ctrl := f.Ctrl(e.CtrlStates[i])
		ctrl.State = piece.State(e.CtrlStates[i+1])
		ctrl.RestartTimer(0)
	}
}

func (e *FieldLost) Equals(ev event.Event) bool {
	q, ok := ev.(*FieldLost)
	return ok && bytes.Equal(e.CtrlStates, q.CtrlStates)
}

func (e *FieldLost) Write(w io.Writer) error {
	if _, err := w.Write([]byte{byte(len(e.CtrlStates))}); err != nil {
		return err
	}
	if _, err := w.Write(e.CtrlStates); err != nil {
		return err
	}
	return nil
}

func (e *FieldLost) Read(r io.Reader) error {
	var buffer [1]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	e.CtrlStates = make([]byte, buffer[0])
	if _, err := io.ReadFull(r, e.CtrlStates); err != nil {
		return err
	}

	return nil
}
