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

type GameDouble struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text

	textHUD render.Text
	fpsHUD  render.FPS

	player1InCh chan<- []byte
	player2InCh chan<- []byte
	model       mgl32.Mat4
	game        core.RenderRequester
	fieldRender *render.Field

	actionCh chan<- action.Action

	waitDoneCh <-chan struct{}
}

var _ screen.Screen = (*GameDouble)(nil)

func NewGameDouble(
	renderer *render.Renderer,
	tex *texture.Manager,
	params types.GameDoubleParams,
) *GameDouble {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)

	textHUD := render.MakeText(tex, render.HudFont)
	fpsHUD := render.NewFPS()

	w, h := render.GetExtendedContent(params.Game.GetSize(0))

	g := &GameDouble{
		BlockBase: base.NewBlockBase(renderer, tex, w, h, true),
		res:       *res,
		text:      *text,

		textHUD: *textHUD,
		fpsHUD:  *fpsHUD,

		// these are set below
		player1InCh: nil,
		player2InCh: nil,
		model:       mgl32.Mat4{},
		game:        nil,
		fieldRender: nil,

		actionCh: params.ActionCh,

		waitDoneCh: params.Done,
	}

	g.player1InCh = params.Player1InCh
	g.player2InCh = params.Player2InCh
	g.model = mgl32.Ident4()
	g.game = params.Game
	g.fieldRender = render.NewField(g.model, &g.res, &g.text, 0, g.game)

	return g
}

func (ft *GameDouble) Release() {
	<-ft.waitDoneCh
	ft.text.Release()
	ft.res.Release()

	close(ft.player1InCh)
	close(ft.player2InCh)

	close(ft.actionCh)
}

func (ft *GameDouble) InputKeyPress(key, scancode int) {
	var a1, a2 action.Action

	switch glfw.Key(key) {
	case glfw.KeyEscape:
		ft.actionCh <- action.Abort
	case glfw.KeyPause:
		ft.actionCh <- action.Pause

	case glfw.KeyLeft:
		a1 = action.MoveLeft
	case glfw.KeyRight:
		a1 = action.MoveRight
	case glfw.KeyUp:
		a1 = action.Activate
	case glfw.KeyDown:
		a1 = action.MoveDown
	case glfw.KeyRightShift:
		a1 = action.Drop

	case glfw.KeyA:
		a2 = action.MoveLeft
	case glfw.KeyD:
		a2 = action.MoveRight
	case glfw.KeyW:
		a2 = action.Activate
	case glfw.KeyS:
		a2 = action.MoveDown
	case glfw.KeyLeftShift:
		a2 = action.Drop
	}

	base.SendAction(a1, ft.waitDoneCh, ft.player1InCh)
	base.SendAction(a2, ft.waitDoneCh, ft.player2InCh)
}

func (ft *GameDouble) Prepare(now time.Time) {
	ft.fieldRender.Prepare(now)
}

func (ft *GameDouble) Render() {
	r := ft.Renderer()

	ft.SetCamera()
	ft.fieldRender.Render(r)

	ft.fpsHUD.Render(r, &ft.textHUD)
}
