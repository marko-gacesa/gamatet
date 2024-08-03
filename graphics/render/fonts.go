// Copyright (c) 2024 by Marko Gaćeša

package render

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"
)

var Font *truetype.Font

func init() {
	var err error

	Font, err = truetype.Parse(gobold.TTF)
	if err != nil {
		panic(err)
	}
}
