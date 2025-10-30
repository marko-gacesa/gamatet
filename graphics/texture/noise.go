// Copyright (c) 2020-2025 by Marko Gaćeša

package texture

import (
	"gamatet/graphics/gutil"
	"math"
	"math/rand"
)

func Perlin2D(nDim, mDim int, seed int64) []float32 {
	if nDim <= 0 || !gutil.IsPow2(nDim) {
		panic("nDim must be a power of 2")
	}

	if mDim <= 0 || !gutil.IsPow2(mDim) {
		panic("mDim must be a power of 2")
	}

	if nDim <= mDim {
		panic("mDim must be smaller than nDim")
	}

	size := nDim * nDim
	values := make([]float32, size)

	mesh := newMesh2D(mDim, seed)
	iterations := gutil.Log2(nDim / mDim)

	for range iterations {
		mDim = mesh.dim

		cellSize := nDim / mDim
		cellSizeF := float32(cellSize)

		for yCell := 0; yCell < mDim; yCell++ {
			for xCell := 0; xCell < mDim; xCell++ {
				mesh.SetMeshCell(xCell, yCell)
				for y := range cellSize {
					for x := range cellSize {
						v := mesh.Interpolate(float32(x)/cellSizeF, float32(y)/cellSizeF)
						idx := (yCell*cellSize+y)*nDim + (xCell*cellSize + x)
						values[idx] += v
					}
				}
			}
		}

		mesh = mesh.Double()
	}

	return values
}

type mesh2D struct {
	dim    int
	amp    float32
	random *rand.Rand
	values []float32
	interX [4]gutil.CatmullRom
}

func newMesh2D(dim int, seed int64) *mesh2D {
	random := rand.New(rand.NewSource(seed))
	m := &mesh2D{
		dim:    dim,
		amp:    1.0,
		random: random,
		values: make([]float32, dim*dim),
		interX: [4]gutil.CatmullRom{},
	}
	n := dim * dim
	for i := range n {
		m.values[i] = random.Float32()
	}
	return m
}

func (m *mesh2D) Double() *mesh2D {
	const ampDiv = 2.0
	dim := m.dim * 2
	amp := m.amp / ampDiv
	values := make([]float32, dim*dim)
	for y := range dim {
		for x := range dim {
			divX := x / 2
			modX := x % 2
			divY := y / 2
			modY := y % 2
			if modX == 0 && modY == 0 {
				values[y*dim+x] = m.values[divY*m.dim+divX] / ampDiv
			} else {
				values[y*dim+x] = m.random.Float32() * amp
			}
		}
	}
	return &mesh2D{dim: dim, amp: amp, random: m.random, values: values}
}

func (m *mesh2D) getValue(x, y int) float32 {
	dim := m.dim
	if x >= 0 {
		x = x % dim
	} else {
		x = (dim - 1) + (x+1)%dim
	}
	if y >= 0 {
		y = y % dim
	} else {
		y = (dim - 1) + (y+1)%dim
	}

	return m.values[y*m.dim+x]
}

func (m *mesh2D) SetMeshCell(x, y int) {
	for i := range 4 {
		m.interX[i].Update(m.getValue(x-1, y-1+i), m.getValue(x, y-1+i), m.getValue(x+1, y-1+i), m.getValue(x+2, y-1+i))
	}
}

func (m *mesh2D) Interpolate(x, y float32) float32 {
	var interY gutil.CatmullRom
	interY.Update(m.interX[0].Value(x), m.interX[1].Value(x), m.interX[2].Value(x), m.interX[3].Value(x))
	return interY.Value(y)
}

func clamp(values []float32, min, max float32) {
	size := len(values)
	var minV float32
	var maxV float32
	minV = math.MaxFloat32
	maxV = -math.MaxFloat32
	for i := range size {
		if values[i] < minV {
			minV = values[i]
		}
		if values[i] > maxV {
			maxV = values[i]
		}
	}

	//fmt.Println("minV, maxV", minV, maxV)

	dV := maxV - minV
	v := max - min
	for i := range size {
		values[i] = ((values[i]-minV)/dV)*v + min
	}
}

func clampSin(values []float32, min, max float32) {
	size := len(values)
	var minV float32
	var maxV float32
	minV = math.MaxFloat32
	maxV = -math.MaxFloat32
	for i := range size {
		if values[i] < minV {
			minV = values[i]
		}
		if values[i] > maxV {
			maxV = values[i]
		}
	}

	//fmt.Println("minV, maxV", minV, maxV)

	dV := maxV - minV
	v := max - min
	for i := range size {
		values[i] = float32(math.Sin(float64((values[i]-minV)/dV)*2*math.Pi)*0.5+0.5)*v + min
	}
}

func symXY(values []float32, size int) {
	for y := 0; y < size/2; y++ {
		k := float32(2*y) / float32(size)
		k = 1 - k*k*(3-2*k)
		for x := y; x < size-y-1; x++ {
			x1, y1 := x, y
			x2, y2 := size-y-1, x
			x3, y3 := size-x-1, size-y-1
			x4, y4 := y, size-x-1

			idx1 := y1*size + x1
			idx2 := y2*size + x2
			idx3 := y3*size + x3
			idx4 := y4*size + x4

			v1 := values[idx1]
			v2 := values[idx2]
			v3 := values[idx3]
			v4 := values[idx4]

			v := (v1 + v2 + v3 + v4) / 4

			values[idx1] = k*v + (1-k)*values[idx1]
			values[idx2] = k*v + (1-k)*values[idx2]
			values[idx3] = k*v + (1-k)*values[idx3]
			values[idx4] = k*v + (1-k)*values[idx4]
		}
	}
}
