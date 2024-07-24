// Copyright (c) 2024 by Marko Gaćeša

package scene

import (
	"context"
	"gamatet/graphics/render"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type Object interface {
	Prepare(ctx context.Context, model *mgl32.Mat4, now time.Time)
	Render(r *render.Renderer)
}

type Objects []Object

func (objs Objects) Prepare(ctx context.Context, model *mgl32.Mat4, now time.Time) {
	for _, object := range objs {
		object.Prepare(ctx, model, now)
	}
}

func (objs Objects) Render(r *render.Renderer) {
	for _, object := range objs {
		object.Render(r)
	}
}
