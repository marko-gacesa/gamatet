// Copyright (c) 2024, 2025 by Marko Gaćeša

package screen

import (
	"context"
	"time"
)

// Screen abstracts screen rendering.
type Screen interface {
	// UpdateViewSize should be called whenever the viewport size has changed.
	UpdateViewSize(w, h int)

	// Release should be called to release any allocated resources.
	Release()

	// InputKeyPress handles keyboard key press event
	InputKeyPress(key, scancode int)

	// InputChar handles keyboard input.
	InputChar(char rune)

	// Prepare should be called prior to the Render and can be used asynchronously prepare render data.
	Prepare(now time.Time)

	// Render presents data onto the screen.
	Render()
}

type Screener interface {
	Screen(ctx context.Context, data any) Screen
}
