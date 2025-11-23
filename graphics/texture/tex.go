// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package texture

import (
	"image"
	"image/color"
	"math/rand"
)

func GrayTex(seed int64) *image.Gray {
	const size = 256

	values := Perlin2D(size, 16, seed)
	symXY(values, size)
	clamp(values, 0.0, 0.99999)

	img := image.NewGray(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			g := color.Gray{Y: byte(values[y*size+x] * 256)}
			img.SetGray(x, y, g)
		}
	}

	return img
}

func Chain1(link *image.Gray) *image.Gray {
	const size = 256

	rect := link.Bounds()
	w := rect.Dx()
	h := rect.Dy()

	img := image.NewGray(image.Rect(0, 0, size, size))
	y := size/2 - h/2
	for x := 0; x < size; x += w {
		drawXY(x, y, link, img)
		drawXYVert(y, x, link, img)
	}

	return img
}

func Chain2(link *image.Gray) *image.Gray {
	const size = 256

	rect := link.Bounds()
	w := rect.Dx()
	h := rect.Dy()

	img := image.NewGray(image.Rect(0, 0, size, size))
	y1 := size/2 - h/2 - h - 1
	y2 := size/2 - h/2 + h + 1
	for x := 0; x < size; x += w {
		drawXY(x, y1, link, img)
		drawXY(x, y2, link, img)
		drawXYVert(y1, x, link, img)
		drawXYVert(y2, x, link, img)
	}

	return img
}

func Chain3(link *image.Gray) *image.Gray {
	const size = 256

	rect := link.Bounds()
	w := rect.Dx()
	h := rect.Dy()

	img := image.NewGray(image.Rect(0, 0, size, size))
	y0 := size/2 - h/2
	y1 := size/2 - h/2 - 2*h - 1
	y2 := size/2 - h/2 + 2*h + 1
	for x := 0; x < size; x += w {
		drawXY(x, y0, link, img)
		drawXY(x, y1, link, img)
		drawXY(x, y2, link, img)
		drawXYVert(y0, x, link, img)
		drawXYVert(y1, x, link, img)
		drawXYVert(y2, x, link, img)
	}

	return img
}

func Link(seed int64) *image.Gray {
	const (
		w = 64
		h = 32
	)

	r := rand.New(rand.NewSource(seed))

	img := image.NewGray(image.Rect(0, 0, w, h))

	part := func(x, y, l, a int) {
		for i := range l {
			v := byte(r.Intn(16) + a)
			img.SetGray(x+i, y, color.Gray{v})
		}
	}

	symPart := func(x, y, l, a int) {
		part(x, y, l, a)
		part(x, h-1-y, l, a)
		part(w-x-l, y, l, a)
		part(w-x-l, h-1-y, l, a)
	}

	symPart(0, 4, 17, 160)
	symPart(0, 5, 19, 176)
	symPart(0, 6, 20, 192)
	symPart(0, 7, 20, 192)
	symPart(0, 8, 20, 176)
	symPart(0, 9, 20, 176)
	symPart(15, 10, 5, 160)
	symPart(16, 11, 4, 144)
	symPart(14, 12, 18, 208)
	symPart(13, 13, 19, 208)
	symPart(12, 14, 20, 224)
	symPart(12, 15, 20, 224)

	isEmpty := func(x, y int) bool {
		return img.GrayAt(x, y).Y == 0
	}

	isLink := func(x, y int) bool {
		if x < 0 || y < 0 || x >= w || y >= h {
			return false
		}
		return !isEmpty(x, y) && img.GrayAt(x, y).Y > 128
	}

	for x := range w {
		for y := range h {
			if isEmpty(x, y) && (isLink(x-1, y) || isLink(x+1, y) || isLink(x, y-1) || isLink(x, y+1)) {
				v := byte(r.Intn(64) + 32)
				img.SetGray(x, y, color.Gray{v})
			}
		}
	}

	for x := range w {
		for y := range h {
			if isEmpty(x, y) {
				img.SetGray(x, y, color.Gray{1})
			}
		}
	}

	return img
}

func drawXY(x, y int, src, target *image.Gray) {
	srcRect := src.Bounds()
	for i := 0; i < srcRect.Dx(); i++ {
		for j := 0; j < srcRect.Dy(); j++ {
			s := src.GrayAt(i, j)
			t := target.GrayAt(x+i, y+j)
			if t.Y < s.Y {
				target.SetGray(x+i, y+j, s)
			}
		}
	}
}

func drawXYVert(x, y int, src, target *image.Gray) {
	srcRect := src.Bounds()
	for i := 0; i < srcRect.Dy(); i++ {
		for j := 0; j < srcRect.Dx(); j++ {
			s := src.GrayAt(j, i)
			t := target.GrayAt(x+i, y+j)
			if t.Y < s.Y {
				target.SetGray(x+i, y+j, s)
			}
		}
	}
}
