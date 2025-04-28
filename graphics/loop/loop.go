// Copyright (c) 2023-2025 by Marko Gaćeša

package loop

import (
	"context"
	"fmt"
	"gamatet/graphics/scene"
	"gamatet/internal/app"
	"gamatet/internal/values"
	"gamatet/logic/screen"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"math"
	"runtime"
	"time"
)

func Loop(globalCtx context.Context, app *app.App) error {
	runtime.LockOSThread() // GLFW event handling must run on the main OS thread

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
		//window, err = windowFullscreen(values.ProgramName, nil)
		window, err = windowResizable(900, 600, values.ProgramName)
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

	log := app.Log()

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Info("OpenGL initialized", "version", version)

	func() {
		w, h := window.GetFramebufferSize()
		log.Info("Framebuffer size", "width", w, "height", h)

		w, h = window.GetSize()
		log.Info("Window size", "width", w, "height", h)

		var videoMode *glfw.VidMode

		monitor := window.GetMonitor()
		if monitor != nil {
			sizeW, sizeH := monitor.GetPhysicalSize()
			d := math.Sqrt(float64(sizeW*sizeW+sizeH*sizeH)) / 25.4
			log.Info("Monitor info",
				"name", monitor.GetName(),
				"width[mm]", sizeW, "height[mm]", sizeH,
				"diagonal[inch]", fmt.Sprintf("%.2f", d))

			videoMode = monitor.GetVideoMode()
		}

		if videoMode != nil {
			log.Info("Video mode",
				"resolution", fmt.Sprintf("%dx%d", videoMode.Width, videoMode.Height),
				"refresh", fmt.Sprintf("%dHz", videoMode.RefreshRate),
				"color", fmt.Sprintf("%dx%dx%d", videoMode.RedBits, videoMode.GreenBits, videoMode.BlueBits))
		}
	}()

	// Configure global settings

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// Main resources

	ctxLoop, cancelLoopCtxFn := context.WithCancel(globalCtx)
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
		func(ctx context.Context) {
			var done <-chan struct{}

			// create the screen
			scr, done = app.MakeScreen(ctx)
			if scr == nil { // no screen means exit the app
				cancelLoopCtxFn()
				return
			}

			defer func() {
				memstats1 := &runtime.MemStats{}
				runtime.ReadMemStats(memstats1)

				scr.Release()
				runtime.GC()

				memstats2 := &runtime.MemStats{}
				runtime.ReadMemStats(memstats2)

				app.Log().Info("Screen done",
					"memory.before.inuse", memstats1.HeapInuse,
					"memory.before.alloc", memstats1.HeapAlloc,
					"memory.after.inuse", memstats2.HeapInuse,
					"memory.after.alloc", memstats2.HeapAlloc,
					"goroutines", runtime.NumGoroutine())
			}()

			scr.UpdateViewSize(window.GetFramebufferSize())

			for isActive(done) {
				scr.Prepare(time.Now())

				gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

				scr.Render()

				window.SwapBuffers()
				glfw.PollEvents()
			}
		}(ctxLoop)

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
	if monitor == nil {
		monitor = glfw.GetPrimaryMonitor()
	}

	videoMode := monitor.GetVideoMode()

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

func isActive(ch <-chan struct{}) bool {
	select {
	case _, ok := <-ch:
		return ok
	default:
		return true
	}
}
