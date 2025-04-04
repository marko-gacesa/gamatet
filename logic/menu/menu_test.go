// Copyright (c) 2025 by Marko Gaćeša

package menu

import (
	"testing"
)

func TestMenu(t *testing.T) {
	type dataStruct struct {
		show bool
		cmd  string
	}
	var data dataStruct
	tests := []struct {
		name   string
		items  []Item
		data   dataStruct
		mutate func(m *Menu)

		expCurrentIdx int
		expCount      int
	}{
		{
			name: "hide@first",
			items: []Item{
				NewBool(&data.show, "X", ""),
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return data.show })),
			},
			data: dataStruct{show: true},
			mutate: func(m *Menu) {
				data.show = false
				m.updateItems()
			},
			expCurrentIdx: 0,
			expCount:      1,
		},
		{
			name: "hide@last",
			items: []Item{
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return data.show })),
				NewBool(&data.show, "X", ""),
			},
			data: dataStruct{show: true},
			mutate: func(m *Menu) {
				m.Next()
				m.Next()
				m.Next()
				data.show = false
				m.updateItems()
			},
			expCurrentIdx: 0,
			expCount:      1,
		},
		{
			name: "hide@mid",
			items: []Item{
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewBool(&data.show, "X", ""),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "d", "D", "", WithVisible(func() bool { return data.show })),
			},
			data: dataStruct{show: true},
			mutate: func(m *Menu) {
				m.Next()
				m.Next()
				data.show = false
				m.updateItems()
			},
			expCurrentIdx: 0,
			expCount:      1,
		},
		{
			name: "show@first",
			items: []Item{
				NewBool(&data.show, "X", ""),
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return data.show })),
			},
			data: dataStruct{},
			mutate: func(m *Menu) {
				data.show = true
				m.updateItems()
			},
			expCurrentIdx: 0,
			expCount:      4,
		},
		{
			name: "show@last",
			items: []Item{
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return data.show })),
				NewBool(&data.show, "X", ""),
			},
			data: dataStruct{},
			mutate: func(m *Menu) {
				data.show = true
				m.updateItems()
			},
			expCurrentIdx: 3,
			expCount:      4,
		},
		{
			name: "show@mid",
			items: []Item{
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewBool(&data.show, "X", ""),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "d", "D", "", WithVisible(func() bool { return data.show })),
			},
			data: dataStruct{},
			mutate: func(m *Menu) {
				data.show = true
				m.updateItems()
			},
			expCurrentIdx: 2,
			expCount:      5,
		},
		{
			name: "show+hide@mid",
			items: []Item{
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return !data.show })),
				NewBool(&data.show, "X", ""),
				NewCommand(&data.cmd, "d", "D", "", WithVisible(func() bool { return !data.show })),
				NewCommand(&data.cmd, "e", "E", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "f", "F", "", WithVisible(func() bool { return data.show })),
			},
			data: dataStruct{show: true},
			mutate: func(m *Menu) {
				m.Next()
				m.Next()
				data.show = false
				m.updateItems()
			},
			expCurrentIdx: 1,
			expCount:      3,
		},
		{
			name: "hide-all",
			items: []Item{
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "d", "D", "", WithVisible(func() bool { return data.show })),
			},
			data: dataStruct{show: true},
			mutate: func(m *Menu) {
				m.Next()
				m.Next()
				data.show = false
				m.updateItems()
			},
			expCurrentIdx: 0,
			expCount:      0,
		},
		{
			name: "show-all",
			items: []Item{
				NewCommand(&data.cmd, "a", "A", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "b", "B", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "c", "C", "", WithVisible(func() bool { return data.show })),
				NewCommand(&data.cmd, "d", "D", "", WithVisible(func() bool { return data.show })),
			},
			data: dataStruct{show: false},
			mutate: func(m *Menu) {
				data.show = true
				m.updateItems()
			},
			expCurrentIdx: 0,
			expCount:      4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data = test.data
			m := New("", nil, test.items...)

			test.mutate(m)

			if want, got := test.expCurrentIdx, m.CurrentIdx(); want != got {
				t.Errorf("CurrentIdx mismatch: want=%d, got=%d", want, got)
			}
			if want, got := test.expCount, m.Count(); want != got {
				t.Errorf("Count mismatch: want=%d, got=%d", want, got)
			}
		})
	}
}
