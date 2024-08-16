// Copyright (c) 2024 by Marko Gaćeša

package screen

import (
	"context"
	"gamatet/graphics/render"
	"github.com/go-gl/glfw/v3.3/glfw"
	"time"
)

type Object interface {
	Prepare(ctx context.Context, now time.Time)
	Render(r *render.Renderer)
}

type Screen interface {
	SetCamera(r *render.Renderer, w, h int)
	Release()
	Shutdown()

	InputKey(key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey)
	InputChar(char rune)

	Object
}
