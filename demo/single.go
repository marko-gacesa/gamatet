// Copyright (c) 2020 by Marko Gaćeša

package demo

import (
	"context"
	"fmt"
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/graphics/camera"
	"gamatet/graphics/geometry"
	"gamatet/graphics/material"
	"gamatet/graphics/renderer"
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
		//defer window.Destroy()
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
	//defer window.Destroy()
	//*/

	window.MakeContextCurrent()

	var windowWidth int
	var windowHeight int
	windowWidth, windowHeight = window.GetSize()
	//window.SetOpacity(0.5)

	func(window *glfw.Window) {
		w, h := window.GetSize()
		fmt.Println("Window size", w, "x", h)
		sx, sy := window.GetContentScale()
		fmt.Println("Content scale", sx, "x", sy)
		w, h = window.GetFramebufferSize()
		fmt.Println("Framebuffer size", w, "x", h)
	}(window)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// calculate distance
	const fovDeg = 45
	const fovRad = fovDeg * math.Pi / 180

	contentW := 10
	contentH := 24
	aspectRatio := float32(windowWidth) / float32(windowHeight)

	fmt.Printf("fov: %f; aspect ratio %f\n", fovRad, aspectRatio)

	cameraDistance := float32(math.Max(float64(contentW+2)/(2*math.Tan(fovRad/2)), float64(contentH+2)/(2*math.Tan(fovRad/2))))

	fmt.Println("camera distance X", float64(contentW)/(2*math.Tan(fovRad/2)))
	fmt.Println("camera distance Y", float64(contentH)/(2*math.Tan(fovRad/2)))
	fmt.Println("camera distance", cameraDistance)

	cam := camera.NewCamera()
	cam.LookAt(mgl32.Vec3{0, 0, cameraDistance}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cam.Perspective(fovRad, aspectRatio, cameraDistance-12, cameraDistance+12)

	texture.Instance = texture.Init()
	geometry.LoadAll()

	img := texture.Tex2D(345)
	texture.Instance = texture.Init()
	tex := texture.Instance.Bind(img)
	matWall := material.NewRock(tex)
	//matRuby := material.NewRuby()
	matLava := material.NewLava(tex)
	matAcid := material.NewAcid(tex)

	fieldRenderer := renderer.NewFieldRender(matWall, matLava, matAcid)

	rend := &renderer.Renderer{}
	rend.Geometry(geometry.Cube)
	rend.Camera(cam)

	// Configure global settings

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.FRONT)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// To render transparent:
	// * gl.Disable(gl.CULL_FACE)
	// * gl.Disable(gl.DEPTH_TEST)
	// * gl.Enable(gl.BLEND)
	// * gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	//   or gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	ctx, cancelFunc := context.WithCancel(appctx.Context)
	defer cancelFunc()

	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		fmt.Printf("Size Callback %dx%d\n", width, height)
		aspectRatio = float32(width) / float32(height)
		gl.Viewport(0, 0, int32(width), int32(height))
		cam.LookAt(mgl32.Vec3{0, 0, cameraDistance}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
		cam.Perspective(fovRad, aspectRatio, cameraDistance-12, cameraDistance+12)
	})

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		fmt.Printf("FramebufferSize Callback\n")
	})

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

	setup := core.Setup{
		Name: "test game",
		Config: core.GameConfig{
			WidthPerPlayer: 10,
			Height:         22,
			Level:          6,
			PlayerZones:    true,
			FieldConfig: field.Config{
				PieceCollision: false,
				Anim:           false,
			},
			RandomSeed:  123,
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

	game := core.MakeSession(setup)

	setupClient := core.Setup{
		Name: "test game",
		Config: core.GameConfig{
			WidthPerPlayer: 10,
			Height:         22,
			Level:          6,
			PlayerZones:    true,
			FieldConfig: field.Config{
				PieceCollision: false,
				Anim:           true,
			},
			RandomSeed:  123,
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

	chRenderInfoServer := make(chan *field.RenderInfo, 1)
	chRenderInfoClient := make(chan *field.RenderInfo, 1)

	previousTime := glfw.GetTime()

out:
	for {
		select {
		case <-ctx.Done():
			break out
		default:
		}

		//*
		t := glfw.GetTime()
		elapsed := t - previousTime
		previousTime = t
		_ = elapsed
		//*/

		now := time.Now()

		game.RenderRequest(ctx, 0, now, chRenderInfoServer)
		gameClient.RenderRequest(ctx, 0, now, chRenderInfoClient)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		center := mgl32.Ident4()

		///////////////////////
		//camera = mgl32.LookAtV(mgl32.Vec3{cameraDistance * float32(math.Sin(t)), 0, cameraDistance * float32(math.Cos(t))}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
		//projection = mgl32.Perspective(fovRad, aspectRatio, cameraDistance-12, cameraDistance+12)
		///////////////////////
		/*
			angle := t
			drawBigBlock := func(position mgl32.Mat4, x, y float32) {
				const scale = 8
				bigBlock := position.
					Mul4(mgl32.HomogRotate3DZ(float32(angle / 6))).
					//Mul4(mgl32.HomogRotate3DY(float32(angle / 2.7))).
					//Mul4(mgl32.HomogRotate3DX(float32(angle / 1.2))).
					Mul4(mgl32.Scale3D(scale, scale, scale)).
					Mul4(mgl32.Translate3D(x, y, 0))
				rend.Render(&bigBlock)
			}
			rend.Geometry(geometry.Gem)
			rend.Material(matAcid)
			matLava.Color(mgl32.Vec4{1, 1, 1, 1})
			drawBigBlock(center, -0.5, -0.5)
			drawBigBlock(center, -0.5, 0.5)
			//rend.Material(matAcid)
			matAcid.Color(mgl32.Vec4{1, 1, 1, 1})
			drawBigBlock(center, 0.5, -0.5)
			drawBigBlock(center, 0.5, 0.5)
			_ = angle
			//*/
		///////////////////////

		//left := center.Mul4(mgl32.Translate3D(-6, 0, 0))
		right := center.Mul4(mgl32.Translate3D(6, 0, 0))

		/*
			select {
			case renderInfo := <-chRenderInfoServer:
				fieldRenderer.Render(rend, &left, renderInfo)
				field.ReturnRenderInfo(renderInfo)
			}
			//*/

		select {
		case renderInfo := <-chRenderInfoClient:
			fieldRenderer.Render(rend, &right, renderInfo)
			field.ReturnRenderInfo(renderInfo)
		}

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	fmt.Println("LOOP DONE, waiting")

	wg.Wait()

	window.Destroy()
}
