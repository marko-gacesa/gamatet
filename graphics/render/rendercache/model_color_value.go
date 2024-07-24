// Copyright (c) 2024 by Marko Gaćeša

package rendercache

import (
	"gamatet/graphics/gtypes"
	"github.com/go-gl/mathgl/mgl32"
	"sort"
	"sync"
)

type modelColorValuePool struct {
	pool sync.Pool
}

var ModelColorValuePool = modelColorValuePool{
	pool: sync.Pool{
		New: func() any {
			list := make([]gtypes.ModelColorValue, 0, 256)
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

type ModelColorValueList []gtypes.ModelColorValue

func (p *ModelColorValueList) Add(model mgl32.Mat4, color mgl32.Vec4, value int) {
	*p = append(*p, gtypes.ModelColorValue{
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
