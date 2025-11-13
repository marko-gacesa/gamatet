// Copyright (c) 2020-2024 by Marko Gaćeša

package rendercache

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type modelPool struct {
	pool sync.Pool
}

var ModelPool = modelPool{
	pool: sync.Pool{
		New: func() any {
			list := make([]mgl32.Mat4, 0, 256)
			return Models(list)
		},
	},
}

func (b *modelPool) Get() Models {
	list := b.pool.Get().(Models)
	list = list[:0]
	return list
}

func (b *modelPool) Put(list Models) {
	b.pool.Put(list)
}

type Models []mgl32.Mat4

func (p *Models) Add(model mgl32.Mat4) {
	*p = append(*p, model)
}
