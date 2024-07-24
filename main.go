// Copyright (c) 2020 by Marko Gaćeša

package main

import (
	"fmt"
	"gamatet/demo"
	"gamatet/graphics/loop"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	demo.Single()
	return
	if err := loop.Loop(); err != nil {
		fmt.Println(err)
	}
	return
}
