// Copyright (c) 2023-2025 by Marko Gaćeša

package loop

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/marko-gacesa/gamatet/graphics/scene"
	"github.com/marko-gacesa/gamatet/internal/app"
	"github.com/marko-gacesa/gamatet/internal/values"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func Loop(globalCtx context.Context, app *app.App) error {
	log := app.Log()

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

	videoConfig := app.VideoConfig()
	if err := func() (err error) {
		if videoConfig.Fullscreen {
			log.Info("Starting fullscreen")
			window, err = windowFullscreen(values.ProgramName, nil)
		} else {
			log.Info("Creating window", "width", videoConfig.WindowWidth, "height", videoConfig.WindowHeight)
			window, err = windowResizable(videoConfig.WindowWidth, videoConfig.WindowHeight, values.ProgramName)
			window.SetOpacity(videoConfig.WindowOpacity)
		}
		return
	}(); err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}

	defer window.Destroy()

	log.Info("Setting current OpenGL context")

	window.MakeContextCurrent()

	log.Info("Initializing OpenGL")

	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize OpenGL bindings: %w", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Info("OpenGL", "version", version)

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

				log.Info("Screen done",
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
