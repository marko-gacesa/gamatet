// Copyright (c) 2024 by Marko Gaćeša

package render

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
)

var FontNumerals *truetype.Font
var FontNormal *truetype.Font

func init() {
	var err error

	FontNumerals, err = truetype.Parse(gobold.TTF)
	if err != nil {
		panic(err)
	}

	FontNormal, err = truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}
}
