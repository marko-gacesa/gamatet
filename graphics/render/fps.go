// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package render

import (
	"strconv"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/marko-gacesa/gamatet/logic/cache"
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
			func(v1 *int, v2 int) bool {
				equal := *v1 == v2
				*v1 = v2
				return equal
			},
			func(fps int) string { return "FPS=" + strconv.Itoa(fps) + "Hz" },
			0,
		),
	}
}

func (f FPS) String() string { return f.cached.String() }
