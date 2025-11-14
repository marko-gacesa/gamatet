// Copyright (c) 2025 by Marko Gaćeša

package scene

import (
	"fmt"
	"time"

	"github.com/marko-gacesa/gamatet/game/setup"
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

type Game struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text

	textHUD      render.Text
	fpsHUD       render.HUD
	latenciesHUD render.HUD

	playersInCh [setup.MaxLocalPlayers]chan<- []byte
	inputs      [setup.MaxLocalPlayers]key.Input

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

	fieldRenders := make([]*render.Field, 0, params.FieldCount)
	var w, h int

	switch params.FieldCount {
	case 1:
		wf, hf, pf := params.Game.GetSize(0)
		w, h = render.GetExtendedContent(wf, hf, render.PreferredSideTop.PosN(pf))
		fieldRenders = append(fieldRenders, render.NewField(center, res, text, 0, params.Game, render.PreferredSideTop))
	case 2:
		wf0, hf0, pf0 := params.Game.GetSize(0)
		wf1, hf1, pf1 := params.Game.GetSize(1)
		w0, h0 := render.GetExtendedContent(wf0, hf0, render.PreferredSideLeft.PosN(pf0))
		w1, h1 := render.GetExtendedContent(wf1, hf1, render.PreferredSideRight.PosN(pf1))

		w, h = w0+w1, max(h0, h1)
		dx := 0.5 * (float32(w) - float32(w1))
		fieldModel0 := center.Mul4(mgl32.Translate3D(-dx, 0, 0))
		fieldModel1 := center.Mul4(mgl32.Translate3D(dx, 0, 0))
		fieldRenders = append(fieldRenders, render.NewField(fieldModel0, res, text, 0, params.Game, render.PreferredSideLeft))
		fieldRenders = append(fieldRenders, render.NewField(fieldModel1, res, text, 1, params.Game, render.PreferredSideRight))
		/*
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
		*/
	default:
		panic("unsupported number of fields")
	}

	g := &Game{
		BlockBase: base.NewBlockBase(renderer, tex, w, h, true),
		res:       *res,
		text:      *text,

		textHUD:      *textHUD,
		fpsHUD:       *fpsHUD,
		latenciesHUD: *latenciesHUD,

		playersInCh: params.PlayerInCh,
		inputs:      params.PlayerInputs,

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

	k := glfw.Key(key)

	switch k {
	case glfw.KeyEscape:
		ft.actionCh <- action.Abort
	case glfw.KeyPause:
		ft.actionCh <- action.Pause
	}

	for i := range setup.MaxLocalPlayers {
		switch k {
		case KeyMap[ft.inputs[i].Left]:
			base.SendAction(action.MoveLeft, ft.waitDoneCh, ft.playersInCh[i])
		case KeyMap[ft.inputs[i].Right]:
			base.SendAction(action.MoveRight, ft.waitDoneCh, ft.playersInCh[i])
		case KeyMap[ft.inputs[i].Activate]:
			base.SendAction(action.Activate, ft.waitDoneCh, ft.playersInCh[i])
		case KeyMap[ft.inputs[i].Drop]:
			base.SendAction(action.Drop, ft.waitDoneCh, ft.playersInCh[i])
		}
	}
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
