// Copyright (c) 2023,2024 by Marko Gaćeša

package render

import (
	"context"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/marko-gacesa/appctx"
	"math"
	"runtime"
)

func windowResizable(w, h int, title string) (*glfw.Window, error) {
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create resizeable window: %w", err)
	}

	return window, nil
}

func windowFullscreen(title string, monitor *glfw.Monitor) (*glfw.Window, error) {
	// https://www.glfw.org/docs/latest/window_guide.html#window_windowed_full_screen

	videoMode := monitor.GetVideoMode()
	fmt.Printf("Current Video Mode: %dx%d@%dHz %dx%dx%d\n",
		videoMode.Width, videoMode.Height, videoMode.RefreshRate,
		videoMode.RedBits, videoMode.GreenBits, videoMode.BlueBits)

	glfw.WindowHint(glfw.RedBits, videoMode.RedBits)
	glfw.WindowHint(glfw.GreenBits, videoMode.GreenBits)
	glfw.WindowHint(glfw.BlueBits, videoMode.BlueBits)
	glfw.WindowHint(glfw.RefreshRate, videoMode.RefreshRate)

	window, err := glfw.CreateWindow(videoMode.Width, videoMode.Height, title, monitor, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create fullscreen window: %w", err)
	}

	return window, nil
}

func getMonitorInfo() {
	for _, monitor := range glfw.GetMonitors() {
		sizeW, sizeH := monitor.GetPhysicalSize()
		workX, workY, workW, workH := monitor.GetWorkarea()
		d := math.Sqrt(float64(sizeW*sizeW+sizeH*sizeH)) / 25.4
		fmt.Printf("Monitor: %s (size=%dmm x %dmm D=%.2f\") area: X=%d Y=%d W=%d H=%d\n", monitor.GetName(), sizeW, sizeH, d, workX, workY, workW, workH)

		for _, videoMode := range monitor.GetVideoModes() {
			fmt.Printf("\tVideo mode: %dx%d@%dHz %dx%dx%d\n", videoMode.Width, videoMode.Height, videoMode.RefreshRate, videoMode.RedBits, videoMode.GreenBits, videoMode.BlueBits)
		}
	}
}

func Loop() error {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		return fmt.Errorf("failed to initialize glfw: %w", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	//window, err := windowFullscreen("GaMaTeT", glfw.GetPrimaryMonitor())
	window, err := windowResizable(900, 600, "GaMaTeT")
	if err != nil {
		return err
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	//window.SetOpacity(0.5)

	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize OpenGL bindings: %w", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

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

	const fieldW = 10
	const fieldH = 22
	const contentW = 2 * (3 + 1 + fieldW)
	const contentH = fieldH + 2

	GenerateResources()
	defer ReleaseResources()

	rend := &Renderer{}
	rend.Geometry(Resources.GeomCube)

	func() {
		w, h := window.GetFramebufferSize()
		rend.CameraSetDistance(w, h, contentW, contentH, 12)
	}()

	window.SetSizeCallback(func(window *glfw.Window, w int, h int) {
		fmt.Printf("Size Callback %dx%d\n", w, h)
	})

	window.SetFramebufferSizeCallback(func(window *glfw.Window, w int, h int) {
		fmt.Printf("FramebufferSize Callback %dx%d\n", w, h)
		rend.CameraSetDistance(w, h, contentW, contentH, 12)
		gl.Viewport(0, 0, int32(w), int32(h))
	})

	ctx, cancelFunc := context.WithCancel(appctx.Context)
	defer cancelFunc()

	window.SetCloseCallback(func(w *glfw.Window) {
		fmt.Printf("Close Callback\n")
		cancelFunc()
	})

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

		/*
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
		*/
	})

	previousTime := glfw.GetTime()

out:
	for {
		select {
		case <-ctx.Done():
			break out
		default:
		}

		//*/
		t := glfw.GetTime()
		elapsed := t - previousTime
		previousTime = t
		_ = elapsed
		//*/

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		center := mgl32.Ident4()

		///////////////////////
		//*/
		angle := t
		drawBigBlock := func(position mgl32.Mat4, x, y float32) {
			const scale = 8
			bigBlock := position.
				Mul4(mgl32.HomogRotate3DZ(float32(angle / 6))).
				Mul4(mgl32.HomogRotate3DY(float32(angle / 2.7))).
				Mul4(mgl32.HomogRotate3DX(float32(angle / 1.2))).
				Mul4(mgl32.Scale3D(scale, scale, scale)).
				Mul4(mgl32.Translate3D(x, y, 0))
			rend.Render(&bigBlock)
		}

		//rend.Material(Resources.MatTexUV)
		//rend.Material(Resources.MatNorm)

		rend.Geometry(Resources.GeomDie)

		rend.Material(Resources.MatWave)
		Resources.MatWave.Color(mgl32.Vec4{1, 1, 1, 1})
		Resources.MatWave.Texture(Resources.TexRock)
		drawBigBlock(center, -0.5, -0.5)

		rend.Geometry(Resources.GeomRoundedCube)

		rend.Material(Resources.MatRock)
		Resources.MatRock.ChainTexture(Resources.TexChain3)
		Resources.MatRock.Color(mgl32.Vec4{1, 1, 1, 1})
		drawBigBlock(center, -0.5, 0.5)

		rend.Geometry(Resources.GeomSphere)

		rend.Material(Resources.MatRock)
		Resources.MatRock.Color(mgl32.Vec4{1, 1, 1, 1})
		drawBigBlock(center, 0.5, -0.5)

		rend.Geometry(Resources.GeomGem)

		rend.Material(Resources.MatAcid)
		Resources.MatAcid.Color(mgl32.Vec4{1, 1, 1, 1})
		drawBigBlock(center, 0.5, 0.5)
		//*/
		///////////////////////

		window.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}
