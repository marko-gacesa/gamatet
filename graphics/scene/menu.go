// Copyright (c) 2024 by Marko Gaćeša

package scene

import (
	"context"
	"gamatet/graphics/render"
	"gamatet/graphics/texture"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"time"
	"unicode"
)

const (
	screenContentW      = 24
	screenContentH      = 16
	screenMaxShownItems = 8
)

var (
	screenModelItem [screenMaxShownItems]mgl32.Mat4
	screenModelDesc mgl32.Mat4
)

func init() {
	for i := 0; i < screenMaxShownItems; i++ {
		screenModelItem[i] = mgl32.Ident4().
			Mul4(mgl32.Translate3D(-screenContentW/2+1, float32(screenMaxShownItems-screenContentH/2+1-i), 0))
	}

	screenModelDesc = mgl32.Ident4().
		Mul4(mgl32.Translate3D(-screenContentW/2+1, float32(-screenContentH/2+0.5), 0)).
		Mul4(mgl32.Scale3D(0.6, 0.6, 0.6))
}

type Menu struct {
	renderer *render.Renderer
	tex      *texture.Manager
	text     render.Text

	stopper *screen.Stopper

	menu *menu.Menu

	strCache  []string
	strCache2 []string

	offset int
}

var _ screen.Screen = (*Menu)(nil)

var (
	colorSelected    = mgl32.Vec4{0.8, 0.7, 0.6, 1}
	colorItem        = colorSelected.Mul(0.8)
	colorDescription = colorSelected.Mul(0.6)
	colorTitle       = colorSelected.Mul(0.4)
)

func NewMenu(renderer *render.Renderer, tex *texture.Manager, m *menu.Menu) *Menu {
	text := render.MakeText(tex, render.Font)

	n := m.Count()
	strCache := make([]string, 2*n+1)

	return &Menu{
		renderer: renderer,
		tex:      tex,
		text:     *text,
		stopper:  screen.NewStopper(),
		menu:     m,
		strCache: strCache,
	}
}

func (m *Menu) Done() <-chan error { return m.stopper.Done() }

func (m *Menu) Release() {
	m.text.Release()
}

func (m *Menu) UpdateViewSize(w, h int) {
	m.renderer.OrthogonalFull(w, h, screenContentW, screenContentH, 2)
}

func (m *Menu) InputKeyPress(key, scancode int) {
	switch glfw.Key(key) {
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
	case glfw.KeyEscape:
		m.menu.Finish()
	}
}

func (m *Menu) InputChar(r rune) {
	if unicode.IsGraphic(r) {
		m.menu.Input(r)
	}
}

func (m *Menu) Prepare(ctx context.Context, now time.Time) {
	if m.menu.Finished() {
		m.stopper.Stop()
	}

	n := m.menu.Count()

	if len(m.strCache) != n {
		m.strCache = make([]string, 2*n+1)
	}

	m.strCache[0] = m.menu.Title()
	for i := 0; i < n; i++ {
		text := m.menu.Item(i).Text()
		desc := m.menu.Item(i).Description()
		m.strCache[2*i+1] = text
		m.strCache[2*i+2] = desc
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

func (m *Menu) Render(ctx context.Context) {
	modelIdxStart := 0
	modelIdxEnd := m.menu.Count()
	if modelIdxEnd > screenMaxShownItems {
		modelIdxEnd = screenMaxShownItems
	} else if modelIdxEnd < screenMaxShownItems {
		modelIdxStart = screenMaxShownItems - modelIdxEnd
		modelIdxEnd = screenMaxShownItems
	}

	idxSelected := m.menu.CurrentIdx() - m.offset
	if idxSelected < 0 {
		m.offset += idxSelected
		idxSelected = 0
	} else if idxSelected >= screenMaxShownItems {
		m.offset += idxSelected - screenMaxShownItems + 1
		idxSelected = screenMaxShownItems - 1
	}

	r := m.renderer

	for modelIdx, idx := modelIdxStart, 0; modelIdx < modelIdxEnd; modelIdx, idx = modelIdx+1, idx+1 {
		var color mgl32.Vec4
		text := m.menu.Item(idx + m.offset).Text()
		if idxSelected == idx {
			color = colorSelected
		} else {
			color = colorItem
		}

		m.text.String(r, screenModelItem[modelIdx], color, text)
	}

	desc := m.menu.Current().Description()
	m.text.String(r, screenModelDesc, colorDescription, desc)

	if title := m.menu.Title(); title != "" {
		modelTitle := screenModelItem[modelIdxStart].
			Mul4(mgl32.Translate3D(-0.5, 1.4, 0)).
			Mul4(mgl32.Scale3D(1.5, 1.5, 1))
		m.text.String(r, modelTitle, colorTitle, title)
	}
}
