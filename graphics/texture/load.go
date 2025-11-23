// Copyright (c) 2020-2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package texture

import (
	"fmt"
	"image"
	"image/draw"
	"os"
)

func LoadFile(file string) (image.Image, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %w", file, err)
	}

	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported texture %q stride=%d", file, rgba.Stride)
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Src)

	return rgba, nil
}
