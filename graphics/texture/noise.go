// Copyright (c) 2020 by Marko Gaćeša

package texture

import (
	"gamatet/graphics/gutil"
	"math"
	"math/rand"
)

func Clamp(values []float32, min, max float32) {
	size := len(values)
	var minV float32
	var maxV float32
	minV = math.MaxFloat32
	maxV = -math.MaxFloat32
	for i := 0; i < size; i++ {
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
	for i := 0; i < size; i++ {
		values[i] = ((values[i]-minV)/dV)*v + min
	}
}

func ClampSin(values []float32, min, max float32) {
	size := len(values)
	var minV float32
	var maxV float32
	minV = math.MaxFloat32
	maxV = -math.MaxFloat32
	for i := 0; i < size; i++ {
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
	for i := 0; i < size; i++ {
		values[i] = float32(math.Sin(float64((values[i]-minV)/dV)*2*math.Pi)*0.5+0.5)*v + min
	}
}

func Perlin2D(nDim, mDimIter1, iterations int, seed int64) []float32 {
	size := nDim * nDim

	values := make([]float32, size)

	mDimMax := mDimIter1 << (iterations - 1)
	if nDim < mDimMax || nDim%mDimMax != 0 {
		panic("invalid parameter combination")
	}

	mesh := NewMesh2D(mDimIter1, seed)
	for i := 0; i < iterations; i++ {
		mDim := mesh.dim

		cellSize := nDim / mDim
		cellSizeF := float32(cellSize)

		for yCell := 0; yCell < mDim; yCell++ {
			for xCell := 0; xCell < mDim; xCell++ {
				mesh.SetMeshCell(xCell, yCell)
				for y := 0; y < cellSize; y++ {
					for x := 0; x < cellSize; x++ {
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

type Mesh2D struct {
	dim    int
	amp    float32
	random *rand.Rand
	values []float32
	interX [4]gutil.CatmullRom
}

func NewMesh2D(dim int, seed int64) *Mesh2D {
	random := rand.New(rand.NewSource(seed))
	m := &Mesh2D{
		dim:    dim,
		amp:    1.0,
		random: random,
		values: make([]float32, dim*dim),
		interX: [4]gutil.CatmullRom{},
	}
	n := dim * dim
	for i := 0; i < n; i++ {
		m.values[i] = random.Float32()
	}
	return m
}

func (m *Mesh2D) Double() *Mesh2D {
	const ampDiv = 2.0
	dim := m.dim * 2
	amp := m.amp / ampDiv
	values := make([]float32, dim*dim)
	for y := 0; y < dim; y++ {
		for x := 0; x < dim; x++ {
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
	return &Mesh2D{dim: dim, amp: amp, random: m.random, values: values}
}

func (m *Mesh2D) getValue(x, y int) float32 {
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

func (m *Mesh2D) SetMeshCell(x, y int) {
	for i := 0; i < 4; i++ {
		m.interX[i].Update(m.getValue(x-1, y-1+i), m.getValue(x, y-1+i), m.getValue(x+1, y-1+i), m.getValue(x+2, y-1+i))
	}
}

func (m *Mesh2D) Interpolate(x, y float32) float32 {
	var interY gutil.CatmullRom
	interY.Update(m.interX[0].Value(x), m.interX[1].Value(x), m.interX[2].Value(x), m.interX[3].Value(x))
	return interY.Value(y)
}
