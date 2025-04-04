// Copyright (c) 2024,2025 by Marko Gaćeša

package menu

type Menu struct {
	title string

	itemsAll     []Item
	itemsVisible []Item
	currentIdx   int

	mutateCallback func()
}

// New creates new Menu. Menu requires some items
// The mutateCallback will be called whenever the state of any item has been changed.
func New(title string, mutateCallback func(), items ...Item) *Menu {
	if len(items) == 0 {
		panic("no menu items provided")
	}

	m := &Menu{
		title:          title,
		itemsAll:       items,
		itemsVisible:   make([]Item, 0, len(items)),
		currentIdx:     0,
		mutateCallback: mutateCallback,
	}

	for _, item := range m.itemsAll {
		b := item.b()
		b.updateDisabled()
		b.updateVisible()
		if b.visible {
			m.itemsVisible = append(m.itemsVisible, item)
		}
	}

	return m
}

func (m *Menu) Title() string {
	return m.title
}
func (m *Menu) Description() string {
	if len(m.itemsVisible) == 0 {
		return ""
	}

	return m.itemsVisible[m.currentIdx].Description()
}

func (m *Menu) Count() int        { return len(m.itemsVisible) }
func (m *Menu) Item(idx int) Item { return m.itemsVisible[idx] }
func (m *Menu) Current() Item     { return m.itemsVisible[m.currentIdx] }
func (m *Menu) CurrentIdx() int   { return m.currentIdx }

func (m *Menu) Previous() {
	n := len(m.itemsVisible)
	if n < 2 {
		return
	}

	m.itemsVisible[m.currentIdx].focusLost()
	m.currentIdx = (m.currentIdx - 1 + n) % n
	m.itemsVisible[m.currentIdx].focus()
}

func (m *Menu) Next() {
	n := len(m.itemsVisible)
	if n < 2 {
		return
	}

	m.itemsVisible[m.currentIdx].focusLost()
	m.currentIdx = (m.currentIdx + 1) % n
	m.itemsVisible[m.currentIdx].focus()
}

func (m *Menu) Decrease() {
	if len(m.itemsVisible) == 0 {
		m.cancel()
		return
	}

	item := m.itemsVisible[m.currentIdx]
	if item.IsDisabled() {
		return
	}
	item.decrease()
	if item.isDirty() {
		m.updateItems()
	}
}

func (m *Menu) Increase() {
	if len(m.itemsVisible) == 0 {
		m.cancel()
		return
	}

	item := m.itemsVisible[m.currentIdx]
	if item.IsDisabled() {
		return
	}
	item.increase()
	if item.isDirty() {
		m.updateItems()
	}
}

func (m *Menu) Input(r rune) {
	if len(m.itemsVisible) == 0 {
		m.cancel()
		return
	}

	item := m.itemsVisible[m.currentIdx]
	if !item.IsDisabled() && item.input(r) && item.isDirty() {
		m.updateItems()
		return
	}

	if r == InputEscape {
		m.cancel()
	}
}

func (m *Menu) cancel() {
	var canceled bool
	for _, it := range m.itemsAll {
		if it.b().canceler {
			it.input(InputEscape)
			canceled = true
		}
	}
	if canceled {
		m.updateItems()
	}
}

func (m *Menu) updateItems() {
	var curr Item
	if len(m.itemsVisible) > 0 {
		curr = m.itemsVisible[m.currentIdx]
	}

	m.itemsVisible = m.itemsVisible[:0]

	var seenCurr bool
	for i := range m.itemsAll {
		item := m.itemsAll[i]
		b := item.b()

		item.fix()
		b.markDirty()

		wasVisible := b.visible
		b.updateDisabled()
		b.updateVisible()
		isVisible := b.visible

		if isVisible {
			m.itemsVisible = append(m.itemsVisible, item)
		}

		if seenCurr = seenCurr || item == curr; seenCurr {
			continue
		}

		if wasVisible && !isVisible {
			m.currentIdx--
		} else if !wasVisible && isVisible {
			m.currentIdx++
		}
	}

	if m.currentIdx < 0 || m.currentIdx >= len(m.itemsVisible) {
		m.currentIdx = 0
	}

	if m.mutateCallback != nil {
		m.mutateCallback()
	}
}
