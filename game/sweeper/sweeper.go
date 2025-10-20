// Copyright (c) 2020, 2025 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/event"
	"time"
)

type Sweeper interface {
	// Timer returns a channel. When a value is returned through the channel the Sweep method should be called.
	Timer() <-chan time.Time

	// Start starts the sweeper. The analyzer is used for conditional start.
	// Returns true if the sweeper has just been started, or false if it was already active.
	Start(analyzer *Analyzer) bool

	// Pause pauses the internal timer of the sweeper.
	Pause()

	// Unpause resumes the internal timer of the sweeper.
	Unpause()

	// Sweep removes blocks from the field.
	Sweep(p event.Pusher)
}
