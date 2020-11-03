// Copyright (c) 2020 by Marko Gaćeša

package texture

import (
	"image"
	"image/png"
	"os"
)

func SavePNG(filename string, image image.Image) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}

	defer f.Close()

	err = png.Encode(f, image)

	return
}
