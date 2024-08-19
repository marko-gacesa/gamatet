// Copyright (c) 2020 by Marko Gaćeša

package main

import (
	"fmt"
	"gamatet/graphics/loop"
	"gamatet/internal/app"
	"gamatet/internal/config"
)

func main() {
	cfg, cfgPath := config.Load()

	app := app.NewApp(cfg, cfgPath)

	err := loop.Loop(app)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
}
