// Copyright (c) 2024, 2025 by Marko Gaćeša

package render

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"strconv"
)

type FPS struct {
	value    int
	valueOld int
	valueStr string

	frame int
	freq  uint64
	prev  uint64
}

func NewFPS() *FPS {
	return &FPS{freq: glfw.GetTimerFrequency()}
}

func (fps *FPS) update() {
	curr := glfw.GetTimerValue() % fps.freq
	if curr > fps.prev {
		fps.frame++
	} else {
		fps.value = fps.frame
		if fps.valueOld != fps.value {
			fps.valueOld = fps.value
			fps.valueStr = "fps=" + strconv.Itoa(fps.value)
		}
		fps.frame = 1
	}
	fps.prev = curr
}

func (fps *FPS) Get() int {
	fps.update()
	return fps.value
}

func (fps *FPS) Render(r *Renderer, text *Text) {
	fps.update()
	s := fps.valueStr
	if s == "" {
		return
	}

	//tw, th := text.Dim(s)
	//_, _ = tw, th

	const contentW = 80
	const contentH = contentW * 9 / 16
	r.OrthogonalFull(contentW, contentH, contentW, contentH, 1)

	model := mgl32.Translate3D(-float32(contentW)/2, -float32(contentH)/2+0.5, 0)

	text.String(r, model, mgl32.Vec4{0.5, 0.5, 0, 0.7}, s)
}
