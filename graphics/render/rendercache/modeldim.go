// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package rendercache

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type modelDim struct {
	Model mgl32.Mat4
	Dim   mgl32.Vec2
}

type modelDimPool struct {
	pool sync.Pool
}

var ModelDimPool = modelDimPool{
	pool: sync.Pool{
		New: func() any {
			list := make([]modelDim, 0, 256)
			return (*ModelDimList)(&list)
		},
	},
}

func (b *modelDimPool) Get() *ModelDimList {
	list := b.pool.Get().(*ModelDimList)
	*list = (*list)[:0]
	return list
}

func (b *modelDimPool) Put(list *ModelDimList) {
	b.pool.Put(list)
}

type ModelDimList []modelDim

func (p *ModelDimList) Len() int {
	return len(*p)
}

func (p *ModelDimList) Add(model mgl32.Mat4, dim mgl32.Vec2) {
	*p = append(*p, modelDim{
		Model: model,
		Dim:   dim,
	})
}
