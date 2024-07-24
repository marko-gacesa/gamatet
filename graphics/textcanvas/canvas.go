// Copyright (c) 2024 by Marko Gaćeša

package textcanvas

import (
	"crypto/md5"
	"encoding/binary"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type TextCanvas struct {
	size     int
	image    *image.RGBA
	isDirty  bool
	hasHoles bool

	pos  fixed.Point26_6
	yMax fixed.Int26_6

	entries []*entry
}

type entry struct {
	key   string
	text  string
	color color.Color
	face  Face
	rect  image.Rectangle
	lrPad bool
}

var hashFunc = md5.New

const (
	padW = 1
	padH = 1
)

func NewTextCanvas(size int) *TextCanvas {
	rect := image.Rect(0, 0, size-1, size-1)
	img := image.NewRGBA(rect)
	c := &TextCanvas{
		size:     size,
		image:    img,
		isDirty:  true,
		hasHoles: false,
		pos:      fixed.Point26_6{},
		yMax:     0,
		entries:  nil,
	}
	return c
}

func (c *TextCanvas) Image() *image.RGBA {
	return c.image
}

func (c *TextCanvas) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, c.image)
}

func (c *TextCanvas) IsDirty() bool {
	return c.isDirty
}

func (c *TextCanvas) ClearDirty() {
	c.isDirty = false
}

func (c *TextCanvas) Clear() {
	c.clear()
	c.hasHoles = false
	c.isDirty = true
	c.pos = fixed.Point26_6{}
	c.yMax = 0
	c.entries = c.entries[:0]
}

func (c *TextCanvas) Remove(text string, face Face, color color.Color, lrPad bool) {
	key := makeKey(text, face, color, lrPad)
	for i, e := range c.entries {
		if e.key == key {
			c.entries = append(c.entries[:i], c.entries[i+1:]...)
			c.hasHoles = true
			return
		}
	}
}

func (c *TextCanvas) TextUV(text string, face Face, color color.Color, lrPad bool) RectUV {
	e := c.store(text, face, color, lrPad)
	if e == nil {
		return [4]float32{0, 0, 0, 0}
	}

	s := float32(c.size)
	x0 := float32(e.rect.Min.X) / s
	y0 := float32(e.rect.Min.Y) / s
	x1 := float32(e.rect.Max.X) / s
	y1 := float32(e.rect.Max.Y) / s

	return [4]float32{x0, y0, x1, y1}
}

func (c *TextCanvas) TextRect(text string, face Face, color color.Color, lrPad bool) image.Rectangle {
	e := c.store(text, face, color, lrPad)
	if e == nil {
		return image.Rectangle{}
	}
	return e.rect
}

func (c *TextCanvas) store(text string, face Face, color color.Color, lrPad bool) *entry {
	key := makeKey(text, face, color, lrPad)
	if e := c.find(key); e != nil {
		return e
	}

	rect, ok := c.draw(text, face, color, lrPad)
	if !ok {
		if c.hasHoles {
			c.redrawAll()
			rect, ok = c.draw(text, face, color, lrPad)
			if !ok {
				return nil
			}
		} else {
			return nil
		}
	}

	e := &entry{
		key:   key,
		text:  text,
		color: color,
		face:  face,
		rect:  rect,
		lrPad: lrPad,
	}

	c.entries = append(c.entries, e)

	return e
}

func (c *TextCanvas) redrawAll() bool {
	c.clear()

	allOk := true

	for _, e := range c.entries {
		rect, ok := c.draw(e.text, e.face, e.color, e.lrPad)
		e.rect = rect
		allOk = allOk && ok
	}

	c.hasHoles = false

	return allOk
}

func (c *TextCanvas) draw(text string, face Face, color color.Color, lrPad bool) (image.Rectangle, bool) {
	rectText := face.measure(text)

	var rectBox fixed.Rectangle26_6
	var textDrawPoint fixed.Point26_6

	if lrPad {
		textDrawPoint = fixed.Point26_6{X: face.protLeft, Y: face.protTop + face.baseHeight}
		rectBox = fixed.Rectangle26_6{
			Min: fixed.Point26_6{X: 0, Y: 0},
			Max: fixed.Point26_6{X: rectText.Max.X + 2*face.protLeft, Y: face.fullHeight},
		}
	} else {
		textDrawPoint = fixed.Point26_6{X: -rectText.Min.X, Y: face.protTop + face.baseHeight}
		rectBox = fixed.Rectangle26_6{
			Min: fixed.Point26_6{X: 0, Y: 0},
			Max: fixed.Point26_6{X: rectText.Max.X - rectText.Min.X, Y: face.fullHeight},
		}
	}

	width := rectBox.Max.X - rectBox.Min.X
	overflowsLine := c.pos.X+width+fixed.I(padW) >= fixed.I(c.size)
	overflowsPage := c.pos.Y+c.yMax+fixed.I(padH) >= fixed.I(c.size)
	if overflowsLine && overflowsPage {
		return image.Rectangle{}, false
	}

	c.yMax = max(c.yMax, face.fullHeight)

	if overflowsLine {
		c.pos.X = 0
		c.pos.Y += c.yMax + fixed.I(padH)
		c.yMax = 0
	}

	rectBox = rectBox.Add(c.pos)
	textDrawPoint = textDrawPoint.Add(c.pos)

	dctx := freetype.NewContext()
	dctx.SetHinting(font.HintingFull)
	dctx.SetDPI(face.dpi)
	dctx.SetFont(face.font)
	dctx.SetFontSize(face.size)
	dctx.SetSrc(image.NewUniform(color))
	dctx.SetDst(c.image)
	dctx.SetClip(c.image.Bounds())

	rect := image.Rectangle{
		Min: image.Point{X: rectBox.Min.X.Floor(), Y: rectBox.Min.Y.Floor()},
		Max: image.Point{X: rectBox.Max.X.Ceil(), Y: rectBox.Max.Y.Ceil()},
	}

	//draw.Draw(c.image, rect, image.Black, image.Point{}, draw.Over)
	dctx.DrawString(text, textDrawPoint)

	c.pos.X += width + fixed.I(padW)
	c.isDirty = true

	return rect, true
}

func (c *TextCanvas) clear() {
	c.pos = fixed.Point26_6{}
	draw.Draw(c.image, c.image.Rect, image.Transparent, image.Point{}, draw.Over)
}

func (c *TextCanvas) find(key string) *entry {
	for _, e := range c.entries {
		if e.key == key {
			return e
		}
	}

	return nil
}

func makeKey(text string, face Face, color color.Color, lrPad bool) string {
	sum := hashFunc()

	sum.Write(face.hashSum)
	sum.Write([]byte(text))

	var shade [4]byte
	r, g, b, a := color.RGBA()
	binary.LittleEndian.PutUint32(shade[:], r)
	sum.Write(shade[:])
	binary.LittleEndian.PutUint32(shade[:], g)
	sum.Write(shade[:])
	binary.LittleEndian.PutUint32(shade[:], b)
	sum.Write(shade[:])
	binary.LittleEndian.PutUint32(shade[:], a)
	sum.Write(shade[:])

	if lrPad {
		sum.Write([]byte{1})
	}

	hashSum := sum.Sum(nil)
	key := string(hashSum)

	return key
}
