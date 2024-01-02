// Copyright (c) 2020 by Marko Gaćeša

package rendercache

import (
	"gamatet/graphics/gtypes"
	"github.com/go-gl/mathgl/mgl32"
	"sync"
)

type pointLightsPool struct {
	pool sync.Pool
}

var PointLightPool = pointLightsPool{
	pool: sync.Pool{
		New: func() any {
			list := make([]gtypes.PointLight, 0, 16)
			return pointLights(list)
		},
	},
}

func (b *pointLightsPool) Get() pointLights {
	list := b.pool.Get().(pointLights)
	list = list[:0]
	return list
}

func (b *pointLightsPool) Put(list pointLights) {
	b.pool.Put(list)
}

type pointLights []gtypes.PointLight

func (p *pointLights) Add(position mgl32.Vec3, color mgl32.Vec3, intensity float32) {
	*p = append(*p, gtypes.PointLight{
		Position:  position,
		Color:     color,
		Intensity: intensity,
	})
}

func (p *pointLights) AddWithModel(model mgl32.Mat4, color mgl32.Vec3, intensity float32) {
	*p = append(*p, gtypes.PointLight{
		Position:  model.Mul4x1(mgl32.Vec4{0, 0, 0, 1}).Vec3(),
		Color:     color,
		Intensity: intensity,
	})
}
