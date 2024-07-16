// Copyright (c) 2020 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/event"
	"gamatet/game/op"
	"time"
)

type Sweeper interface {
	// Timer returns a channel. When a value is returned through the channel the Sweep method should be called.
	Timer() <-chan time.Time

	// Start starts the sweeper. The analyzer is used for conditional start.
	Start(analyzer op.Analyzer)

	// Pause pauses the internal timer of the sweeper.
	Pause()

	// Unpause resumes the internal timer of the sweeper.
	Unpause()

	// Sweep removes blocks from the field.
	Sweep(p event.Pusher)
}
