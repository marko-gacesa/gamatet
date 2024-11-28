// Copyright (c) 2024 by Marko Gaćeša

package menu

const (
	InputEnter     = '\n'
	InputBackspace = '\b'
	InputDelete    = 0xFF
	InputEscape    = 0x1B
)

type Item interface {
	Text() string
	Description() string

	Editable() bool
	Increase()
	Decrease()

	Input(r rune) bool

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

func (base) Input(rune) bool { return false }
func (base) Focus()          {}
func (base) FocusLost()      {}

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
