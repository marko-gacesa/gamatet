// Copyright (c) 2020 by Marko Gaćeša

package texture

import (
	"gamatet/graphics/symbol"
	"golang.org/x/image/font/gofont/gobold"
	"log"
)

var SymbolCache *symbol.Cache

func init() {
	var err error

	SymbolCache, err = symbol.NewCache(gobold.TTF, 110)
	if err != nil {
		log.Fatalf("failed to initialize font: %v" + err.Error())
	}
}
