// Copyright (c) 2020-2024 by Marko Gaćeša

package fieldtest2

import (
	"context"
	"fmt"
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/graphics/render"
	"gamatet/graphics/screen"
	"gamatet/graphics/texture"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"runtime/debug"
	"sync"
	"time"
)

const (
	fieldW   = 10
	fieldH   = 22
	contentW = 3 + 1 + 2*fieldW
	contentH = fieldH + 2
)

type FieldTest struct {
	tex  *texture.Manager
	res  render.FieldResources
	text render.Text
	fps  render.FPS

	player1InCh chan<- []byte
	player2InCh chan<- []byte

	field *render.Field

	wg *sync.WaitGroup
}

var _ screen.Screen = (*FieldTest)(nil)

func NewFieldTest(ctx context.Context, tex *texture.Manager) *FieldTest {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)
	fps := render.NewFPS()

	// event transfer between the host and the client

	fieldCh := make(chan []byte)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-fieldCh:
			}
		}
	}()

	// game setup

	player1InCh, player1OutCh := core.ChPair(ctx)
	player2InCh, player2OutCh := core.ChPair(ctx)

	const seed = 101
	const level = 2

	setup := core.Setup{
		Name: "test game",
		Config: core.GameConfig{
			WidthPerPlayer: fieldW,
			Height:         fieldH,
			Level:          level,
			PlayerZones:    false,
			FieldConfig: field.Config{
				PieceCollision: true,
				Anim:           true,
			},
			RandomSeed:  seed,
			FeedBagSize: 2,
		},
		Fields: []core.FieldSetup{
			{
				OutCh: fieldCh,
				Players: []core.PlayerSetup{
					{
						Name: "marko",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideEnabled:        true,
							MaxWallKick:         2,
						},
						InCh: player1OutCh,
					},
					{
						Name: "ogi",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideEnabled:        true,
							MaxWallKick:         2,
						},
						InCh: player2OutCh,
					},
				},
			},
		},
	}

	game := core.MakeHost(setup)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context) {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Printf("PANIC: %v\n%s\n", r, debug.Stack())
			}
		}()
		defer wg.Done()

		game.Perform(ctx)

		fmt.Println("GAME STOPPED")
	}(ctx)

	ft := &FieldTest{
		tex:         tex,
		res:         *res,
		text:        *text,
		fps:         *fps,
		player1InCh: player1InCh,
		player2InCh: player2InCh,
		wg:          wg,
	}

	ft.field = render.NewField(mgl32.Ident4(), &ft.res, &ft.text, 0, game)

	return ft
}

func (ft *FieldTest) Release() {
	ft.res.Release()
	ft.text.Release()
}

func (ft *FieldTest) Shutdown() {
	ft.wg.Wait()
}

func (ft *FieldTest) SetCamera(r *render.Renderer, w, h int) {
	r.PerspectiveFull(w, h, contentW, contentH, 2)
}

func (ft *FieldTest) InputKey(key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey) {
	if act != glfw.Press {
		return
	}

	switch key {
	case glfw.KeyA:
		ft.player1InCh <- []byte{byte(action.MoveLeft)}
	case glfw.KeyD:
		ft.player1InCh <- []byte{byte(action.MoveRight)}
	case glfw.KeyW:
		ft.player1InCh <- []byte{byte(action.RotateCCW)}
	case glfw.KeyS:
		ft.player1InCh <- []byte{byte(action.MoveDown)}
	case glfw.KeyZ:
		ft.player1InCh <- []byte{byte(action.Drop)}
	case glfw.KeyR:
		ft.player1InCh <- []byte{byte(action.Pause)}

	case glfw.KeyLeft:
		ft.player2InCh <- []byte{byte(action.MoveLeft)}
	case glfw.KeyRight:
		ft.player2InCh <- []byte{byte(action.MoveRight)}
	case glfw.KeyUp:
		ft.player2InCh <- []byte{byte(action.RotateCCW)}
	case glfw.KeyDown:
		ft.player2InCh <- []byte{byte(action.MoveDown)}
	case glfw.KeySpace:
		ft.player2InCh <- []byte{byte(action.Drop)}
	case glfw.KeyP, glfw.KeyPause:
		ft.player2InCh <- []byte{byte(action.Pause)}
	}
}

func (ft *FieldTest) InputChar(char rune) {}

func (ft *FieldTest) Prepare(ctx context.Context, now time.Time) {
	ft.field.Prepare(ctx, now)
}

func (ft *FieldTest) Render(r *render.Renderer) {
	ft.field.Render(r)
	ft.fps.Render(r, &ft.text, mgl32.Translate3D(-contentW/2+0.5, -contentH/2+1.5, 1.5))
}
