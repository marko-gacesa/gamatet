// Copyright (c) 2020-2024 by Marko Gaćeša

package demoblocks

import (
	"context"
	"gamatet/graphics/render"
	"gamatet/graphics/screen"
	"gamatet/graphics/texture"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

const (
	contentW = 40
	contentH = 26
)

type DemoBlocks struct {
	tex  *texture.Manager
	res  render.FieldResources
	text render.Text
	fps  render.FPS

	chReady    chan struct{}
	modelBlock [9]mgl32.Mat4
	modelDigit [contentW * contentH]mgl32.Mat4
}

var _ screen.Screen = (*DemoBlocks)(nil)

var t0 = time.Now()

func NewDemoBlocks(tex *texture.Manager) *DemoBlocks {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)
	fps := render.NewFPS()

	return &DemoBlocks{
		tex:        tex,
		res:        *res,
		text:       *text,
		fps:        *fps,
		chReady:    make(chan struct{}),
		modelBlock: [9]mgl32.Mat4{},
		modelDigit: [contentW * contentH]mgl32.Mat4{},
	}
}

func (d *DemoBlocks) SetCamera(r *render.Renderer, w, h int) {
	r.PerspectiveFull(w, h, contentW, contentH, 12)
}

func (d *DemoBlocks) Release() {
	d.res.Release()
	d.text.Release()
}

func (d *DemoBlocks) Shutdown() {}

func (d *DemoBlocks) KeyAction(key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey) {

}

func (d *DemoBlocks) Prepare(ctx context.Context, now time.Time) {
	go func() {
		defer func() { d.chReady <- struct{}{} }()

		modelCenter := mgl32.Ident4()

		angle := now.Sub(t0).Seconds()
		const scale = 4
		modelSpin := modelCenter.
			Mul4(mgl32.HomogRotate3DZ(float32(angle / 6))).
			Mul4(mgl32.HomogRotate3DY(float32(angle / 2.7))).
			Mul4(mgl32.HomogRotate3DX(float32(angle / 1.2))).
			Mul4(mgl32.Scale3D(scale, scale, scale))

		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				idx := 3*(y+1) + x + 1
				d.modelBlock[idx] = modelSpin.Mul4(mgl32.Translate3D(float32(x), float32(y), 0))
			}
		}

		for i := 0; i < contentW; i++ {
			for j := 0; j < contentH; j++ {
				modelDigit := modelCenter.Mul4(mgl32.Translate3D(-contentW/2+0.5+float32(i), -contentH/2+0.5+float32(j), 0))
				d.modelDigit[j*contentW+i] = modelDigit
			}
		}
	}()
}

func (d *DemoBlocks) Render(r *render.Renderer) {
	<-d.chReady

	center := mgl32.Ident4()

	for i := 0; i < contentW; i++ {
		for j := 0; j < contentH; j++ {
			d.text.Rune(r, d.modelDigit[j*contentW+i], mgl32.Vec4{0.1, 0.1, 0, 1}, '0'+rune(i+j)%10)
		}
	}

	r.Geometry(d.res.GeomDie)
	r.Material(d.res.MatWave)
	d.res.MatWave.Color(mgl32.Vec4{1, 1, 0, 1})
	r.Render(&d.modelBlock[0])

	r.Geometry(d.res.GeomStar6)
	r.Material(d.res.MatTexUV)
	r.Render(&d.modelBlock[1])

	r.Geometry(d.res.GeomSphere)
	r.Material(d.res.MatColor)
	d.res.MatColor.Color(mgl32.Vec4{0, 1, 1, 1})
	r.Render(&d.modelBlock[2])

	r.Geometry(d.res.GeomGem)
	r.Material(d.res.MatAcid)
	r.Render(&d.modelBlock[3])

	r.Geometry(d.res.GeomStar8)
	r.Material(d.res.MatNorm)
	r.Render(&d.modelBlock[4])

	r.Geometry(d.res.GeomCube)
	r.Material(d.res.MatIron)
	r.Render(&d.modelBlock[5])

	r.Geometry(d.res.GeomFrame)
	r.Material(d.res.MatRock)
	d.res.MatRock.ChainTexture(d.res.TexChain3)
	r.Render(&d.modelBlock[6])

	r.Geometry(d.res.GeomRoundedCube)
	r.Material(d.res.MatLava)
	r.Render(&d.modelBlock[7])

	r.Geometry(d.res.GeomRoundedCube)
	d.text.Material(r, mgl32.Vec4{0.5, 1, 0.7, 1}, 'M')
	r.Render(&d.modelBlock[8])

	d.fps.Render(r, &d.text, center.Mul4(mgl32.Translate3D(-contentW/2+0.5, -contentH/2+0.5, 0)))
}
