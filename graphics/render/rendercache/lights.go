// Copyright (c) 2020-2024 by Marko Gaćeša

package rendercache

import (
	"gamatet/graphics/material"
	"github.com/go-gl/mathgl/mgl32"
	"sync"
)

type PointLightsPool struct {
	pool sync.Pool
}

var PointLightPool = PointLightsPool{
	pool: sync.Pool{
		New: func() any {
			list := make([]material.PointLight, 0, 16)
			return PointLights(list)
		},
	},
}

func (b *PointLightsPool) Get() PointLights {
	list := b.pool.Get().(PointLights)
	list = list[:0]
	return list
}

func (b *PointLightsPool) Put(list PointLights) {
	b.pool.Put(list)
}

type PointLights []material.PointLight

func (p *PointLights) Add(position mgl32.Vec3, color mgl32.Vec3, intensity float32) {
	*p = append(*p, material.PointLight{
		Position:  position,
		Color:     color,
		Intensity: intensity,
	})
}

func (p *PointLights) AddWithModel(model mgl32.Mat4, color mgl32.Vec3, intensity float32) {
	*p = append(*p, material.PointLight{
		Position:  model.Mul4x1(mgl32.Vec4{0, 0, 0, 1}).Vec3(),
		Color:     color,
		Intensity: intensity,
	})
}
