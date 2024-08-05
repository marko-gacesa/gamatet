// Copyright (c) 2024 by Marko Gaćeša

package render

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"strconv"
	"strings"
)

type FPS struct {
	value int
	frame int
	freq  uint64
	prev  uint64
}

func NewFPS() *FPS {
	return &FPS{freq: glfw.GetTimerFrequency()}
}

func (fps *FPS) Get() int {
	curr := glfw.GetTimerValue() % fps.freq
	if curr > fps.prev {
		fps.frame++
	} else {
		fps.value = fps.frame
		fps.frame = 1
	}
	fps.prev = curr
	return fps.value
}

func (fps *FPS) Render(r *Renderer, text *Text, model mgl32.Mat4) {
	v := fps.Get()
	if v == 0 {
		return
	}

	sb := strings.Builder{}
	sb.WriteString("fps=")
	sb.WriteString(strconv.Itoa(v))

	text.String(r, model, mgl32.Vec4{0.5, 0.5, 0, 0.7}, sb.String())
}
