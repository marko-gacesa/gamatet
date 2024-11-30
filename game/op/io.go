// Copyright (c) 2020-2024 by Marko Gaćeša

package op

import (
	"fmt"
	"gamatet/game/event"
	"gamatet/game/serialize"
	"io"
)

const (
	// field events
	codeFieldStop event.Code = iota
	codeFieldPause
	codeFieldUnpause
	codeFieldDestroyRow
	codeFieldDestroyColumn
	codeFieldBlockSet
	codeFieldBlockHardness
	codeFieldBlockTransform
	codeFieldExBlock
	codeFieldLost
	codeFieldQuake

	// piece events
	codePieceState
	codePieceSet
	codePieceMove
	codePieceTransform
	codePieceFall
)

type Type byte

const (
	TypeClear Type = 0
	TypeSet   Type = 1
)

var (
	FieldStopBytes    = []byte{byte(codeFieldStop)}
	FieldPauseBytes   = []byte{byte(codeFieldPause)}
	FieldUnpauseBytes = []byte{byte(codeFieldUnpause)}
)

func Write(w io.Writer, e event.Event) error {
	err := serialize.Write8(w, byte(e.TypeID()))
	if err != nil {
		return err
	}

	return e.Write(w)
}

func Read(r io.Reader) (event.Event, error) {
	code, err := serialize.Read8(r)
	if err != nil {
		return nil, err
	}

	e := instance(event.Code(code))

	err = e.Read(r)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func instance(code event.Code) event.Event {
	var e event.Event

	switch code {
	case codeFieldStop:
		e = FieldStop{}

	case codeFieldPause:
		e = FieldPause{}
	case codeFieldUnpause:
		e = FieldUnpause{}

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
	case codeFieldExBlock:
		e = &FieldExBlock{}
	case codeFieldLost:
		e = &FieldLost{}
	case codeFieldQuake:
		e = &FieldQuake{}

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
		panic(fmt.Sprintf("unrecognized event code=%d", code))
	}

	return e
}
