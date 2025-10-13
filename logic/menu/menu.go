// Copyright (c) 2024, 2025 by Marko Gaćeša

package menu

import "sync"

type Menu struct {
	mx sync.Mutex

	title string

	itemsAll     []Item
	itemsVisible []Item
	currentIdx   int

	mutateCallback func(m *Menu)
}

// New creates new Menu. Menu requires some items
// The mutateCallback will be called whenever the state of any item has been changed.
func New(title string, mutateCallback func(m *Menu), items ...Item) *Menu {
	m := &Menu{
		title:          title,
		mutateCallback: mutateCallback,
	}

	m.SetItems(items...)

	return m
}

func (m *Menu) SetItems(items ...Item) {
	m.mx.Lock()

	m.itemsAll = items
	m.itemsVisible = make([]Item, 0, len(items))

	for _, item := range m.itemsAll {
		b := item.b()
		b.updateDisabled()
		b.updateVisible()
		if b.visible {
			m.itemsVisible = append(m.itemsVisible, item)
		}
	}

	if m.currentIdx >= len(m.itemsVisible) {
		m.currentIdx = 0
	}

	m.mx.Unlock()
}

func (m *Menu) Iteration(iter *Iter) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.updateItems()

	iter.Title = m.title
	iter.Description = ""
	iter.Current = 0
	iter.Items = iter.Items[:0]
	iter.IsDisabled = iter.IsDisabled[:0]

	count := len(m.itemsVisible)

	if count == 0 {
		return
	}

	currentIdx := m.currentIdx
	desc := m.itemsVisible[currentIdx].Description()

	iter.Description = desc
	iter.Current = currentIdx

	for _, item := range m.itemsVisible {
		iter.Items = append(iter.Items, item.Text())
		iter.IsDisabled = append(iter.IsDisabled, item.IsDisabled())
	}
}

func (m *Menu) Previous() {
	m.mx.Lock()
	defer m.mx.Unlock()

	n := len(m.itemsVisible)
	if n < 2 {
		return
	}

	m.itemsVisible[m.currentIdx].focusLost()
	m.currentIdx = (m.currentIdx - 1 + n) % n
	m.itemsVisible[m.currentIdx].focus()
}

func (m *Menu) Next() {
	m.mx.Lock()
	defer m.mx.Unlock()

	n := len(m.itemsVisible)
	if n < 2 {
		return
	}

	m.itemsVisible[m.currentIdx].focusLost()
	m.currentIdx = (m.currentIdx + 1) % n
	m.itemsVisible[m.currentIdx].focus()
}

func (m *Menu) Focus(idx int) {
	m.mx.Lock()
	defer m.mx.Unlock()

	n := len(m.itemsVisible)
	if idx < 0 || idx >= n || idx == m.currentIdx {
		return
	}

	m.itemsVisible[m.currentIdx].focusLost()
	m.currentIdx = idx
	m.itemsVisible[m.currentIdx].focus()
}

func (m *Menu) Decrease() {
	shouldCallback := func() bool {
		m.mx.Lock()
		defer m.mx.Unlock()

		if len(m.itemsVisible) == 0 {
			return false
		}

		item := m.itemsVisible[m.currentIdx]
		if item.IsDisabled() {
			return false
		}
		item.decrease()
		if item.isDirty() {
			m.updateItems()
			return true
		}

		return false
	}()
	if shouldCallback {
		m.callback()
	}
}

func (m *Menu) Increase() {
	shouldCallback := func() bool {
		m.mx.Lock()
		defer m.mx.Unlock()

		if len(m.itemsVisible) == 0 {
			return false
		}

		item := m.itemsVisible[m.currentIdx]
		if item.IsDisabled() {
			return false
		}
		item.increase()
		if item.isDirty() {
			m.updateItems()
			return true
		}

		return false
	}()
	if shouldCallback {
		m.callback()
	}
}

func (m *Menu) Input(r rune) {
	shouldCallback := func() bool {
		m.mx.Lock()
		defer m.mx.Unlock()

		if len(m.itemsVisible) > 0 {
			if item := m.itemsVisible[m.currentIdx]; !item.IsDisabled() && item.input(r) && item.isDirty() {
				m.updateItems()
				return true
			}
		}

		for _, item := range m.itemsAll {
			if item.b().global && item.input(r) {
				m.updateItems()
				return true
			}
		}

		return false
	}()
	if shouldCallback {
		m.callback()
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
}

func (m *Menu) callback() {
	if m.mutateCallback != nil {
		m.mutateCallback(m)
	}
}
