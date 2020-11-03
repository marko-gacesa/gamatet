// Copyright (c) 2020 by Marko Gaćeša

package event

import (
	"gamatet/game/field"
	"io"
)

type Event interface {
	Do(f *field.Field)
	Undo(f *field.Field)
	Equals(ev Event) bool
	Read(r io.Reader) error
	Write(w io.Writer) error
}

type Pusher interface {
	Push(Event)
}
