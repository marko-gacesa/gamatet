// Copyright (c) 2024 by Marko Gaćeša

package screen

import (
	"context"
	"time"
)

// Screen abstracts screen rendering.
type Screen interface {
	// UpdateView should be called whenever the viewport size has changed.
	UpdateView(w, h int)

	// Done returns a channel that is closed when the screen finished.
	// If the screen finishes abnormally the error would be placed to the channel.
	Done() <-chan error

	// Release should be called to release any allocated resources.
	Release()

	// InputKeyPress handles keyboard key press event
	InputKeyPress(key, scancode int)

	// InputChar handles keyboard input.
	InputChar(char rune)

	// Prepare should be called prior to the Render and can be used asynchronously prepare render data.
	Prepare(ctx context.Context, now time.Time)

	// Render presents data onto the screen.
	Render(ctx context.Context)
}

type Screener interface {
	Screen(ctx context.Context, data any) Screen
}
