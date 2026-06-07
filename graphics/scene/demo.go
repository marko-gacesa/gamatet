// Copyright (c) 2025, 2026 by Marko Gaćeša
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
	"github.com/marko-gacesa/gamatet/logic/anim"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

type Demo struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text

	model       mgl32.Mat4
	demo        core.RenderRequester
	fieldRender *render.Field

	animLast   time.Time
	animPeriod time.Duration
	animFunc   func() anim.Anim

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
		BlockBase: base.NewBlockBaseWithZ(
			renderer,
			tex,
			params.FullW,
			params.FullH,
			max(params.FullW/2, params.FullH/2),
			true,
		),
		res:  *res,
		text: *text,

		model: mgl32.Ident4().
			Mul4(mgl32.Translate3D(params.OffsX, params.OffsY, -0.25)).
			Mul4(mgl32.HomogRotate3DX(params.RotX)).
			Mul4(mgl32.HomogRotate3DY(params.RotY)).
			Mul4(mgl32.HomogRotate3DZ(params.RotZ)),
		demo: params.Demo,

		animLast:   time.Now(),
		animPeriod: params.AnimPeriod,
		animFunc:   params.AnimFunc,

		waitDoneCh: params.Done,
	}

	g.fieldRender = render.NewField(params.Done, g.model, &g.res, &g.text, str, 0, g.demo, render.PreferredSideTopL2R)

	return g
}

func (ft *Demo) Release() {
	<-ft.waitDoneCh

	ft.text.Release()
	ft.res.Release()
}

func (ft *Demo) Prepare(now time.Time) {
	if ft.animFunc != nil && ft.animPeriod > 0 && now.Sub(ft.animLast) >= ft.animPeriod {
		ft.fieldRender.AddAnim(ft.animFunc())
		ft.animLast = now
	}

	ft.fieldRender.Prepare(now)
}

func (ft *Demo) Render() {
	r := ft.Renderer()
	ft.SetCamera()
	ft.fieldRender.Render(r)
}
