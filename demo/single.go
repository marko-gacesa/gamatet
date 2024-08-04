// Copyright (c) 2020-2024 by Marko Gaćeša

package demo

import (
	"context"
	"fmt"
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/graphics/render"
	"gamatet/graphics/scene"
	"gamatet/graphics/texture"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/marko-gacesa/appctx"
	"log"
	"math"
	"math/rand"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

func Single() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	for _, monitor := range glfw.GetMonitors() {
		sizeW, sizeH := monitor.GetPhysicalSize()
		workX, workY, workW, workH := monitor.GetWorkarea()
		d := math.Sqrt(float64(sizeW*sizeW+sizeH*sizeH)) / 25.4
		fmt.Printf("Monitor: %s (size=%dmm x %dmm D=%.2f\") area: X=%d Y=%d W=%d H=%d\n", monitor.GetName(), sizeW, sizeH, d, workX, workY, workW, workH)

		for _, videoMode := range monitor.GetVideoModes() {
			fmt.Printf("\tVideo mode: %dx%d@%dHz %dx%dx%d\n", videoMode.Width, videoMode.Height, videoMode.RefreshRate, videoMode.RedBits, videoMode.GreenBits, videoMode.BlueBits)
		}
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	/*
		window := func(monitor *glfw.Monitor) *glfw.Window {
			// https://www.glfw.org/docs/latest/window_guide.html#window_windowed_full_screen

			videoMode := monitor.GetVideoMode()
			fmt.Printf("Current Video Mode: %dx%d@%dHz %dx%dx%d\n", videoMode.Width, videoMode.Height, videoMode.RefreshRate, videoMode.RedBits, videoMode.GreenBits, videoMode.BlueBits)

			glfw.WindowHint(glfw.RedBits, videoMode.RedBits)
			glfw.WindowHint(glfw.GreenBits, videoMode.GreenBits)
			glfw.WindowHint(glfw.BlueBits, videoMode.BlueBits)
			glfw.WindowHint(glfw.RefreshRate, videoMode.RefreshRate)

			window, err := glfw.CreateWindow(videoMode.Width, videoMode.Height, "GaMaTet", monitor, nil)
			if err != nil {
				panic(err)
			}

			return window
		}(glfw.GetPrimaryMonitor())
		defer window.Destroy()
	//*/

	//*
	window := func() *glfw.Window {
		glfw.WindowHint(glfw.Resizable, glfw.True)

		windowWidth := 900
		windowHeight := 600

		window, err := glfw.CreateWindow(windowWidth, windowHeight, "GaMaTet", nil, nil)
		if err != nil {
			panic(err)
		}

		return window
	}()
	defer window.Destroy()
	//*/

	window.MakeContextCurrent()
	//window.SetOpacity(0.5)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure global settings

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// To render transparent:
	// * gl.Disable(gl.CULL_FACE)
	// * gl.Disable(gl.DEPTH_TEST)
	// * gl.Enable(gl.BLEND)
	// * gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	//   or gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	const fieldW = 10
	const fieldH = 22
	const contentW = 2 * (3 + 1 + fieldW)
	const contentH = fieldH + 2

	rend := render.NewRenderer()
	defer rend.Release()

	setupCamera := func(w, h int) {
		rend.PerspectiveFull(w, h, contentW, contentH, 2)
	}
	setupCamera(window.GetFramebufferSize())

	window.SetSizeCallback(func(window *glfw.Window, w int, h int) {
		fmt.Printf("Size Callback %dx%d\n", w, h)
	})

	window.SetFramebufferSizeCallback(func(window *glfw.Window, w int, h int) {
		fmt.Printf("FramebufferSize Callback %dx%d\n", w, h)
		setupCamera(w, h)
		gl.Viewport(0, 0, int32(w), int32(h))
	})

	ctx, cancelFunc := context.WithCancel(appctx.Context)
	defer cancelFunc()

	window.SetCloseCallback(func(w *glfw.Window) {
		fmt.Printf("Close Callback\n")
		cancelFunc()
	})

	// GAME SETUP

	fieldServerCh := make(chan []byte, 100)
	fieldClientCh := make(chan []byte, 100)
	r := rand.New(rand.NewSource(123))
	_ = r
	go func() {
		for e := range fieldServerCh {
			//time.Sleep(time.Millisecond * time.Duration(30+r.Intn(100)))
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

	// DONE GAME SETUP

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey) {
		//fmt.Println("KEY", key, scancode, action, mods)

		if act != glfw.Press {
			return
		}

		if key == glfw.KeyEscape {
			fmt.Println("ESCAPE")
			cancelFunc()
			return
		}

		switch key {
		case glfw.KeyLeft:
			playerInCh <- []byte{byte(action.MoveLeft)}
		case glfw.KeyRight:
			playerInCh <- []byte{byte(action.MoveRight)}
		case glfw.KeyUp:
			playerInCh <- []byte{byte(action.RotateCCW)}
		case glfw.KeyDown:
			playerInCh <- []byte{byte(action.MoveDown)}
		case glfw.KeySpace:
			playerInCh <- []byte{byte(action.Drop)}
		case glfw.KeyP, glfw.KeyPause:
			playerInCh <- []byte{byte(action.Pause)}
		}
	})

	center := mgl32.Ident4()
	left := center.Mul4(mgl32.Translate3D(-float32(contentW)/4, 0, 0))
	right := center.Mul4(mgl32.Translate3D(float32(contentW)/4, 0, 0))

	texManager := texture.Init()

	fieldResources := render.GenerateFieldResources(texManager)
	defer fieldResources.Release()

	textRender := render.MakeText(texManager, render.Font)

	objLeftField := scene.NewField(fieldResources, textRender, 0, game)
	defer objLeftField.Release()

	objRightField := scene.NewField(fieldResources, textRender, 0, gameClient)
	defer objRightField.Release()

out:
	for {
		select {
		case <-ctx.Done():
			break out
		default:
		}

		now := time.Now()

		objLeftField.Prepare(ctx, &left, now)
		objRightField.Prepare(ctx, &right, now)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		objLeftField.Render(rend)
		objRightField.Render(rend)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	fmt.Println("LOOP DONE, waiting")

	wg.Wait()
}
