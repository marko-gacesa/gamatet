// Copyright (c) 2024, 2025 by Marko Gaćeša

package render

import (
	"gamatet/logic/cache"
	"github.com/go-gl/glfw/v3.3/glfw"
	"strconv"
)

type FPS struct {
	cached *cache.String[int]
}

func NewFPS() FPS {
	freq := glfw.GetTimerFrequency()
	var frame uint64
	var prev uint64
	var value int
	return FPS{
		cached: cache.NewString[int](
			func() int {
				curr := glfw.GetTimerValue() % freq
				if curr > prev {
					frame++
				} else {
					value = int(frame)
					frame = 1
				}
				prev = curr
				return value
			},
			func(v1 int, v2 int) bool { return v1 == v2 },
			func(fps int) string { return "fps=" + strconv.Itoa(fps) },
			0,
		),
	}
}

func (f FPS) String() string { return f.cached.String() }
