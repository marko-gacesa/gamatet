// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package scene

import (
	"time"

	"github.com/marko-gacesa/gamatet/internal/config/key"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/graphics/render"
	"github.com/marko-gacesa/gamatet/graphics/scene/base"
	"github.com/marko-gacesa/gamatet/graphics/texture"
	"github.com/marko-gacesa/gamatet/internal/types"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

type GameOne struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text

	textHUD render.Text
	fpsHUD  render.HUD

	playerInCh chan<- []byte
	input      key.Input

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

	textHUD := render.MakeText(tex, render.HudFont)
	fpsHUD := render.NewHUD(render.NewFPS(), HUDPosFPS, textHUD)

	str := fieldStrings()

	preferredSide := render.PreferredSideTopL2R

	wf, hf, pf := params.Game.GetSize(0)
	w, h := render.GetExtendedContent(wf, hf, preferredSide.PieceCorners(pf))

	g := &GameOne{
		BlockBase: base.NewBlockBase(renderer, tex, w, h, true),
		res:       *res,
		text:      *text,

		textHUD: *textHUD,
		fpsHUD:  *fpsHUD,

		playerInCh: params.PlayerInCh,
		input:      params.PlayerInput,

		model:       mgl32.Ident4(),
		game:        params.Game,
		fieldRender: nil,

		actionCh: params.ActionCh,

		waitDoneCh: params.Done,
	}

	g.fieldRender = render.NewField(g.model, &g.res, &g.text, str, 0, g.game, preferredSide)

	return g
}

func (ft *GameOne) Release() {
	<-ft.waitDoneCh
	ft.text.Release()
	ft.res.Release()

	close(ft.playerInCh)

	close(ft.actionCh)
}

func (ft *GameOne) InputKeyPress(key, scancode int) {
	ft.BlockBase.InputKeyPress(key, scancode)

	switch glfw.Key(key) {
	case glfw.KeyEscape:
		ft.actionCh <- action.Abort
	case glfw.KeyPause:
		ft.actionCh <- action.Pause

	case KeyMap[ft.input.Left]:
		base.SendAction(action.MoveLeft, ft.waitDoneCh, ft.playerInCh)
	case KeyMap[ft.input.Right]:
		base.SendAction(action.MoveRight, ft.waitDoneCh, ft.playerInCh)
	case KeyMap[ft.input.Activate]:
		base.SendAction(action.Activate, ft.waitDoneCh, ft.playerInCh)
	case KeyMap[ft.input.Drop]:
		base.SendAction(action.Drop, ft.waitDoneCh, ft.playerInCh)
	}
}

func (ft *GameOne) Prepare(now time.Time) {
	ft.fieldRender.Prepare(now)
	ft.fpsHUD.Prepare()
}

func (ft *GameOne) Render() {
	r := ft.Renderer()

	ft.SetCamera()
	ft.fieldRender.Render(r)

	ft.fpsHUD.Render(r)
}
