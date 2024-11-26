// Copyright (c) 2024 by Marko Gaćeša

package scene

import (
	"context"
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
		return NewGameOne(r.rend, r.tex, v)
	case string:
		switch v {
		case "test-blocks":
			return demoblocks.NewDemoBlocks(r.rend, r.tex)
		case "test-fields":
			return fieldtest.NewFieldTest(ctx, r.rend, r.tex)
		}
	}

	return nil
}
