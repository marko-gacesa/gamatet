// Copyright (c) 2020, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package material

import "github.com/go-gl/mathgl/mgl32"

const (
	MaxLights = 16
)

type PointLight struct {
	Position  mgl32.Vec3
	Color     mgl32.Vec3
	Intensity float32
}

func setLights(uniPointLights int32, lights []PointLight) {
	n := min(int32(len(lights)), MaxLights)

	for i := range n {
		uniformVec3(uniPointLights+i*3, lights[i].Position)
		uniformVec3(uniPointLights+i*3+1, lights[i].Color)
		uniform1f(uniPointLights+i*3+2, lights[i].Intensity)
	}
	for i := n; i < MaxLights; i++ {
		uniform1f(uniPointLights+i*3+2, 0)
	}
}
