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

type GameOne struct {
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

	playerInCh  chan<- []byte
	model       mgl32.Mat4
	game        core.RenderRequester
	fieldRender *render.Field

	waitDoneCh <-chan struct{}
}

var _ screen.Screen = (*GameOne)(nil)

func NewGameOne(
	renderer *render.Renderer,
	tex *texture.Manager,
	params core.GameOneParams,
) *GameOne {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)
	fps := render.NewFPS()

	w, h := render.GetExtendedContent(params.Game.GetSize(0))

	g := &GameOne{
		renderer:       renderer,
		tex:            tex,
		res:            *res,
		text:           *text,
		fps:            *fps,
		usePerspective: true, // TODO: Read from config

		contentW: w,
		contentH: h,
		viewW:    0,
		viewH:    0,

		// these are set below
		playerInCh:  nil,
		model:       mgl32.Mat4{},
		game:        nil,
		fieldRender: nil,

		waitDoneCh: params.Done,
	}

	g.playerInCh = params.PlayerInCh
	g.model = mgl32.Ident4()
	g.game = params.Game
	g.fieldRender = render.NewField(g.model, &g.res, &g.text, 0, g.game)

	return g
}

func (ft *GameOne) Release() {
	<-ft.waitDoneCh
	ft.text.Release()
	ft.res.Release()
}

func (ft *GameOne) UpdateViewSize(w, h int) {
	ft.viewW, ft.viewH = w, h
	if ft.usePerspective {
		ft.renderer.PerspectiveFull(w, h, ft.contentW, ft.contentH, 2)
	} else {
		ft.renderer.OrthogonalFull(w, h, ft.contentW, ft.contentH, 2)
	}
}

func (ft *GameOne) InputKeyPress(key, scancode int) {
	switch glfw.Key(key) {
	case glfw.KeyEscape:
		ft.playerInCh <- []byte{byte(action.Abort)}
	case glfw.KeyPause:
		ft.playerInCh <- []byte{byte(action.Pause)}

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

	case glfw.KeyF12:
		ft.usePerspective = !ft.usePerspective
		if ft.usePerspective {
			ft.renderer.PerspectiveFull(ft.viewW, ft.viewH, ft.contentW, ft.contentH, 2)
		} else {
			ft.renderer.OrthogonalFull(ft.viewW, ft.viewH, ft.contentW, ft.contentH, 2)
		}
	}
}

func (ft *GameOne) InputChar(char rune) {}

func (ft *GameOne) Prepare(ctx context.Context, now time.Time) {
	ft.fieldRender.Prepare(ctx, now)
}

func (ft *GameOne) Render(ctx context.Context) {
	r := ft.renderer
	ft.fieldRender.Render(r)
	ft.fps.Render(r, &ft.text, mgl32.Translate3D(float32(-ft.contentW)/2+0.5, float32(-ft.contentH)/2+1.5, 1.5))
}
