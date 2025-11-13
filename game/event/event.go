// Copyright (c) 2020 by Marko Gaćeša

package event

import (
	"io"

	"github.com/marko-gacesa/gamatet/game/field"
)

type Code byte

type Event interface {
	Do(f *field.Field)
	Undo(f *field.Field)
	Equals(ev Event) bool
	Read(r io.Reader) error
	Write(w io.Writer) error
	TypeID() Code
}

type Pusher interface {
	Push(Event)
}
