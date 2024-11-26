// Copyright (c) 2023,2024 by Marko Gaćeša

package loop

import (
	"context"
	"errors"
	"fmt"
	"gamatet/graphics/scene"
	"gamatet/internal/app"
	"gamatet/logic/screen"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/marko-gacesa/appctx"
	"math"
	"runtime"
	"time"
)

func Loop(app *app.App) error {
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

	var window *glfw.Window

	if err := func() (err error) {
		//window, err = windowFullscreen("GaMaTeT", glfw.GetPrimaryMonitor())
		window, err = windowResizable(900, 600, "GaMaTeT")
		return
	}(); err != nil {
		return fmt.Errorf("failed to create window: %w", err)
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
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// Main resources

	ctxLoop, cancelLoopCtxFn := context.WithCancel(appctx.Context)
	defer cancelLoopCtxFn()

	resources := scene.InitResources()
	app.SetScreener(resources)

	var scr screen.Screen

	// Callbacks

	window.SetSizeCallback(func(window *glfw.Window, w int, h int) {})
	window.SetFramebufferSizeCallback(func(window *glfw.Window, w int, h int) {
		if scr != nil {
			scr.UpdateViewSize(w, h)
		}
		gl.Viewport(0, 0, int32(w), int32(h))
	})

	window.SetCloseCallback(func(w *glfw.Window) {
		cancelLoopCtxFn()
	})

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey) {
		if act != glfw.Press {
			return
		}

		//fmt.Printf("KEY key=%d (%c), scan=%04x, act=%d, mods=%06b\n", key, key, scancode, act, mods)
		/*
			if key == glfw.KeyEscape {
				fmt.Println("ESCAPE")
				if fn := cancelScreenFunc; fn != nil {
					fn()
				}
				return
			}
		*/

		if scr != nil {
			scr.InputKeyPress(int(key), scancode)
		}
	})

	window.SetCharCallback(func(w *glfw.Window, char rune) {
		if scr != nil {
			scr.InputChar(char)
		}
	})

	// Application loop

	for ctxLoop.Err() == nil {
		err := func(ctx context.Context) error {
			// create the screen's context
			ctx, cancelFn := context.WithCancel(ctxLoop)

			// create the screen
			scr = app.MakeScreen(ctx)
			if scr == nil { // no screen means exit the app
				cancelFn()        // first cancel the screen's context
				cancelLoopCtxFn() // then cancel the loop's context
				return nil
			}

			defer func() {
				cancelFn()
				scr.Release()
			}()

			scr.UpdateViewSize(window.GetFramebufferSize())

			for {
				select {
				case <-ctx.Done():
					// If the context is done (before the screen) that means that the termination came from outside.
					return ctx.Err()
				case err := <-scr.Done():
					// If the screen is done that means that the termination came from the screen itself.
					return err
				default:
				}

				now := time.Now()

				scr.Prepare(ctx, now)

				gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

				scr.Render(ctx)

				window.SwapBuffers()
				glfw.PollEvents()
			}
		}(ctxLoop)
		if err != nil && !errors.Is(err, context.Canceled) {
			return err
		}

		app.ScreenFinish()
	}

	return nil
}

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
