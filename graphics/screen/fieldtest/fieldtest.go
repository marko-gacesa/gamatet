// Copyright (c) 2020-2024 by Marko Gaćeša

package fieldtest

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
	"math/rand"
	"runtime/debug"
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
	tex  *texture.Manager
	res  render.FieldResources
	text render.Text
	fps  render.FPS

	playerInCh chan<- []byte

	modelLeft  mgl32.Mat4
	modelRight mgl32.Mat4

	fieldLeft  *render.Field
	fieldRight *render.Field

	wg *sync.WaitGroup
}

var _ screen.Screen = (*FieldTest)(nil)

func NewFieldTest(ctx context.Context, tex *texture.Manager) *FieldTest {
	res := render.GenerateFieldResources(tex)
	text := render.MakeText(tex, render.Font)
	fps := render.NewFPS()

	modelCenter := mgl32.Ident4()
	modelLeft := modelCenter.Mul4(mgl32.Translate3D(-float32(contentW)/4, 0, 0))
	modelRight := modelCenter.Mul4(mgl32.Translate3D(float32(contentW)/4, 0, 0))

	// event transfer between the host and the client

	fieldServerCh := make(chan []byte, 100) // capacity because of the artificial delay to simulate latency
	fieldClientCh := make(chan []byte, 100)

	r := rand.New(rand.NewSource(123))
	_ = r
	go func() {
		for e := range fieldServerCh {
			time.Sleep(time.Millisecond * time.Duration(30+r.Intn(100)))
			select {
			case <-ctx.Done():
				return
			case fieldClientCh <- e:
				/*
					if len(e) > 1 && e[0] == 'Z' {
						fmt.Printf("cli<-srv fIdx=1 len=%-3d compressed=%x\n", len(e), e)
						continue
					}

					b := bytes.NewBuffer(nil)
					w := gzip.NewWriter(b)
					w.Write(e)
					w.Close()
					compressed := b.Len()
					fmt.Printf("cli<-srv fIdx=1 len=%-3d comp=%-3d raw=%x\n", len(e), compressed, e)
				*/
			}
		}
	}()

	// game setup

	playerInCh, playerOutCh := core.ChPair(ctx)

	const seed = 101
	const level = 2

	setup := core.Setup{
		Name: "test game",
		Config: core.GameConfig{
			WidthPerPlayer: fieldW,
			Height:         fieldH,
			Level:          level,
			PlayerZones:    true,
			FieldConfig: field.Config{
				PieceCollision: false,
				Anim:           false,
			},
			RandomSeed:  seed,
			FeedBagSize: 2,
		},
		Fields: []core.FieldSetup{
			{
				OutCh: fieldServerCh,
				Players: []core.PlayerSetup{
					{
						Name: "marko",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideEnabled:        true,
							MaxWallKick:         2,
						},
						InCh: playerOutCh,
					},
				},
			},
		},
	}

	game := core.MakeHost(setup)

	setupClient := core.Setup{
		Name: "test game",
		Config: core.GameConfig{
			WidthPerPlayer: fieldW,
			Height:         fieldH,
			Level:          level,
			PlayerZones:    true,
			FieldConfig: field.Config{
				PieceCollision: false,
				Anim:           true,
			},
			RandomSeed:  seed,
			FeedBagSize: 2,
		},
		Fields: []core.FieldSetup{
			{
				InCh: fieldClientCh,
				Players: []core.PlayerSetup{
					{
						Name: "marko",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideEnabled:        true,
							MaxWallKick:         2,
						},
					},
				},
			},
		},
	}

	gameClient := core.MakeInterpreter(setupClient)

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

	wg.Add(1)
	go func(ctx context.Context) {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Printf("PANIC: %v\n%s\n", r, debug.Stack())
			}
		}()
		defer wg.Done()

		gameClient.Perform(ctx)

		fmt.Println("GAME CLIENT STOPPED")
	}(ctx)

	ft := &FieldTest{
		tex:        tex,
		res:        *res,
		text:       *text,
		fps:        *fps,
		playerInCh: playerInCh,
		modelLeft:  modelLeft,
		modelRight: modelRight,
		fieldLeft:  nil,
		fieldRight: nil,
		wg:         wg,
	}

	ft.fieldLeft = render.NewField(ft.modelLeft, &ft.res, &ft.text, 0, game)
	ft.fieldRight = render.NewField(ft.modelRight, &ft.res, &ft.text, 0, gameClient)

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

func (ft *FieldTest) KeyAction(key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey) {
	if act != glfw.Press {
		return
	}

	switch key {
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
	}
}

func (ft *FieldTest) Prepare(ctx context.Context, now time.Time) {
	ft.fieldLeft.Prepare(ctx, now)
	ft.fieldRight.Prepare(ctx, now)
}

func (ft *FieldTest) Render(r *render.Renderer) {
	ft.fieldLeft.Render(r)
	ft.fieldRight.Render(r)
	ft.fps.Render(r, &ft.text, mgl32.Translate3D(-contentW/2+0.5, -contentH/2+1.5, 1.5))
}
