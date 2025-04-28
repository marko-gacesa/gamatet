// Copyright (c) 2024, 2025 by Marko Gaćeša

package render

import (
	"github.com/go-gl/glfw/v3.3/glfw"
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
