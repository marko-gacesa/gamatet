// Copyright (c) 2023,2024 by Marko Gaćeša

package geometry

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

func makeSphereModel(radius float64, wSegs, hSegs int) []vertex {
	radius = max(0.001, radius)
	wSegs = max(3, wSegs)
	hSegs = max(2, hSegs)

	wDelta := (2 * math.Pi) / float64(wSegs)
	hDelta := math.Pi / float64(hSegs)

	result := make([]vertex, 0, wSegs*(hSegs-1)*6)
	cache := make(map[int]vertex, (wSegs+1)*(hSegs+1))

	spherePoint := func(i, j int) vertex {
		idx := i*(hSegs+1) + j // i is from 0..wSeg, j is from 0..hSeg (in total (wSeg+1)*(hSeg+1) different values)
		if v, ok := cache[idx]; ok {
			return v
		}

		wa := float64(i)*wDelta - math.Pi/2
		ha := float64(j) * hDelta

		px := -radius * math.Cos(wa) * math.Sin(ha)
		py := radius * math.Cos(ha)
		pz := radius * math.Sin(wa) * math.Sin(ha)

		p := mgl32.Vec3{float32(px), float32(py), float32(pz)}
		n := p.Normalize()
		uv := mgl32.Vec3{float32(i) / float32(wSegs), float32(j) / float32(hSegs), 1}

		v := gen(p, n, uv)
		cache[idx] = v

		return v
	}

	for i := 0; i < wSegs; i++ {
		result = append(result, spherePoint(i, 1))
		result = append(result, spherePoint(i, 0))
		result = append(result, spherePoint(i+1, 1))
		for j := 1; j < hSegs-1; j++ {
			v0 := spherePoint(i, j)
			v1 := spherePoint(i+1, j)
			v2 := spherePoint(i+1, j+1)
			v3 := spherePoint(i, j+1)
			result = append(result, v0)
			result = append(result, v1)
			result = append(result, v2)
			result = append(result, v0)
			result = append(result, v2)
			result = append(result, v3)
		}
		result = append(result, spherePoint(i, hSegs-1))
		result = append(result, spherePoint(i+1, hSegs-1))
		result = append(result, spherePoint(i, hSegs))
	}

	return result
}
