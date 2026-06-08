// Copyright (c) 2024, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package rendercache

import (
	"slices"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type modelColorValue[T any] struct {
	Model mgl32.Mat4
	Color mgl32.Vec4
	Value T
}

type modelColorValuePool[T any] struct {
	pool sync.Pool
}

var ModelColorIntPool = newModelColorValuePool[int]()
var ModelColorStringPool = newModelColorValuePool[string]()

func newModelColorValuePool[T any]() modelColorValuePool[T] {
	return modelColorValuePool[T]{
		pool: sync.Pool{
			New: func() any {
				list := make([]modelColorValue[T], 0, 256)
				return (*ModelColorValueList[T])(&list)
			},
		},
	}
}

func (b *modelColorValuePool[T]) Get() *ModelColorValueList[T] {
	list := b.pool.Get().(*ModelColorValueList[T])
	*list = (*list)[:0]
	return list
}

func (b *modelColorValuePool[T]) Put(list *ModelColorValueList[T]) {
	b.pool.Put(list)
}

type ModelColorValueList[T any] []modelColorValue[T]

func (p *ModelColorValueList[T]) Len() int {
	return len(*p)
}

func (p *ModelColorValueList[T]) Add(model mgl32.Mat4, color mgl32.Vec4, value T) {
	*p = append(*p, modelColorValue[T]{
		Model: model,
		Color: color,
		Value: value,
	})
}

func OrderByIntValue(p []modelColorValue[int]) {
	slices.SortFunc(p, func(u, v modelColorValue[int]) int {
		return u.Value - v.Value
	})
}
