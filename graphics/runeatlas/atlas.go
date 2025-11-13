// Copyright (c) 2024, 2025 by Marko Gaćeša

package runeatlas

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/marko-gacesa/gamatet/graphics/gutil"
)

type RuneAtlas struct {
	size     int
	image    *image.Alpha
	isDirty  bool
	hasHoles bool

	face *Face

	pos    fixed.Point26_6
	glyphs map[rune]image.Rectangle

	// mx is to make the RuneAtlas thread safe.
	// Also, because calls to font.Face (for glyph measurement) are not thread safe.
	mx *sync.Mutex
}

const (
	padW = 1
	padH = 1
)

func NewRuneAtlas(face *Face, size int) *RuneAtlas {
	size = gutil.CeilPow2(size)
	img := image.NewAlpha(image.Rect(0, 0, size, size))
	c := &RuneAtlas{
		size:     size,
		image:    img,
		isDirty:  true,
		hasHoles: false,
		face:     face,
		pos:      fixed.Point26_6{},
		glyphs:   map[rune]image.Rectangle{},
		mx:       &sync.Mutex{},
	}
	return c
}

func (a *RuneAtlas) Image() image.Image {
	return a.image
}

func (a *RuneAtlas) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	a.mx.Lock()
	defer a.mx.Unlock()

	return png.Encode(f, a.image)
}

func (a *RuneAtlas) IsDirty() bool {
	return a.isDirty
}

func (a *RuneAtlas) ClearDirty() {
	a.isDirty = false
}

func (a *RuneAtlas) KernToHeight(r0, r1 rune) float32 {
	a.mx.Lock()
	defer a.mx.Unlock()

	if r0 == 0 || r1 == 0 {
		return 0
	}
	kern := a.face.kern(r0, r1)
	return float32(kern) / float32(a.face.fullHeight)
}

func (a *RuneAtlas) Clear() {
	a.mx.Lock()
	defer a.mx.Unlock()

	a.clear()
	a.hasHoles = false
	a.isDirty = true
	a.pos = fixed.Point26_6{}
	for glyph := range a.glyphs {
		delete(a.glyphs, glyph)
	}
}

func (a *RuneAtlas) Remove(r rune) {
	a.mx.Lock()
	defer a.mx.Unlock()

	if _, ok := a.glyphs[r]; !ok {
		return
	}
	delete(a.glyphs, r)
	a.hasHoles = true
}

func (a *RuneAtlas) TextUV(r rune) (RectUV, bool) {
	rect, ok := a.Store(r)
	if !ok {
		return [4]float32{0, 0, 0, 0}, false
	}

	s := float32(a.size)
	x0 := float32(rect.Min.X) / s
	y0 := float32(rect.Min.Y) / s
	x1 := float32(rect.Max.X) / s
	y1 := float32(rect.Max.Y) / s

	return [4]float32{x0, y0, x1, y1}, true
}

func (a *RuneAtlas) Store(r rune) (image.Rectangle, bool) {
	a.mx.Lock()
	defer a.mx.Unlock()

	rect, ok := a.glyphs[r]
	if ok {
		return rect, true
	}

	rect, ok = a.draw(r)
	if !ok {
		if a.hasHoles {
			a.redrawAll()
			rect, ok = a.draw(r)
			if !ok {
				return image.Rectangle{}, false
			}
		} else {
			return image.Rectangle{}, false
		}
	}

	a.glyphs[r] = rect

	return rect, true
}

func (a *RuneAtlas) redrawAll() bool {
	a.clear()

	allOk := true

	for glyph := range a.glyphs {
		rect, ok := a.draw(glyph)
		a.glyphs[glyph] = rect
		allOk = allOk && ok
	}

	a.hasHoles = false

	return allOk
}

func (a *RuneAtlas) draw(r rune) (image.Rectangle, bool) {
	rectText := a.face.measure(r)

	height := a.face.fullHeight
	textDrawPoint := fixed.Point26_6{X: -rectText.Min.X, Y: a.face.protTop + a.face.baseHeight}
	rectBox := fixed.Rectangle26_6{
		Min: fixed.Point26_6{X: 0, Y: 0},
		Max: fixed.Point26_6{X: rectText.Max.X - rectText.Min.X, Y: height},
	}

	width := rectBox.Max.X - rectBox.Min.X
	overflowsLine := a.pos.X+width+fixed.I(padW) >= fixed.I(a.size)
	overflowsPage := a.pos.Y+height+fixed.I(padH) >= fixed.I(a.size)
	if overflowsLine && overflowsPage {
		return image.Rectangle{}, false
	}

	if overflowsLine {
		a.pos.X = 0
		a.pos.Y += height + fixed.I(padH)
	}

	rectBox = rectBox.Add(a.pos)
	textDrawPoint = textDrawPoint.Add(a.pos)

	dctx := freetype.NewContext()
	dctx.SetHinting(font.HintingFull)
	dctx.SetDPI(a.face.dpi)
	dctx.SetFont(a.face.font)
	dctx.SetFontSize(a.face.size)
	dctx.SetSrc(image.NewUniform(color.White))
	dctx.SetDst(a.image)
	dctx.SetClip(a.image.Bounds())

	rect := image.Rectangle{
		Min: image.Point{X: rectBox.Min.X.Floor(), Y: rectBox.Min.Y.Floor()},
		Max: image.Point{X: rectBox.Max.X.Ceil(), Y: rectBox.Max.Y.Ceil()},
	}

	//draw.Draw(a.image, rect, image.Black, image.Point{}, draw.Over)
	dctx.DrawString(string(r), textDrawPoint)

	a.pos.X += width + fixed.I(padW)
	a.isDirty = true

	return rect, true
}

func (a *RuneAtlas) clear() {
	a.pos = fixed.Point26_6{}
	draw.Draw(a.image, a.image.Rect, image.Transparent, image.Point{}, draw.Over)
}
