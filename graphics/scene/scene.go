// Copyright (c) 2024 by Marko Gaćeša

package scene

import (
	"context"
	"gamatet/graphics/render"
	"time"
)

type Object interface {
	StartPrepare(ctx context.Context, now time.Time)
	EndPrepare(ctx context.Context)
	Render(r *render.Renderer)
}

type Scene struct {
	Renderer *render.Renderer
	Objects  []Object
}

func (s *Scene) Render(ctx context.Context) {
	now := time.Now()

	for _, object := range s.Objects {
		object.StartPrepare(ctx, now)
	}

	for _, object := range s.Objects {
		object.EndPrepare(ctx)
		object.Render(s.Renderer)
	}
}
