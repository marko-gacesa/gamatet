// Copyright (c) 2024, 2025 by Marko Gaćeša

package scene

import (
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/graphics/render"
	"gamatet/graphics/scene/base"
	"gamatet/graphics/texture"
	"gamatet/internal/types"
	"gamatet/logic/screen"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type GameOne struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text

	textHUD render.Text
	fpsHUD  render.HUD

	playerInCh  chan<- []byte
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

	preferredSide := render.PreferredSideTop

	wf, hf, pf := params.Game.GetSize(0)
	w, h := render.GetExtendedContent(wf, hf, preferredSide.PosN(pf))

	g := &GameOne{
		BlockBase: base.NewBlockBase(renderer, tex, w, h, true),
		res:       *res,
		text:      *text,

		textHUD: *textHUD,
		fpsHUD:  *fpsHUD,

		// these are set below
		playerInCh:  nil,
		model:       mgl32.Mat4{},
		game:        nil,
		fieldRender: nil,

		actionCh: params.ActionCh,

		waitDoneCh: params.Done,
	}

	g.playerInCh = params.PlayerInCh
	g.model = mgl32.Ident4()
	g.game = params.Game
	g.fieldRender = render.NewField(g.model, &g.res, &g.text, 0, g.game, preferredSide)

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

	var a action.Action

	switch glfw.Key(key) {
	case glfw.KeyEscape:
		ft.actionCh <- action.Abort
	case glfw.KeyPause, glfw.KeyP:
		ft.actionCh <- action.Pause

	case glfw.KeyLeft:
		a = action.MoveLeft
	case glfw.KeyRight:
		a = action.MoveRight
	case glfw.KeyUp:
		a = action.Activate
	case glfw.KeyDown:
		a = action.MoveDown
	case glfw.KeySpace:
		a = action.Drop
	}

	base.SendAction(a, ft.waitDoneCh, ft.playerInCh)
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
