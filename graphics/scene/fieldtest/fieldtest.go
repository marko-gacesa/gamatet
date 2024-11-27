// Copyright (c) 2020-2024 by Marko Gaćeša

package fieldtest

import (
	"context"
	"gamatet/game"
	"gamatet/game/action"
	"gamatet/graphics/render"
	"gamatet/graphics/texture"
	"gamatet/logic/screen"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"sync"
	"time"
)

const (
	fieldW   = 10
	fieldH   = 22
	contentW = 2 * (3 + 1 + fieldW)
	contentH = fieldH + 2
)

type FieldTest struct {
	renderer  *render.Renderer
	tex       *texture.Manager
	res       render.FieldResources
	textLeft  render.Text
	textRight render.Text
	fps       render.FPS

	stopper *screen.Stopper

	playerInCh chan<- []byte

	modelLeft  mgl32.Mat4
	modelRight mgl32.Mat4

	fieldLeft  *render.Field
	fieldRight *render.Field

	wg *sync.WaitGroup
}

var _ screen.Screen = (*FieldTest)(nil)

func NewFieldTest(ctx context.Context, renderer *render.Renderer, tex *texture.Manager) *FieldTest {
	res := render.GenerateFieldResources(tex)
	textLeft := render.MakeText(tex, render.Font)
	textRight := render.MakeText(tex, render.Font)
	fps := render.NewFPS()

	modelCenter := mgl32.Ident4()
	modelLeft := modelCenter.Mul4(mgl32.Translate3D(-float32(contentW)/4, 0, 0))
	modelRight := modelCenter.Mul4(mgl32.Translate3D(float32(contentW)/4, 0, 0))

	gameHost, gameInterpreter, playerInCh, wg := game.NewFieldTest(ctx, fieldW, fieldH)

	ft := &FieldTest{
		renderer:   renderer,
		tex:        tex,
		res:        *res,
		textLeft:   *textLeft,
		textRight:  *textRight,
		fps:        *fps,
		stopper:    screen.NewStopper(),
		playerInCh: playerInCh,
		modelLeft:  modelLeft,
		modelRight: modelRight,
		fieldLeft:  nil,
		fieldRight: nil,
		wg:         wg,
	}

	ft.fieldLeft = render.NewField(ft.modelLeft, &ft.res, &ft.textLeft, 0, gameHost)
	ft.fieldRight = render.NewField(ft.modelRight, &ft.res, &ft.textRight, 0, gameInterpreter)

	return ft
}

func (ft *FieldTest) Done() <-chan error { return ft.stopper.Done() }

func (ft *FieldTest) Release() {
	ft.wg.Wait()
	ft.textLeft.Release()
	ft.textRight.Release()
	ft.res.Release()
}

func (ft *FieldTest) UpdateViewSize(w, h int) {
	ft.renderer.PerspectiveFull(w, h, contentW, contentH, 2)
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
		ft.stopper.Stop()
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
	ft.fps.Render(r, &ft.textLeft, mgl32.Translate3D(-contentW/2+0.5, -contentH/2+1.5, 1.5))
}
