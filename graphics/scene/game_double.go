// Copyright (c) 2024 by Marko Gaćeša

package scene

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

type GameDouble struct {
	renderer *render.Renderer
	tex      *texture.Manager
	res      render.FieldResources
	text     render.Text
	fps      render.FPS

	usePerspective bool

	contentW int
	contentH int
	viewW    int
	viewH    int

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
	fps := render.NewFPS()

	w, h := render.GetExtendedContent(params.Game.GetSize(0))

	g := &GameDouble{
		renderer:       renderer,
		tex:            tex,
		res:            *res,
		text:           *text,
		fps:            *fps,
		usePerspective: false, // TODO: Read from config

		contentW: w,
		contentH: h,
		viewW:    0,
		viewH:    0,

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
}

func (ft *GameDouble) UpdateViewSize(w, h int) {
	ft.viewW, ft.viewH = w, h
	if ft.usePerspective {
		ft.renderer.PerspectiveFull(w, h, ft.contentW, ft.contentH, 2)
	} else {
		ft.renderer.OrthogonalFull(w, h, ft.contentW, ft.contentH, 2)
	}
}

func (ft *GameDouble) InputKeyPress(key, scancode int) {
	switch glfw.Key(key) {
	case glfw.KeyEscape:
		ft.player1InCh <- []byte{byte(action.Abort)}
	case glfw.KeyPause:
		ft.player1InCh <- []byte{byte(action.Pause)}

	case glfw.KeyLeft:
		ft.player1InCh <- []byte{byte(action.MoveLeft)}
	case glfw.KeyRight:
		ft.player1InCh <- []byte{byte(action.MoveRight)}
	case glfw.KeyUp:
		ft.player1InCh <- []byte{byte(action.RotateCCW)}
	case glfw.KeyDown:
		ft.player1InCh <- []byte{byte(action.MoveDown)}
	case glfw.KeyRightShift:
		ft.player1InCh <- []byte{byte(action.Drop)}

	case glfw.KeyA:
		ft.player2InCh <- []byte{byte(action.MoveLeft)}
	case glfw.KeyD:
		ft.player2InCh <- []byte{byte(action.MoveRight)}
	case glfw.KeyW:
		ft.player2InCh <- []byte{byte(action.RotateCCW)}
	case glfw.KeyS:
		ft.player2InCh <- []byte{byte(action.MoveDown)}
	case glfw.KeyLeftShift:
		ft.player2InCh <- []byte{byte(action.Drop)}

	case glfw.KeyF12:
		ft.usePerspective = !ft.usePerspective
		if ft.usePerspective {
			ft.renderer.PerspectiveFull(ft.viewW, ft.viewH, ft.contentW, ft.contentH, 2)
		} else {
			ft.renderer.OrthogonalFull(ft.viewW, ft.viewH, ft.contentW, ft.contentH, 2)
		}
	}
}

func (ft *GameDouble) InputChar(char rune) {}

func (ft *GameDouble) Prepare(ctx context.Context, now time.Time) {
	ft.fieldRender.Prepare(ctx, now)
}

func (ft *GameDouble) Render(ctx context.Context) {
	r := ft.renderer
	ft.fieldRender.Render(r)
	ft.fps.Render(r, &ft.text, mgl32.Translate3D(float32(-ft.contentW)/2+0.5, float32(-ft.contentH)/2+1.5, 1.5))
}
