// Copyright (c) 2024 by Marko Gaćeša

package menu

type Item interface {
	Text() string
	Description() string
	Increase()
	Decrease()
	Input(r rune)

	Focus()
	FocusLost()
}

type base struct {
	description string
}

func (b base) Description() string {
	return b.description
}

func (base) Input(rune) {}
func (base) Focus()     {}
func (base) FocusLost() {}

type editable struct {
	base
	label   string
	current string
}

func (e *editable) dirty() {
	e.current = ""
}
