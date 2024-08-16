// Copyright (c) 2023,2024 by Marko Gaćeša

package loop

import (
	"context"
	"fmt"
	"gamatet/graphics/render"
	"gamatet/graphics/screen"
	"gamatet/graphics/screen/menu"
	"gamatet/graphics/texture"
	"gamatet/internal/config"
	"gamatet/internal/router"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/marko-gacesa/appctx"
	"math"
	"runtime"
	"time"
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

func Loop(cfg *config.Config, router router.Router) error {
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
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	tex := texture.Init()

	ctx, cancelFunc := context.WithCancel(appctx.Context)
	defer cancelFunc()

	var scr screen.Screen
	r := render.NewRenderer()

	window.SetSizeCallback(func(window *glfw.Window, w int, h int) {})
	window.SetFramebufferSizeCallback(func(window *glfw.Window, w int, h int) {
		if scr != nil {
			scr.SetCamera(r, w, h)
		}
		gl.Viewport(0, 0, int32(w), int32(h))
	})

	window.SetCloseCallback(func(w *glfw.Window) {
		cancelFunc()
	})

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, act glfw.Action, mods glfw.ModifierKey) {
		//fmt.Printf("KEY key=%d (%c), scan=%04x, act=%d, mods=%06b\n", key, key, scancode, act, mods)
		if key == glfw.KeyEscape {
			fmt.Println("ESCAPE")
			cancelFunc()
			return
		}

		if scr != nil {
			scr.InputKey(key, scancode, act, mods)
		}
	})

	window.SetCharCallback(func(w *glfw.Window, char rune) {
		if scr != nil {
			scr.InputChar(char)
		}
	})

	//scr = demoblocks.NewDemoBlocks(tex)
	//scr = fieldtest2.NewFieldTest(ctx, tex)
	scr = menu.NewMenu(ctx, tex, router[""]...)
	defer scr.Release()

	func() {
		w, h := window.GetFramebufferSize()
		scr.SetCamera(r, w, h)
	}()

	/////////////
	/*
		var route string
		for {
			scrGen := getScreenGen(route)
			if scrGen == nil {
				break
			}

			func(appCtx context.Context, scrGen func(ctx context.Context, tex *texture.Manager)) {
				ctx, cancelFn := context.WithCancel(appCtx)
				defer cancelFn()

				scr = scrGen(ctx, tex)
				defer scr.Release()

			out:
				for {
					select {
					case <-ctx.Done():
						break out
					default:
					}

					now := time.Now()

					scr.Prepare(ctx, now)

					gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

					scr.Render(r)

					window.SwapBuffers()
					glfw.PollEvents()
				}

				scr.Shutdown()
			}(appCtx, scrGen)
		}
		//*/
	///////////////
out:
	for {
		select {
		case <-ctx.Done():
			break out
		default:
		}

		now := time.Now()

		scr.Prepare(ctx, now)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		scr.Render(r)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	scr.Shutdown()

	return nil
}
