// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package scene

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/graphics/render"
	"github.com/marko-gacesa/gamatet/graphics/scene/base"
	"github.com/marko-gacesa/gamatet/graphics/texture"
	"github.com/marko-gacesa/gamatet/internal/types"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

type Demo struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text

	model       mgl32.Mat4
	demo        core.RenderRequester
	fieldRender *render.Field

	waitDoneCh <-chan struct{}
}

var _ screen.Screen = (*Demo)(nil)

func NewDemo(
	renderer *render.Renderer,
	tex *texture.Manager,
	params types.DemoParams,
) *Demo {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)

	str := fieldStrings()

	g := &Demo{
		BlockBase: base.NewBlockBase(renderer, tex, params.FullW, params.FullH, true),
		res:       *res,
		text:      *text,

		model: mgl32.Ident4().
			Mul4(mgl32.Translate3D(params.OffsX, params.OffsY, -0.25)).
			Mul4(mgl32.HomogRotate3DX(params.RotX)).
			Mul4(mgl32.HomogRotate3DY(params.RotY)).
			Mul4(mgl32.HomogRotate3DZ(params.RotZ)),
		demo: params.Demo,

		waitDoneCh: params.Done,
	}

	g.fieldRender = render.NewField(g.model, &g.res, &g.text, str, 0, g.demo, render.PreferredSideTopL2R)

	return g
}

func (ft *Demo) Release() {
	<-ft.waitDoneCh

	ft.text.Release()
	ft.res.Release()
}

func (ft *Demo) InputKeyPress(key int, act screen.KeyAction) {}

func (ft *Demo) Prepare(now time.Time) {
	ft.fieldRender.Prepare(now)
}

func (ft *Demo) Render() {
	r := ft.Renderer()
	ft.SetCamera()
	ft.fieldRender.Render(r)
}
