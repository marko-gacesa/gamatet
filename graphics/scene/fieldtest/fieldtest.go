// Copyright (c) 2020-2024 by Marko Gaćeša

package fieldtest

import (
	"context"
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/graphics/render"
	"gamatet/graphics/texture"
	"gamatet/logic/screen"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type FieldTest struct {
	renderer  *render.Renderer
	tex       *texture.Manager
	res       render.FieldResources
	textLeft  render.Text
	textRight render.Text
	fps       render.FPS

	usePerspective bool

	contentW int
	contentH int
	viewW    int
	viewH    int

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
	textLeft := render.MakeText(tex, render.Font)
	textRight := render.MakeText(tex, render.Font)
	fps := render.NewFPS()

	contentWLeft, contentHLeft := render.GetExtendedContent(gameHost.GetSize(0))
	contentWRight, contentHRight := render.GetExtendedContent(gameInterpreter.GetSize(0))

	contentW := contentWLeft + contentWRight
	contentH := max(contentHLeft, contentHRight)

	modelCenter := mgl32.Ident4()
	modelLeft := modelCenter.Mul4(mgl32.Translate3D(-float32(contentW)/4, 0, 0))
	modelRight := modelCenter.Mul4(mgl32.Translate3D(float32(contentW)/4, 0, 0))

	ft := &FieldTest{
		renderer:   renderer,
		tex:        tex,
		res:        *res,
		textLeft:   *textLeft,
		textRight:  *textRight,
		fps:        *fps,
		contentW:   contentW,
		contentH:   contentH,
		viewW:      0,
		viewH:      0,
		playerInCh: playerInCh,
		modelLeft:  modelLeft,
		modelRight: modelRight,
		fieldLeft:  nil,
		fieldRight: nil,
		waitDoneCh: waitDoneCh,
	}

	ft.fieldLeft = render.NewField(ft.modelLeft, &ft.res, &ft.textLeft, 0, gameHost)
	ft.fieldRight = render.NewField(ft.modelRight, &ft.res, &ft.textRight, 0, gameInterpreter)

	return ft
}

func (ft *FieldTest) Release() {
	<-ft.waitDoneCh
	ft.textLeft.Release()
	ft.textRight.Release()
	ft.res.Release()
}

func (ft *FieldTest) UpdateViewSize(w, h int) {
	ft.viewW, ft.viewH = w, h
	if ft.usePerspective {
		ft.renderer.PerspectiveFull(w, h, ft.contentW, ft.contentH, 2)
	} else {
		ft.renderer.OrthogonalFull(w, h, ft.contentW, ft.contentH, 2)
	}
}

func (ft *FieldTest) InputKeyPress(key, scancode int) {
	switch glfw.Key(key) {
	case glfw.KeyLeft:
		ft.playerInCh <- []byte{byte(action.MoveLeft)}
	case glfw.KeyRight:
		ft.playerInCh <- []byte{byte(action.MoveRight)}
	case glfw.KeyUp:
		ft.playerInCh <- []byte{byte(action.RotateCCW)}
	case glfw.KeyDown:
		ft.playerInCh <- []byte{byte(action.MoveDown)}
	case glfw.KeySpace:
		ft.playerInCh <- []byte{byte(action.Drop)}
	case glfw.KeyP, glfw.KeyPause:
		ft.playerInCh <- []byte{byte(action.Pause)}
	case glfw.KeyEscape:
		ft.playerInCh <- []byte{byte(action.Abort)}
	case glfw.KeyF12:
		ft.usePerspective = !ft.usePerspective
		if ft.usePerspective {
			ft.renderer.PerspectiveFull(ft.viewW, ft.viewH, ft.contentW, ft.contentH, 2)
		} else {
			ft.renderer.OrthogonalFull(ft.viewW, ft.viewH, ft.contentW, ft.contentH, 2)
		}
	}
}

func (ft *FieldTest) InputChar(char rune) {}

func (ft *FieldTest) Prepare(ctx context.Context, now time.Time) {
	ft.fieldLeft.Prepare(ctx, now)
	ft.fieldRight.Prepare(ctx, now)
}

func (ft *FieldTest) Render(ctx context.Context) {
	r := ft.renderer
	ft.fieldLeft.Render(r)
	ft.fieldRight.Render(r)
	ft.fps.Render(r, &ft.textLeft, mgl32.Translate3D(float32(-ft.contentW)/2+0.5, -float32(ft.contentH)/2+1.5, 1.5))
}
