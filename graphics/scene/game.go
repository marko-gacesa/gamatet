// Copyright (c) 2025 by Marko Gaćeša

package scene

import (
	"fmt"
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

type Game struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text

	textHUD      render.Text
	fpsHUD       render.HUD
	latenciesHUD render.HUD

	playersInCh  [4]chan<- []byte
	model        mgl32.Mat4
	game         core.RenderRequester
	fieldRenders []*render.Field

	actionCh chan<- action.Action

	waitDoneCh <-chan struct{}
}

var _ screen.Screen = (*Game)(nil)

func NewGame(
	renderer *render.Renderer,
	tex *texture.Manager,
	params types.GameParams,
) *Game {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)

	var latencies fmt.Stringer
	if params.Latencies != nil {
		latencies = params.Latencies
	}

	textHUD := render.MakeText(tex, render.HudFont)
	fpsHUD := render.NewHUD(render.NewFPS(), HUDPosFPS, textHUD)
	latenciesHUD := render.NewHUD(latencies, HUDPosLatencies, textHUD)

	center := mgl32.Ident4()
	fieldModels := make([]mgl32.Mat4, params.FieldCount)

	// assuming all fields have identical dimensions
	w1, h1 := render.GetExtendedContent(params.Game.GetSize(0))

	var w, h int

	switch params.FieldCount {
	case 1:
		w, h = w1, h1
		fieldModels[0] = center
	case 2:
		w, h = 2*w1, h1
		dx := 0.5 * (float32(w) - float32(w1))
		fieldModels[0] = center.Mul4(mgl32.Translate3D(-dx, 0, 0))
		fieldModels[1] = center.Mul4(mgl32.Translate3D(dx, 0, 0))
	case 3:
		w, h = 3*w1+2, h1
		dx := 0.5 * (float32(w) - float32(w1))
		fieldModels[0] = center.Mul4(mgl32.Translate3D(-dx, 0, 0))
		fieldModels[1] = center
		fieldModels[2] = center.Mul4(mgl32.Translate3D(dx, 0, 0))
	case 4:
		w, h = 2*w1+1, 2*h1+1
		dx := 0.5 * (float32(w) - float32(w1))
		dy := 0.5 * (float32(h) - float32(h1))
		fieldModels[0] = center.Mul4(mgl32.Translate3D(-dx, dy, 0))
		fieldModels[1] = center.Mul4(mgl32.Translate3D(dx, dy, 0))
		fieldModels[2] = center.Mul4(mgl32.Translate3D(-dx, -dy, 0))
		fieldModels[3] = center.Mul4(mgl32.Translate3D(dx, -dy, 0))
	case 5:
		w, h = 3*w1+2, 2*h1+1
		dx := 0.5 * (float32(w) - float32(w1))
		dy := 0.5 * (float32(h) - float32(h1))
		fieldModels[0] = center.Mul4(mgl32.Translate3D(-dx, dy, 0))
		fieldModels[1] = center.Mul4(mgl32.Translate3D(0, dy, 0))
		fieldModels[2] = center.Mul4(mgl32.Translate3D(dx, dy, 0))
		dx = 0.5 * (float32(2*w1+1) - float32(w1))
		fieldModels[3] = center.Mul4(mgl32.Translate3D(-dx, -dy, 0))
		fieldModels[4] = center.Mul4(mgl32.Translate3D(dx, -dy, 0))
	case 6:
		w, h = 3*w1+2, 2*h1+1
		dx := 0.5 * (float32(w) - float32(w1))
		dy := 0.5 * (float32(h) - float32(h1))
		fieldModels[0] = center.Mul4(mgl32.Translate3D(-dx, dy, 0))
		fieldModels[1] = center.Mul4(mgl32.Translate3D(0, dy, 0))
		fieldModels[2] = center.Mul4(mgl32.Translate3D(dx, dy, 0))
		fieldModels[3] = center.Mul4(mgl32.Translate3D(-dx, -dy, 0))
		fieldModels[4] = center.Mul4(mgl32.Translate3D(0, -dy, 0))
		fieldModels[5] = center.Mul4(mgl32.Translate3D(dx, -dy, 0))
	case 7, 8:
		panic("TODO")
	default:
		panic("unsupported number of fields")
	}

	fieldRenders := make([]*render.Field, params.FieldCount)
	for i := range params.FieldCount {
		fieldRenders[i] = render.NewField(fieldModels[i], res, text, int(i), params.Game)
	}

	g := &Game{
		BlockBase: base.NewBlockBase(renderer, tex, w, h, true),
		res:       *res,
		text:      *text,

		textHUD:      *textHUD,
		fpsHUD:       *fpsHUD,
		latenciesHUD: *latenciesHUD,

		// these are set below
		playersInCh:  params.PlayerInCh,
		model:        center,
		game:         params.Game,
		fieldRenders: fieldRenders,

		actionCh: params.ActionCh,

		waitDoneCh: params.Done,
	}

	return g
}

func (ft *Game) Release() {
	<-ft.waitDoneCh

	ft.textHUD.Release()

	ft.text.Release()
	ft.res.Release()

	for i := range ft.playersInCh {
		if ft.playersInCh[i] != nil {
			close(ft.playersInCh[i])
		}
	}

	close(ft.actionCh)
}

func (ft *Game) InputKeyPress(key, scancode int) {
	ft.BlockBase.InputKeyPress(key, scancode)

	var a1, a2, a3, a4 action.Action

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

	base.SendAction(a1, ft.waitDoneCh, ft.playersInCh[0])
	base.SendAction(a2, ft.waitDoneCh, ft.playersInCh[1])
	base.SendAction(a3, ft.waitDoneCh, ft.playersInCh[2])
	base.SendAction(a4, ft.waitDoneCh, ft.playersInCh[3])
}

func (ft *Game) Prepare(now time.Time) {
	for i := range ft.fieldRenders {
		ft.fieldRenders[i].Prepare(now)
	}

	ft.fpsHUD.Prepare()
	ft.latenciesHUD.Prepare()
}

func (ft *Game) Render() {
	r := ft.Renderer()

	ft.SetCamera()
	for i := range ft.fieldRenders {
		ft.fieldRenders[i].Render(r)
	}

	ft.fpsHUD.Render(r)
	ft.latenciesHUD.Render(r)
}
