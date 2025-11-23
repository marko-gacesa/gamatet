// Copyright (c) 2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package scene

import (
	"github.com/marko-gacesa/gamatet/graphics/render"
	"github.com/marko-gacesa/gamatet/graphics/texture"
	"github.com/marko-gacesa/gamatet/internal/types"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
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
	}

	return nil
}
