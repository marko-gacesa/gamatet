// Copyright (c) 2024 by Marko Gaćeša

package menu

type Item interface {
	Text() string
	Description() string

	Editable() bool
	Increase()
	Decrease()

	Input(r rune)

	Focus()
	FocusLost()

	setParent(menu *Menu)
}

type base struct {
	description string
}

func (b base) Description() string {
	return b.description
}
func (b base) Editable() bool { return false }

func (base) Input(rune) {}
func (base) Focus()     {}
func (base) FocusLost() {}

func (base) setParent(*Menu) {}

type editable struct {
	base
	label   string
	current string
}

func (e *editable) dirty() {
	e.current = ""
}

func (editable) Editable() bool { return true }
