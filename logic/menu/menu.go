// Copyright (c) 2024 by Marko Gaćeša

package menu

type Menu struct {
	title      string
	currentIdx int
	items      []Item
	stopFn     func()
}

func New(title string, stopFn func(), items ...Item) *Menu {
	m := &Menu{}

	m.SetTitle(title)
	m.SetItems(items...)
	m.stopFn = stopFn

	return m
}

func (m *Menu) SetTitle(title string) {
	m.title = title
}

func (m *Menu) SetItems(items ...Item) {
	if len(items) == 0 {
		panic("no menu items provided")
	}

	for _, item := range items {
		item.setParent(m)
	}

	m.items = items
}

func (m *Menu) Title() string {
	return m.title
}

func (m *Menu) Count() int        { return len(m.items) }
func (m *Menu) Item(idx int) Item { return m.items[idx] }
func (m *Menu) Current() Item     { return m.items[m.currentIdx] }
func (m *Menu) CurrentIdx() int   { return m.currentIdx }

func (m *Menu) Previous() {
	n := len(m.items)
	if n == 1 {
		return
	}

	m.items[m.currentIdx].FocusLost()
	m.currentIdx = (m.currentIdx - 1 + n) % n
	m.items[m.currentIdx].Focus()
}

func (m *Menu) Next() {
	n := len(m.items)
	if n == 1 {
		return
	}

	m.items[m.currentIdx].FocusLost()
	m.currentIdx = (m.currentIdx + 1) % n
	m.items[m.currentIdx].Focus()
}

func (m *Menu) Focus(idx int) {
	n := len(m.items)

	if n == 1 || idx < 0 || idx >= n || idx == m.currentIdx {
		return
	}

	m.items[m.currentIdx].FocusLost()
	m.currentIdx = idx
	m.items[m.currentIdx].Focus()
}

func (m *Menu) Decrease() { m.items[m.currentIdx].Decrease() }
func (m *Menu) Increase() { m.items[m.currentIdx].Increase() }

func (m *Menu) Input(r rune) {
	if consumed := m.items[m.currentIdx].Input(r); consumed {
		return
	}

	if r == InputEscape {
		m.stopFn()
	}
}
