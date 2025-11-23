// Copyright (c) 2020-2023 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package texture

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func SavePNG(filename string, q image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create png file %q: %w", filename, err)
	}

	defer f.Close()

	err = png.Encode(f, q)
	if err != nil {
		return fmt.Errorf("failed to encode png image to %q: %w", filename, err)
	}

	return nil
}

func LoadPNG(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open png file %q: %w", filename, err)
	}

	defer f.Close()

	q, err := png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode png from file %q: %w", filename, err)
	}

	return q, nil
}
