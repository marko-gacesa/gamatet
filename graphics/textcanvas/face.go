// Copyright (c) 2024 by Marko Gaćeša

package textcanvas

import (
	"encoding/binary"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"math"
)

type Face struct {
	font    *truetype.Font
	face    font.Face
	dpi     float64
	size    float64
	hashSum []byte
	// character height and protrusions
	baseHeight fixed.Int26_6
	protTop    fixed.Int26_6
	protBottom fixed.Int26_6
	protLeft   fixed.Int26_6
	fullHeight fixed.Int26_6
}

func NewFace(ttf *truetype.Font, size float64, dpi float64) Face {
	face := truetype.NewFace(ttf, &truetype.Options{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	// Get the font face's hash

	sum := hashFunc()
	sum.Write([]byte(ttf.Name(truetype.NameIDFontFullName)))
	var binFloat [8]byte
	binary.LittleEndian.PutUint64(binFloat[:], math.Float64bits(size))
	sum.Write(binFloat[:])
	binary.LittleEndian.PutUint64(binFloat[:], math.Float64bits(dpi))
	sum.Write(binFloat[:])

	hashSum := sum.Sum(nil)

	// Calculate size of the font's base height

	baseBounds, _, _ := face.GlyphBounds('A')
	baseHeight := baseBounds.Max.Y - baseBounds.Min.Y

	// Find max top and bottom protrusions by analysing the letters
	// with the largest protrusions from the base block rectangle.
	// We interested in the top, bottom and left protrusion.
	// The protrusion on the right side is of less interest
	// because the text is written to that side anyway.

	var boundsProtruded fixed.Rectangle26_6
	for _, r := range "ğđjÅŠßq" {
		rBounds, _, ok := face.GlyphBounds(r)
		if !ok {
			continue
		}
		boundsProtruded = boundsProtruded.Union(rBounds)
	}

	protTop := -boundsProtruded.Min.Y + baseBounds.Min.Y
	protBottom := boundsProtruded.Max.Y - baseBounds.Max.Y
	protLeft := -boundsProtruded.Min.X

	return Face{
		font:       ttf,
		face:       face,
		dpi:        dpi,
		size:       size,
		hashSum:    hashSum,
		baseHeight: baseHeight,
		protTop:    protTop,
		protBottom: protBottom,
		protLeft:   protLeft,
		fullHeight: baseHeight + protTop + protBottom,
	}
}

func (face *Face) measure(text string) fixed.Rectangle26_6 {
	var textBounds fixed.Rectangle26_6
	var xPos fixed.Int26_6
	for _, r := range text {
		rBounds, advance, ok := face.face.GlyphBounds(r)
		if !ok {
			continue // skip unsupported
		}

		rBounds.Min.X += xPos
		rBounds.Max.X += xPos
		xPos += advance

		textBounds = textBounds.Union(rBounds)
	}

	return textBounds
}
