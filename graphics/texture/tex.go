// Copyright (c) 2020 by Marko Gaćeša

package texture

import (
	"image"
	"image/color"
)

func Tex2D(seed int64) image.Image {
	const size = 256

	values := Perlin2D(size, 8, 6, seed)
	Clamp(values, 0.0, 0.99999)

	img := image.NewGray(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			g := color.Gray{Y: byte(values[y*size+x] * 256)}
			img.SetGray(x, y, g)
		}
	}

	return img
}
