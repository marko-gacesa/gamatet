// Copyright (c) 2024 by Marko Gaćeša

package rendercache

import (
	"github.com/go-gl/mathgl/mgl32"
	"sort"
	"sync"
)

type modelColorValue struct {
	Model mgl32.Mat4
	Color mgl32.Vec4
	Value int
}

type modelColorValuePool struct {
	pool sync.Pool
}

var ModelColorValuePool = modelColorValuePool{
	pool: sync.Pool{
		New: func() any {
			list := make([]modelColorValue, 0, 256)
			return ModelColorValueList(list)
		},
	},
}

func (b *modelColorValuePool) Get() ModelColorValueList {
	list := b.pool.Get().(ModelColorValueList)
	list = list[:0]
	return list
}

func (b *modelColorValuePool) Put(list ModelColorValueList) {
	b.pool.Put(list)
}

type ModelColorValueList []modelColorValue

func (p *ModelColorValueList) Add(model mgl32.Mat4, color mgl32.Vec4, value int) {
	*p = append(*p, modelColorValue{
		Model: model,
		Color: color,
		Value: value,
	})
}

func (p *ModelColorValueList) OrderByValue() {
	sort.Slice(*p, func(i, j int) bool {
		return (*p)[i].Value < (*p)[j].Value
	})
}
