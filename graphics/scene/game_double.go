// Copyright (c) 2024, 2025 by Marko Gaćeša

package scene

import (
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/graphics/render"
	"gamatet/graphics/scene/base"
	"gamatet/graphics/texture"
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

	waitDoneCh <-chan struct{}
}

var _ screen.Screen = (*GameDouble)(nil)

func NewGameDouble(
	renderer *render.Renderer,
	tex *texture.Manager,
	params core.GameDoubleParams,
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
}

func (ft *GameDouble) InputKeyPress(key, scancode int) {
	var cmd1, cmd2 []byte

	switch glfw.Key(key) {
	case glfw.KeyEscape:
		cmd1 = []byte{byte(action.Abort)}
	case glfw.KeyPause:
		cmd1 = []byte{byte(action.Pause)}

	case glfw.KeyLeft:
		cmd1 = []byte{byte(action.MoveLeft)}
	case glfw.KeyRight:
		cmd1 = []byte{byte(action.MoveRight)}
	case glfw.KeyUp:
		cmd1 = []byte{byte(action.Activate)}
	case glfw.KeyDown:
		cmd1 = []byte{byte(action.MoveDown)}
	case glfw.KeyRightShift:
		cmd1 = []byte{byte(action.Drop)}

	case glfw.KeyA:
		cmd2 = []byte{byte(action.MoveLeft)}
	case glfw.KeyD:
		cmd2 = []byte{byte(action.MoveRight)}
	case glfw.KeyW:
		cmd2 = []byte{byte(action.Activate)}
	case glfw.KeyS:
		cmd2 = []byte{byte(action.MoveDown)}
	case glfw.KeyLeftShift:
		cmd2 = []byte{byte(action.Drop)}
	}

	base.SendAction(cmd1, ft.waitDoneCh, ft.player1InCh)
	base.SendAction(cmd2, ft.waitDoneCh, ft.player2InCh)
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
