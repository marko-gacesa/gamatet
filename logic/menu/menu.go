// Copyright (c) 2024 by Marko Gaćeša

package menu

type Menu struct {
	current int
	items   []Item
}

func NewMenu(items ...Item) *Menu {
	if len(items) == 0 {
		panic("no menu items provided")
	}

	return &Menu{items: items}
}

func (m *Menu) Current() int      { return m.current }
func (m *Menu) Count() int        { return len(m.items) }
func (m *Menu) Item(idx int) Item { return m.items[idx] }

func (m *Menu) Previous() {
	n := len(m.items)
	if n == 1 {
		return
	}

	m.items[m.current].FocusLost()
	m.current = (m.current - 1 + n) % n
	m.items[m.current].Focus()
}

func (m *Menu) Next() {
	n := len(m.items)
	if n == 1 {
		return
	}

	m.items[m.current].FocusLost()
	m.current = (m.current + 1) % n
	m.items[m.current].Focus()
}

func (m *Menu) Decrease() { m.items[m.current].Decrease() }
func (m *Menu) Increase() { m.items[m.current].Increase() }

func (m *Menu) Input(r rune) { m.items[m.current].Input(r) }
