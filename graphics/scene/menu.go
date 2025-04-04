// Copyright (c) 2024,2025 by Marko Gaćeša

package scene

import (
	"gamatet/graphics/render"
	"gamatet/graphics/scene/base"
	"gamatet/graphics/texture"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"time"
	"unicode"
)

const (
	screenContentW      = 32
	screenContentH      = 20
	screenMaxShownItems = 16
	currItemMarker      = '●' //'◦'
)

var (
	screenModelItem   [screenMaxShownItems]mgl32.Mat4
	screenModelMarker [screenMaxShownItems]mgl32.Mat4
	screenModelDesc   mgl32.Mat4
)

func init() {
	for i := 0; i < screenMaxShownItems; i++ {
		screenModelMarker[i] = mgl32.Ident4().
			Mul4(mgl32.Translate3D(-screenContentW/2+1, float32(screenMaxShownItems-screenContentH/2+1-i), 0))
		screenModelItem[i] = screenModelMarker[i].
			Mul4(mgl32.Translate3D(0.4, 0, 0))
	}

	screenModelDesc = mgl32.Ident4().
		Mul4(mgl32.Translate3D(-screenContentW/2+1, float32(-screenContentH/2+0.5), 0)).
		Mul4(mgl32.Scale3D(0.6, 0.6, 0.6))
}

type Menu struct {
	base.Base
	text render.Text

	menu *menu.Menu

	strCache  []string
	strCache2 []string

	offset        int
	selectedColor mgl32.Vec4
}

var _ screen.Screen = (*Menu)(nil)

var (
	colorItemSelected         = mgl32.Vec4{0.8, 0.7, 0.6, 1}
	colorItem                 = colorItemSelected.Mul(0.8)
	colorItemDisabledSelected = mgl32.Vec4{0.4, 0.4, 0.4, 1}
	colorItemDisabled         = colorItemDisabledSelected.Mul(0.8)
	colorDescription          = colorItemSelected.Mul(0.6)
	colorTitle                = colorItemSelected.Mul(0.4)
)

func NewMenu(renderer *render.Renderer, tex *texture.Manager, m *menu.Menu) *Menu {
	text := render.MakeText(tex, render.Font)

	text.Prepare(string(currItemMarker))

	n := m.Count()
	strCache := make([]string, 2*n+1)

	return &Menu{
		Base:     base.NewBase(renderer, tex),
		text:     *text,
		menu:     m,
		strCache: strCache,
	}
}

func (m *Menu) Release() {
	m.text.Release()
}

func (m *Menu) UpdateViewSize(w, h int) {
	m.Renderer().OrthogonalFull(w, h, screenContentW, screenContentH, 2)
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
		m.menu.Input(menu.InputEnter)
	case glfw.KeyBackspace:
		m.menu.Input(menu.InputBackspace)
	case glfw.KeyDelete:
		m.menu.Input(menu.InputDelete)
	case glfw.KeyEscape:
		m.menu.Input(menu.InputEscape)
	}
}

func (m *Menu) InputChar(r rune) {
	if unicode.IsGraphic(r) {
		m.menu.Input(r)
	}
}

func (m *Menu) Prepare(now time.Time) {
	t64 := math.Sin(4 * glfw.GetTime())
	t := float32(t64 * t64)
	m.selectedColor = colorItem.Mul(t).Add(colorItemSelected.Mul(1 - t))

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

func (m *Menu) Render() {
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

	for modelIdxEnd-modelIdxStart+m.offset > m.menu.Count() && m.offset > 0 {
		if idxSelected > 0 {
			idxSelected--
		}
		m.offset--
	}

	r := m.Renderer()

	for modelIdx, idx := modelIdxStart, 0; modelIdx < modelIdxEnd; modelIdx, idx = modelIdx+1, idx+1 {
		item := m.menu.Item(idx + m.offset)
		text := item.Text()
		disabled := item.IsDisabled()

		model := screenModelItem[modelIdx]

		var color mgl32.Vec4
		if idxSelected == idx {
			if disabled {
				color = colorItemDisabledSelected
			} else {
				color = colorItemSelected
			}
			m.text.Rune(r, screenModelMarker[modelIdx], m.selectedColor, currItemMarker)
		} else {
			if disabled {
				color = colorItemDisabled
			} else {
				color = colorItem
			}
		}

		m.text.String(r, model, color, text)
	}

	desc := m.menu.Description()
	m.text.String(r, screenModelDesc, colorDescription, desc)

	screenModeTitle := screenModelItem[screenMaxShownItems-1]
	if modelIdxEnd != modelIdxStart {
		screenModeTitle = screenModelMarker[modelIdxStart]
	}

	if title := m.menu.Title(); title != "" {
		modelTitle := screenModeTitle.
			Mul4(mgl32.Translate3D(-0.5, 1.4, 0)).
			Mul4(mgl32.Scale3D(1.5, 1.5, 1))
		m.text.String(r, modelTitle, colorTitle, title)
	}
}
