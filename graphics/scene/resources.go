// Copyright (c) 2024 by Marko Gaćeša

package scene

import (
	"gamatet/graphics/render"
	"gamatet/graphics/texture"
	"gamatet/internal/types"
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

func (r Resources) Screen(ctx screen.Context, data any) screen.Screen {
	switch v := data.(type) {
	case *menu.Menu:
		return NewMenu(r.rend, r.tex, v)
	case types.GameParams:
		return NewGame(r.rend, r.tex, v)
	case types.GameOneParams:
		return NewGameOne(r.rend, r.tex, v)
	case types.GameDoubleParams:
		return NewGameDouble(r.rend, r.tex, v)
	}

	return nil
}
