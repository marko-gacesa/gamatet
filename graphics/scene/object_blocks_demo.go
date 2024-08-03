// Copyright (c) 2020-2024 by Marko Gaćeša

package scene

import (
	"context"
	"gamatet/graphics/gtypes"
	"gamatet/graphics/render"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type BlocksDemo struct {
	resources *render.FieldResources
	model     mgl32.Mat4
}

var _ Object = (*BlocksDemo)(nil)

func NewBlocksDemo(resources *render.FieldResources) *BlocksDemo {
	return &BlocksDemo{
		resources: resources,
	}
}

func (b *BlocksDemo) Prepare(ctx context.Context, model *mgl32.Mat4, now time.Time) {
	angle := now.Sub(gtypes.Time).Seconds()
	const scale = 4
	b.model = model.
		Mul4(mgl32.HomogRotate3DZ(float32(angle / 6))).
		Mul4(mgl32.HomogRotate3DY(float32(angle / 2.7))).
		Mul4(mgl32.HomogRotate3DX(float32(angle / 1.2))).
		Mul4(mgl32.Scale3D(scale, scale, scale))
}

func (b *BlocksDemo) Render(r *render.Renderer) {
	var m mgl32.Mat4

	r.Geometry(b.resources.GeomDie)
	r.Material(b.resources.MatWave)
	b.resources.MatWave.Color(mgl32.Vec4{1, 1, 0, 1})
	m = b.model.Mul4(mgl32.Translate3D(-1, -1, 0))
	r.Render(&m)

	r.Geometry(b.resources.GeomRoundedCube)
	r.Material(b.resources.MatRock)
	b.resources.MatRock.Color(mgl32.Vec4{0.5, 1, 0.7, 1})
	m = b.model.Mul4(mgl32.Translate3D(0, -1, 0))
	r.Render(&m)

	r.Geometry(b.resources.GeomSphere)
	r.Material(b.resources.MatColor)
	b.resources.MatColor.Color(mgl32.Vec4{0, 1, 1, 1})
	m = b.model.Mul4(mgl32.Translate3D(1, -1, 0))
	r.Render(&m)

	r.Geometry(b.resources.GeomGem)
	r.Material(b.resources.MatAcid)
	m = b.model.Mul4(mgl32.Translate3D(-1, 0, 0))
	r.Render(&m)

	r.Geometry(b.resources.GeomStar8)
	r.Material(b.resources.MatNorm)
	m = b.model.Mul4(mgl32.Translate3D(0, 0, 0))
	r.Render(&m)

	r.Geometry(b.resources.GeomCube)
	r.Material(b.resources.MatIron)
	m = b.model.Mul4(mgl32.Translate3D(1, 0, 0))
	r.Render(&m)

	r.Geometry(b.resources.GeomFrame)
	r.Material(b.resources.MatRock)
	b.resources.MatRock.ChainTexture(b.resources.TexChain3)
	m = b.model.Mul4(mgl32.Translate3D(-1, 1, 0))
	r.Render(&m)

	r.Geometry(b.resources.GeomRoundedCube)
	r.Material(b.resources.MatLava)
	m = b.model.Mul4(mgl32.Translate3D(0, 1, 0))
	r.Render(&m)

	r.Geometry(b.resources.GeomStar6)
	r.Material(b.resources.MatWave)
	m = b.model.Mul4(mgl32.Translate3D(1, 1, 0))
	r.Render(&m)
}
