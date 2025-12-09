// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package op

import (
	"fmt"
	"io"

	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/serialize"
)

const (
	// field events
	codeFieldStop event.Code = iota
	codeFieldMode
	codeFieldDestroyRow
	codeFieldDestroyColumn
	codeFieldBlockSet
	codeFieldBlockHardness
	codeFieldBlockTransform
	codeFieldColumnShift
	codeFieldExBlock
	codeFieldStat
	codeFieldEffect
	codeFieldQuake

	// piece events
	codePieceState
	codePieceSet
	codePieceMove
	codePieceRotate
	codePieceFlip
	codePieceShoot
	codePieceFall
	codePieceLevelBoost
	codePieceOverride
	codePieceSpeedUp
)

type Type byte

const (
	TypeClear Type = 0
	TypeSet   Type = 1
)

var FieldStopBytes = []byte{byte(codeFieldStop)}

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

	case codeFieldMode:
		e = &FieldMode{}
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
	case codeFieldColumnShift:
		e = &FieldColumnShift{}
	case codeFieldExBlock:
		e = &FieldExBlock{}
	case codeFieldStat:
		e = &FieldStat{}
	case codeFieldEffect:
		e = &FieldEffect{}
	case codeFieldQuake:
		e = &FieldQuake{}

	case codePieceState:
		e = &PieceState{}
	case codePieceSet:
		e = &PieceSet{}
	case codePieceMove:
		e = &PieceMove{}
	case codePieceRotate:
		e = &PieceRotate{}
	case codePieceFlip:
		e = &PieceFlip{}
	case codePieceShoot:
		e = &PieceShoot{}
	case codePieceFall:
		e = &PieceFall{}
	case codePieceLevelBoost:
		e = &PieceLevelBoost{}
	case codePieceOverride:
		e = &PieceOverride{}
	case codePieceSpeedUp:
		e = &PieceSpeedUp{}

	default:
		panic(fmt.Sprintf("unrecognized event code=%d", code))
	}

	return e
}
