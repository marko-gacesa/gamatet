// Copyright (c) 2024 by Marko Gaćeša

package menu

import (
	"context"
	"gamatet/graphics/render"
	"gamatet/graphics/screen"
	"gamatet/graphics/texture"
	"gamatet/logic/menu"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

const (
	contentW = 20
	contentH = 14
)

type Menu struct {
	tex  *texture.Manager
	text render.Text

	menu *menu.Menu

	strCache  []string
	strCache2 []string

	modelItem []mgl32.Mat4
	modelDesc mgl32.Mat4
}

var _ screen.Screen = (*Menu)(nil)

var (
	colorSelected    = mgl32.Vec4{0.8, 0.7, 0.6, 1}
	colorItem        = colorSelected.Mul(0.8)
	colorDescription = colorSelected.Mul(0.6)
)

func NewMenu(ctx context.Context, tex *texture.Manager, items ...menu.Item) *Menu {
	text := render.MakeText(tex, render.Font)

	n := len(items)

	m := menu.NewMenu(items...)
	strCache := make([]string, 2*n)

	modelItem := make([]mgl32.Mat4, n)
	for i := 0; i < n; i++ {
		modelItem[i] = mgl32.Ident4().
			Mul4(mgl32.Translate3D(-contentW/2+1, float32(n-contentH/2+1-i), 0))
	}

	modelDesc := mgl32.Ident4().
		Mul4(mgl32.Translate3D(-contentW/2+1, float32(-contentH/2+0.5), 0)).
		Mul4(mgl32.Scale3D(0.6, 0.6, 0.6))

	return &Menu{
		tex:       tex,
		text:      *text,
		menu:      m,
		strCache:  strCache,
		modelItem: modelItem,
		modelDesc: modelDesc,
	}
}

func (m *Menu) Release() {
	m.text.Release()
}

func (m *Menu) Shutdown() {
}

func (m *Menu) SetCamera(r *render.Renderer, w, h int) {
	r.OrthogonalFull(w, h, contentW, contentH, 2)
}

func (m *Menu) InputKey(key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey) {
	if act != glfw.Press {
		return
	}

	switch key {
	case glfw.KeyUp:
		m.menu.Previous()
	case glfw.KeyDown:
		m.menu.Next()
	case glfw.KeyLeft:
		m.menu.Decrease()
	case glfw.KeyRight:
		m.menu.Increase()
	case glfw.KeyEnter, glfw.KeyKPEnter:
		m.menu.Input('\n')
	case glfw.KeyBackspace:
		m.menu.Input('\b')
	case glfw.KeyDelete:
		m.menu.Input('\xFF')
	}
}

func (m *Menu) InputChar(r rune) { m.menu.Input(r) }

func (m *Menu) Prepare(ctx context.Context, now time.Time) {
	n := m.menu.Count()
	for i := 0; i < n; i++ {
		text := m.menu.Item(i).Text()
		desc := m.menu.Item(i).Description()
		m.strCache[2*i] = text
		m.strCache[2*i+1] = desc
	}

	if len(m.strCache2) != len(m.strCache) {
		m.strCache2 = make([]string, len(m.strCache))
	}

	same := true
	for i := range m.strCache {
		same = same && m.strCache[i] == m.strCache2[i]
		m.strCache2[i] = m.strCache[i]
	}

	if same {
		return
	}

	m.text.Prepare(m.strCache...)
}

func (m *Menu) Render(r *render.Renderer) {
	n := m.menu.Count()
	for i := 0; i < n; i++ {
		var color mgl32.Vec4
		text := m.menu.Item(i).Text()
		selected := i == m.menu.Current()
		if selected {
			color = colorSelected
		} else {
			color = colorItem
		}

		m.text.String(r, m.modelItem[i], color, text)
	}

	desc := m.menu.Item(m.menu.Current()).Description()
	m.text.String(r, m.modelDesc, colorDescription, desc)
}
