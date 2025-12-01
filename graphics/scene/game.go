// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package scene

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/graphics/render"
	"github.com/marko-gacesa/gamatet/graphics/scene/base"
	"github.com/marko-gacesa/gamatet/graphics/texture"
	"github.com/marko-gacesa/gamatet/internal/config/key"
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

	str := fieldStrings()

	center := mgl32.Ident4()

	var fieldRenders []*render.Field
	var w, h int

	switch params.FieldCount {
	case 1:
		wf, hf, pf := params.Game.GetSize(0)
		w, h = render.GetExtendedContent(wf, hf, render.PreferredSideTopL2R.PieceCorners(pf))
		fieldRenders = []*render.Field{
			render.NewField(center, res, text, str, 0, params.Game, render.PreferredSideTopL2R),
		}
	case 2:
		psides := []render.PreferredSide{
			render.PreferredSideLeftT2B, render.PreferredSideRightT2B,
		}
		w, h, fieldRenders = positionFieldRenderers(params.Game, res, text, str, psides, 2)
	case 3:
		psides := []render.PreferredSide{
			render.PreferredSideLeftT2B, render.PreferredSideLeftB2T, render.PreferredSideLeftT2B,
		}
		w, h, fieldRenders = positionFieldRenderers(params.Game, res, text, str, psides, 3)
	case 4:
		w0, h0, _ := params.Game.GetSize(0)
		if h0 >= 2*w0 {
			psides := []render.PreferredSide{
				render.PreferredSideTopL2R, render.PreferredSideTopL2R, render.PreferredSideTopL2R, render.PreferredSideTopL2R,
			}
			w, h, fieldRenders = positionFieldRenderers(params.Game, res, text, str, psides, 4)
		} else {
			psides := []render.PreferredSide{
				render.PreferredSideTopL2R, render.PreferredSideTopR2L,
				render.PreferredSideBottomL2R, render.PreferredSideBottomR2L,
			}
			w, h, fieldRenders = positionFieldRenderers(params.Game, res, text, str, psides, 2)
		}
	case 5:
		psides := []render.PreferredSide{
			render.PreferredSideTopL2R, render.PreferredSideTopL2R, render.PreferredSideTopL2R,
			render.PreferredSideBottomR2L, render.PreferredSideBottomR2L,
		}
		w, h, fieldRenders = positionFieldRenderers(params.Game, res, text, str, psides, 3)
	case 6:
		psides := []render.PreferredSide{
			render.PreferredSideTopL2R, render.PreferredSideTopL2R, render.PreferredSideTopL2R,
			render.PreferredSideBottomL2R, render.PreferredSideBottomL2R, render.PreferredSideBottomL2R,
		}
		w, h, fieldRenders = positionFieldRenderers(params.Game, res, text, str, psides, 3)
	case 7:
		psides := []render.PreferredSide{
			render.PreferredSideTopL2R, render.PreferredSideTopL2R, render.PreferredSideTopL2R, render.PreferredSideTopL2R,
			render.PreferredSideBottomL2R, render.PreferredSideBottomL2R, render.PreferredSideBottomL2R,
		}
		w, h, fieldRenders = positionFieldRenderers(params.Game, res, text, str, psides, 4)
	case 8:
		psides := []render.PreferredSide{
			render.PreferredSideTopL2R, render.PreferredSideTopL2R, render.PreferredSideTopL2R, render.PreferredSideTopL2R,
			render.PreferredSideBottomL2R, render.PreferredSideBottomL2R, render.PreferredSideBottomL2R, render.PreferredSideBottomL2R,
		}
		w, h, fieldRenders = positionFieldRenderers(params.Game, res, text, str, psides, 4)
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

func (ft *Game) InputKeyPress(key int, act screen.KeyAction) {
	ft.BlockBase.InputKeyPress(key, act)

	k := glfw.Key(key)

	if act == screen.KeyActionPress {
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
			case KeyMap[ft.inputs[i].Boost]:
				base.SendAction(action.SpeedUp, ft.waitDoneCh, ft.playersInCh[i])
			case KeyMap[ft.inputs[i].Drop]:
				base.SendAction(action.Drop, ft.waitDoneCh, ft.playersInCh[i])
			}
		}
	} else if act == screen.KeyActionRelease {
		for i := range setup.MaxLocalPlayers {
			switch k {
			case KeyMap[ft.inputs[i].Boost]:
				base.SendAction(action.SpeedDown, ft.waitDoneCh, ft.playersInCh[i])
			}
		}
	}
}

func (ft *Game) Prepare(now time.Time) {
	for i := range ft.fieldRenders {
		ft.fieldRenders[i].Prepare(now)
	}

	ft.fpsHUD.Prepare(ft.ViewSize())
	ft.latenciesHUD.Prepare(ft.ViewSize())
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

func positionFieldRenderers(
	rr core.RenderRequester,
	res *render.FieldResources,
	text *render.Text,
	str render.FieldStrings,
	psides []render.PreferredSide,
	fieldsInRow int,
) (int, int, []*render.Field) {
	n := len(psides)

	fieldSizes := make([]struct{ w, h int }, n)
	for idx := range n {
		fw, fh, fpc := rr.GetSize(idx)
		w, h := render.GetExtendedContent(fw, fh, psides[idx].PieceCorners(fpc))
		fieldSizes[idx] = struct{ w, h int }{w: w, h: h}
	}

	var gridW, gridH int
	var gridWCurr, gridHCurr, k int

	gridWCurr, gridHCurr, k = 0, 0, 0
	for idx := range n {
		gridWCurr += fieldSizes[idx].w
		gridHCurr = max(fieldSizes[idx].h, gridHCurr)
		k++
		if k >= fieldsInRow {
			gridW = max(gridWCurr, gridW)
			gridH += gridHCurr
			gridWCurr, k = 0, 0
		}
	}

	fieldsInColumn := len(fieldSizes) / fieldsInRow
	if len(fieldSizes)%fieldsInRow > 0 {
		fieldsInColumn++
	}

	center := mgl32.Ident4()
	pos0 := center.Mul4(mgl32.Translate3D(-0.5*float32(gridW), 0.5*float32(gridH), 0))

	fieldRenderers := make([]*render.Field, n)

	curr := pos0
	gridWCurr, gridHCurr, k = 0, 0, 0
	for idx := range n {
		gridWCurr = fieldSizes[idx].w
		gridHCurr = max(fieldSizes[idx].h, gridHCurr)
		k++

		curr = curr.Mul4(mgl32.Translate3D(0.5*float32(gridWCurr), -0.5*float32(gridHCurr), 0))
		fieldRenderers[idx] = render.NewField(curr, res, text, str, idx, rr, psides[idx])
		curr = curr.Mul4(mgl32.Translate3D(0.5*float32(gridWCurr), 0.5*float32(gridHCurr), 0))

		if k >= fieldsInRow {
			pos0 = pos0.Mul4(mgl32.Translate3D(0, -float32(gridHCurr), 0))
			curr = pos0
			gridWCurr, k = 0, 0
		}
	}

	return gridW, gridH, fieldRenderers
}
