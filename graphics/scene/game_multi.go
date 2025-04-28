// Copyright (c) 2025 by Marko Gaćeša

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

type GameMulti struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text
	fps  render.FPS

	playersInCh  [4]chan<- []byte
	model        mgl32.Mat4
	game         core.RenderRequester
	fieldRenders []*render.Field

	waitDoneCh <-chan struct{}
}

var _ screen.Screen = (*GameMulti)(nil)

func NewGameMulti(
	renderer *render.Renderer,
	tex *texture.Manager,
	params core.GameParams,
) *GameMulti {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)
	fps := render.NewFPS()

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

	g := &GameMulti{
		BlockBase: base.NewBlockBase(renderer, tex, w, h, true),
		res:       *res,
		text:      *text,
		fps:       *fps,

		// these are set below
		playersInCh:  params.PlayerInCh,
		model:        center,
		game:         params.Game,
		fieldRenders: fieldRenders,

		waitDoneCh: params.Done,
	}

	return g
}

func (ft *GameMulti) Release() {
	<-ft.waitDoneCh
	ft.text.Release()
	ft.res.Release()

	for i := range ft.playersInCh {
		if ft.playersInCh[i] != nil {
			close(ft.playersInCh[i])
		}
	}
}

func (ft *GameMulti) InputKeyPress(key, scancode int) {
	ft.BlockBase.InputKeyPress(key, scancode)

	var cmd1, cmd2, cmd3, cmd4 []byte

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
		cmd1 = []byte{byte(action.RotateCCW)}
	case glfw.KeyDown:
		cmd1 = []byte{byte(action.MoveDown)}
	case glfw.KeyRightShift:
		cmd1 = []byte{byte(action.Drop)}

	case glfw.KeyA:
		cmd2 = []byte{byte(action.MoveLeft)}
	case glfw.KeyD:
		cmd2 = []byte{byte(action.MoveRight)}
	case glfw.KeyW:
		cmd2 = []byte{byte(action.RotateCCW)}
	case glfw.KeyS:
		cmd2 = []byte{byte(action.MoveDown)}
	case glfw.KeyLeftShift:
		cmd2 = []byte{byte(action.Drop)}
	}

	base.SendAction(cmd1, ft.waitDoneCh, ft.playersInCh[0])
	base.SendAction(cmd2, ft.waitDoneCh, ft.playersInCh[1])
	base.SendAction(cmd3, ft.waitDoneCh, ft.playersInCh[2])
	base.SendAction(cmd4, ft.waitDoneCh, ft.playersInCh[3])
}

func (ft *GameMulti) Prepare(now time.Time) {
	for i := range ft.fieldRenders {
		ft.fieldRenders[i].Prepare(now)
	}
}

func (ft *GameMulti) Render() {
	ft.SetCamera()
	r := ft.Renderer()
	for i := range ft.fieldRenders {
		ft.fieldRenders[i].Render(r)
	}
}
