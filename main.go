// Copyright (c) 2020 by Marko Gaćeša

package main

import (
	"fmt"
	"gamatet/demo"
	"gamatet/graphics/render"

	_ "image/jpeg"
	_ "image/png"
)

func main() {
	demo.Single()
	return
	if err := render.Loop(); err != nil {
		fmt.Println(err)
	}
	return
}
