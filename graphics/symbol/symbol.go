// Copyright (c) 2020 by Marko Gaćeša

package symbol

import (
	"errors"
	"gamatet/graphics/gutil"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"sync"
	"unicode"
)

type Cache struct {
	ttf  *truetype.Font
	face font.Face

	size float64
	dpi  float64

	dim int
	bg  image.Image
	fg  image.Image

	drawPos    fixed.Point26_6
	baseBounds fixed.Rectangle26_6

	mx    sync.Mutex
	cache map[rune]image.Image
}

func NewCache(data []byte, size float64) (c *Cache, err error) {
	const dpi = 72

	ttf, err := truetype.Parse(data)
	if err != nil {
		return
	}

	face := truetype.NewFace(ttf, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	var baseBounds fixed.Rectangle26_6

	// calculate size of the base bound box by analysing the widest block shaped letters

	for _, r := range "WM" {
		rBounds, _, ok := face.GlyphBounds(r)
		if !ok {
			err = errors.New("unsupported font")
			return
		}

		baseBounds = baseBounds.Union(rBounds)
	}

	baseW := baseBounds.Max.X - baseBounds.Min.X
	baseH := baseBounds.Max.Y - baseBounds.Min.Y

	// find max top and bottom protrusions by analysing the letters with the largest protrusions from the base block rectangle

	var boundsProtruded fixed.Rectangle26_6
	for _, r := range "QqÅgpJWMiŠ" {
		rBounds, _, ok := face.GlyphBounds(r)
		if !ok {
			err = errors.New("unsupported font")
			return
		}
		boundsProtruded = boundsProtruded.Union(rBounds)
	}

	topPad := -boundsProtruded.Min.Y + baseBounds.Min.Y
	bottomPad := boundsProtruded.Max.Y - baseBounds.Max.Y
	maxH := baseH + topPad + bottomPad

	// calculate image dimension

	var dim int

	if baseW > maxH {
		dim = baseW.Ceil()
	} else {
		dim = maxH.Ceil()
	}

	dim = gutil.CeilPow2(dim)

	// calculate draw position

	var drawPos fixed.Point26_6

	drawPos.X = (fixed.I(dim)-baseW)/2 - baseBounds.Min.X
	drawPos.Y = (fixed.I(dim) + baseH) / 2

	c = &Cache{
		ttf:        ttf,
		face:       face,
		size:       size,
		dpi:        dpi,
		dim:        dim,
		bg:         image.White,
		fg:         image.Black,
		drawPos:    drawPos,
		baseBounds: baseBounds,
		mx:         sync.Mutex{},
		cache:      make(map[rune]image.Image, 128),
	}

	return
}

func (c *Cache) Dim() int {
	return c.dim
}

func (c *Cache) SetFG(fg image.Image) {
	c.fg = fg
}

func (c *Cache) SetBG(bg image.Image) {
	c.bg = bg
}

func (c *Cache) PreCache(from, to rune) {
	dctx := c._drawContext()
	for r := from; r < to; r++ {
		c.cache[r] = c._symbol(dctx, r)
	}
}

func (c *Cache) PreCacheLatinAndDigits() {
	dctx := c._drawContext()
	for r := '0'; r < '9'; r++ {
		c.cache[r] = c._symbol(dctx, r)
	}
	for r := 'a'; r < 'z'; r++ {
		c.cache[r] = c._symbol(dctx, r)
	}
	for r := 'A'; r < 'Z'; r++ {
		c.cache[r] = c._symbol(dctx, r)
	}
}

func (c *Cache) Symbol(r rune) image.Image {
	if unicode.IsControl(r) || unicode.IsSpace(r) {
		return nil
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	if img, ok := c.cache[r]; ok {
		return img
	}

	img := c._symbol(c._drawContext(), r)
	if img == nil {
		return nil
	}

	c.cache[r] = img

	return img
}

func (c *Cache) _drawContext() *freetype.Context {
	dctx := freetype.NewContext()
	dctx.SetDPI(c.dpi)
	dctx.SetFont(c.ttf)
	dctx.SetFontSize(c.size)
	dctx.SetSrc(c.fg)
	dctx.SetHinting(font.HintingFull)
	return dctx
}

func (c *Cache) _symbol(dctx *freetype.Context, r rune) image.Image {
	bounds, _, ok := c.face.GlyphBounds(r)
	if !ok {
		return nil
	}

	drawPos := c.drawPos
	drawPos.X += (c.baseBounds.Max.X - c.baseBounds.Min.X - bounds.Max.X + bounds.Min.X - bounds.Min.X) / 2

	img := image.NewGray(image.Rect(0, 0, c.dim, c.dim))
	draw.Draw(img, img.Bounds(), c.bg, image.Point{}, draw.Src)

	/*
		// draw base bounding box
		b := image.Rect(
			(c.baseBounds.Min.X + c.drawPos.X).Floor(), (c.baseBounds.Min.Y + c.drawPos.Y).Floor(),
			(c.baseBounds.Max.X + c.drawPos.X).Ceil(), (c.baseBounds.Max.Y + c.drawPos.Y).Ceil())
		draw.Draw(img, b, image.NewUniform(color.Gray{Y: 192}), image.Point{}, draw.Src)

		// draw symbol bounding box
		b = image.Rect(
			(bounds.Min.X + drawPos.X).Floor(), (bounds.Min.Y + drawPos.Y).Floor(),
			(bounds.Max.X + drawPos.X).Ceil(), (bounds.Max.Y + drawPos.Y).Ceil())
		draw.Draw(img, b, image.NewUniform(color.Gray{Y: 224}), image.Point{}, draw.Src)
	*/

	dctx.SetClip(img.Bounds())
	dctx.SetDst(img)

	_, err := dctx.DrawString(string(r), drawPos)
	if err != nil {
		return nil
	}

	return img
}

func (c *Cache) Empty() image.Image {
	img := image.NewGray(image.Rect(0, 0, c.dim, c.dim))
	draw.Draw(img, img.Bounds(), c.bg, image.Point{}, draw.Src)
	return img
}
