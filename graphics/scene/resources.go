// Copyright (c) 2024 by Marko Gaćeša

package scene

import (
	"context"
	"gamatet/game"
	"gamatet/game/core"
	"gamatet/graphics/render"
	"gamatet/graphics/scene/demoblocks"
	"gamatet/graphics/scene/fieldtest"
	"gamatet/graphics/texture"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
)

type Resources struct {
	rend *render.Renderer
	tex  *texture.Manager
}

func InitResources() *Resources {
	return &Resources{
		rend: render.NewRenderer(),
		tex:  texture.Init(), // Texture manager
	}
}

func (r Resources) Screen(ctx context.Context, data any) screen.Screen {
	switch v := data.(type) {
	case *menu.Menu:
		return NewMenu(r.rend, r.tex, v)
	case core.GameOneParams:
		//return NewGameOne(r.rend, r.tex, v)
		m := ScreenMap{}
		m["g"] = NewGameOne(r.rend, r.tex, v)
		m["h"] = NewHud(r.rend, r.tex)
		return m
	case core.GameDoubleParams:
		return NewGameDouble(r.rend, r.tex, v)
	case DemoScreenConfig:
		switch v.Name {
		case DemoBlocks:
			return demoblocks.NewDemoBlocks(r.rend, r.tex, v.Stop)
		case DemoFields:
			gameHost, gameInterpreter, playerInCh, waitDoneCh := game.NewFieldTest(ctx, 10, 22, v.Stop)
			return fieldtest.NewFieldTest(r.rend, r.tex, playerInCh, gameHost, gameInterpreter, waitDoneCh)
		}
	}

	return nil
}
