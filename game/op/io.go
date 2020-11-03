// Copyright (c) 2020 by Marko Gaćeša

package op

import (
	"fmt"
	"gamatet/game/event"
	"gamatet/util"
	"io"
)

const (
	codeNoOp = iota

	// field events
	codeFieldStop
	codeFieldPause
	codeFieldUnpause
	codeFieldDestroyRow
	codeFieldDestroyColumn
	codeFieldBlockSet
	codeFieldBlockHardness
	codeFieldBlockTransform

	// piece events
	codePieceState
	codePieceSet
	codePieceMove
	codePieceTransform
	codePieceFall
)

type OpType byte

const (
	OpClear OpType = 0
	OpSet   OpType = 1
)

var (
	FieldStopBytes    = []byte{codeFieldStop}
	FieldPauseBytes   = []byte{codeFieldPause}
	FieldUnpauseBytes = []byte{codeFieldUnpause}
)

func Write(w io.Writer, e event.Event) (err error) {
	switch e.(type) {

	case *PieceState:
		err = util.Write8(w, codePieceState)
	case *PieceSet:
		err = util.Write8(w, codePieceSet)
	case *PieceMove:
		err = util.Write8(w, codePieceMove)
	case *PieceTransform:
		err = util.Write8(w, codePieceTransform)
	case *PieceFall:
		err = util.Write8(w, codePieceFall)

	case *FieldPause:
		err = util.Write8(w, codeFieldPause)
	case *FieldUnpause:
		err = util.Write8(w, codeFieldUnpause)
	case *FieldDestroyRow:
		err = util.Write8(w, codeFieldDestroyRow)
	case *FieldDestroyColumn:
		err = util.Write8(w, codeFieldDestroyColumn)
	case *FieldBlockSet:
		err = util.Write8(w, codeFieldBlockSet)
	case *FieldBlockHardness:
		err = util.Write8(w, codeFieldBlockHardness)
	case *FieldBlockTransform:
		err = util.Write8(w, codeFieldBlockTransform)

	default:
		err = fmt.Errorf("unrecognized event: %T", e)
	}
	if err != nil {
		return
	}

	err = e.Write(w)

	return
}

func Read(r io.Reader) (event.Event, error) {
	code, err := util.Read8(r)
	if err != nil {
		return nil, err
	}

	var e event.Event

	switch code {
	case codeFieldPause:
		e = &FieldPause{}
	case codeFieldUnpause:
		e = &FieldUnpause{}
	case codeFieldDestroyRow:
		e = &FieldDestroyRow{}
	case codeFieldDestroyColumn:
		e = &FieldDestroyColumn{}
	case codeFieldBlockSet:
		e = &FieldBlockSet{}
	case codeFieldBlockHardness:
		e = &FieldBlockHardness{}
	case codeFieldBlockTransform:
		e = &FieldBlockTransform{}

	case codePieceState:
		e = &PieceState{}
	case codePieceSet:
		e = &PieceSet{}
	case codePieceMove:
		e = &PieceMove{}
	case codePieceTransform:
		e = &PieceTransform{}
	case codePieceFall:
		e = &PieceFall{}

	default:
		return nil, fmt.Errorf("unrecognized event code: %d", code)
	}

	err = e.Read(r)
	if err != nil {
		return nil, err
	}

	return e, nil
}
