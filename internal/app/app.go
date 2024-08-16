// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"fmt"
	"gamatet/graphics/loop"
	"gamatet/internal/config"
	"gamatet/internal/router"
)

type App struct {
	Config         config.Config
	ConfigFileName string
}

func Run() {
	cfg, cfgFN := config.Load()

	var app App
	app.Config = cfg
	app.ConfigFileName = cfgFN

	router := router.NewRouter(&app.Config)

	if err := loop.Loop(&app.Config, router); err != nil {
		fmt.Println(err)
	}
}
