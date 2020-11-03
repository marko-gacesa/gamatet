// Copyright (c) 2020 by Marko Gaćeša

package rendercache

import (
	"gamatet/graphics/gtypes"
	"github.com/go-gl/mathgl/mgl32"
	"sync"
)

type modelColorPool struct {
	pool sync.Pool
}

var ModelColorPool = modelColorPool{
	pool: sync.Pool{
		New: func() any {
			list := make([]gtypes.ModelColor, 0, 256)
			return modelColorList(list)
		},
	},
}

func (b *modelColorPool) Get() modelColorList {
	list := b.pool.Get().(modelColorList)
	list = list[:0]
	return list
}

func (b *modelColorPool) Put(list modelColorList) {
	b.pool.Put(list)
}

type modelColorList []gtypes.ModelColor

func (p *modelColorList) Add(model mgl32.Mat4, color mgl32.Vec4) {
	*p = append(*p, gtypes.ModelColor{
		Model: model,
		Color: color,
	})
}
