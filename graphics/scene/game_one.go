// Copyright (c) 2024-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package scene

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/graphics/render"
	"github.com/marko-gacesa/gamatet/graphics/scene/base"
	"github.com/marko-gacesa/gamatet/graphics/scene/hud"
	"github.com/marko-gacesa/gamatet/graphics/texture"
	"github.com/marko-gacesa/gamatet/internal/config"
	"github.com/marko-gacesa/gamatet/internal/types"
	"github.com/marko-gacesa/gamatet/logic/gamepad"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

type GameOne struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text
	huds *hud.HUDs

	playerInCh chan<- []byte
	input      config.Input

	model       mgl32.Mat4
	game        core.RenderRequester
	fieldRender *render.Field

	actionCh chan<- action.Action

	waitDoneCh <-chan struct{}
}

var _ screen.Screen = (*GameOne)(nil)

func NewGameOne(
	renderer *render.Renderer,
	tex *texture.Manager,
	params types.GameOneParams,
) *GameOne {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)
	huds := hud.NewHUDs(tex)
	huds.Add(render.NewFPS(), hud.PosFPS)

	str := fieldStrings()

	preferredSide := render.PreferredSideTopL2R

	wf, hf, pf := params.Game.GetSize(0)
	w, h := render.GetExtendedContent(wf, hf, preferredSide.PieceCorners(pf))

	g := &GameOne{
		BlockBase: base.NewBlockBase(renderer, tex, w, h, true),
		res:       *res,
		text:      *text,
		huds:      huds,

		playerInCh: params.PlayerInCh,
		input:      params.PlayerInput,

		model: mgl32.Ident4().Mul4(mgl32.Translate3D(0, 0, -0.25)), // -0.5 for ideal 2D/3D transition
		game:  params.Game,

		actionCh: params.ActionCh,

		waitDoneCh: params.Done,
	}

	g.fieldRender = render.NewField(g.model, &g.res, &g.text, str, 0, g.game, preferredSide)

	return g
}

func (ft *GameOne) Release() {
	<-ft.waitDoneCh

	ft.huds.Release()
	ft.text.Release()
	ft.res.Release()

	close(ft.playerInCh)

	close(ft.actionCh)
}

func (ft *GameOne) InputKeyPress(key int, act screen.KeyAction) {
	ft.BlockBase.InputKeyPress(key, act)
	ft.huds.InputKeyPress(key, act)

	inputKeyboardGlobal(key, act, ft.actionCh)

	if ft.input.Source == config.InputSourceKeyboard {
		inputKeyboardPlayer(key, act, ft.input.Keys, ft.waitDoneCh, ft.playerInCh)
	}
}

func (ft *GameOne) InputGamepadPress(gamepadIdx int, b gamepad.ButtonChange) {
	if ft.input.Source == config.InputSourceGamepad && ft.input.Gamepad == gamepadIdx {
		inputGamepad(b, ft.actionCh, ft.waitDoneCh, ft.playerInCh)
	}
}

func (ft *GameOne) Prepare(now time.Time) {
	ft.fieldRender.Prepare(now)
	ft.huds.Prepare(ft.ViewSize())
}

func (ft *GameOne) Render() {
	r := ft.Renderer()

	ft.SetCamera()
	ft.fieldRender.Render(r)

	ft.huds.Render(r)
}
