// Copyright (c) 2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package rendercache

import (
	"cmp"
	"sort"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type modelColorValue[T cmp.Ordered] struct {
	Model mgl32.Mat4
	Color mgl32.Vec4
	Value T
}

type modelColorValuePool[T cmp.Ordered] struct {
	pool sync.Pool
}

var ModelColorIntPool = newModelColorValuePool[int]()
var ModelColorStringPool = newModelColorValuePool[string]()

func newModelColorValuePool[T cmp.Ordered]() modelColorValuePool[T] {
	return modelColorValuePool[T]{
		pool: sync.Pool{
			New: func() any {
				list := make([]modelColorValue[T], 0, 256)
				return ModelColorValueList[T](list)
			},
		},
	}
}

func (b *modelColorValuePool[T]) Get() ModelColorValueList[T] {
	list := b.pool.Get().(ModelColorValueList[T])
	list = list[:0]
	return list
}

func (b *modelColorValuePool[T]) Put(list ModelColorValueList[T]) {
	b.pool.Put(list)
}

type ModelColorValueList[T cmp.Ordered] []modelColorValue[T]

func (p *ModelColorValueList[T]) Add(model mgl32.Mat4, color mgl32.Vec4, value T) {
	*p = append(*p, modelColorValue[T]{
		Model: model,
		Color: color,
		Value: value,
	})
}

func (p *ModelColorValueList[T]) OrderByValue() {
	sort.Slice(*p, func(i, j int) bool {
		return (*p)[i].Value < (*p)[j].Value
	})
}
