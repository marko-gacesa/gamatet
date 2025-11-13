// Copyright (c) 2024, 2025 by Marko Gaćeša

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

	IsDisabled() bool
	isVisible() bool

	fix()
	increase()
	decrease()
	input(r rune) bool

	focus()
	focusLost()

	updateDisabled()
	updateVisible()

	markDirty()
	isDirty() bool

	b() *base
}

const strNilPointer = "need non-nil pointer"

func makeBase(label, description string) base {
	return base{
		label:       label,
		description: description,
		visible:     true,
	}
}

type base struct {
	label       string
	description string
	current     string

	global   bool
	disabled bool
	visible  bool

	disabledFn    func() bool
	visibleFn     func() bool
	labelFn       func() string
	descriptionFn func() string
}

func (b *base) Text() string        { return b.getLabel() }
func (b *base) Description() string { return b.getDescription() }

func (b *base) getLabel() string {
	if b.labelFn != nil {
		return b.labelFn()
	}
	return b.label
}

func (b *base) getDescription() string {
	if b.descriptionFn != nil {
		return b.descriptionFn()
	}
	return b.description
}

func (b *base) IsDisabled() bool {
	return b.disabled
}
func (b *base) isVisible() bool { return b.visible }

func (*base) focus()     {}
func (*base) focusLost() {}

func (*base) fix()      {}
func (*base) increase() {}
func (*base) decrease() {}

func (b *base) updateDisabled() {
	if b.disabledFn != nil {
		b.disabled = b.disabledFn()
		return
	}
	b.disabled = false
}

func (b *base) updateVisible() {
	if b.visibleFn != nil {
		b.visible = b.visibleFn()
		return
	}
	b.visible = true
}

func (b *base) markDirty() { b.current = "" }
func (b *base) isDirty() bool {
	return b.current == ""
}

func (b *base) b() *base { return b }

func (b *base) String() string {
	return b.label
}
