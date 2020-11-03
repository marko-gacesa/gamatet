// Copyright (c) 2020 by Marko Gaćeša

package rendercache

import (
	"github.com/go-gl/mathgl/mgl32"
	"sync"
)

type modelPool struct {
	pool sync.Pool
}

var ModelPool = modelPool{
	pool: sync.Pool{
		New: func() any {
			list := make([]mgl32.Mat4, 0, 256)
			return models(list)
		},
	},
}

func (b *modelPool) Get() models {
	list := b.pool.Get().(models)
	list = list[:0]
	return list
}

func (b *modelPool) Put(list models) {
	b.pool.Put(list)
}

type models []mgl32.Mat4

func (p *models) Add(model mgl32.Mat4) {
	*p = append(*p, model)
}
