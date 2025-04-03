// Copyright (c) 2020-2025 by Marko Gaćeša

package fieldtest

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

type FieldTest struct {
	base.BlockBase
	res  render.FieldResources
	text render.Text
	fps  render.FPS

	playerInCh chan<- []byte

	modelLeft  mgl32.Mat4
	modelRight mgl32.Mat4

	fieldLeft  *render.Field
	fieldRight *render.Field

	waitDoneCh <-chan struct{}
}

var _ screen.Screen = (*FieldTest)(nil)

func NewFieldTest(
	renderer *render.Renderer,
	tex *texture.Manager,
	playerInCh chan<- []byte,
	gameHost, gameInterpreter core.RenderRequester,
	waitDoneCh <-chan struct{},
) *FieldTest {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)
	fps := render.NewFPS()

	contentWLeft, contentHLeft := render.GetExtendedContent(gameHost.GetSize(0))
	contentWRight, contentHRight := render.GetExtendedContent(gameInterpreter.GetSize(0))

	contentW := contentWLeft + contentWRight
	contentH := max(contentHLeft, contentHRight)

	modelCenter := mgl32.Ident4()
	modelLeft := modelCenter.Mul4(mgl32.Translate3D(-float32(contentW)/4, 0, 0))
	modelRight := modelCenter.Mul4(mgl32.Translate3D(float32(contentW)/4, 0, 0))

	ft := &FieldTest{
		BlockBase:  base.NewBlockBase(renderer, tex, contentW, contentH, true),
		res:        *res,
		text:       *text,
		fps:        *fps,
		playerInCh: playerInCh,
		modelLeft:  modelLeft,
		modelRight: modelRight,
		fieldLeft:  nil,
		fieldRight: nil,
		waitDoneCh: waitDoneCh,
	}

	ft.fieldLeft = render.NewField(ft.modelLeft, &ft.res, &ft.text, 0, gameHost)
	ft.fieldRight = render.NewField(ft.modelRight, &ft.res, &ft.text, 0, gameInterpreter)

	return ft
}

func (ft *FieldTest) Release() {
	<-ft.waitDoneCh
	ft.text.Release()
	ft.res.Release()
}

func (ft *FieldTest) InputKeyPress(key, scancode int) {
	var cmd []byte
	switch glfw.Key(key) {
	case glfw.KeyLeft:
		cmd = []byte{byte(action.MoveLeft)}
	case glfw.KeyRight:
		cmd = []byte{byte(action.MoveRight)}
	case glfw.KeyUp:
		cmd = []byte{byte(action.RotateCCW)}
	case glfw.KeyDown:
		cmd = []byte{byte(action.MoveDown)}
	case glfw.KeySpace:
		cmd = []byte{byte(action.Drop)}
	case glfw.KeyP, glfw.KeyPause:
		cmd = []byte{byte(action.Pause)}
	case glfw.KeyEscape:
		cmd = []byte{byte(action.Abort)}
	}

	base.SendAction(cmd, ft.waitDoneCh, ft.playerInCh)
}

func (ft *FieldTest) Prepare(now time.Time) {
	ft.fieldLeft.Prepare(now)
	ft.fieldRight.Prepare(now)
}

func (ft *FieldTest) Render() {
	r := ft.Renderer()
	ft.fieldLeft.Render(r)
	ft.fieldRight.Render(r)
}
