// Copyright (c) 2020 by Marko Gaćeša

package main

import (
	"fmt"
	"gamatet/graphics/loop"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	if err := loop.Loop(); err != nil {
		fmt.Println(err)
	}
	return
}
