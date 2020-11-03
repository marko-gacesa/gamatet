// Copyright (c) 2020 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/event"
	"gamatet/game/field"
	"time"
)

type Sweeper interface {
	// Field returns the field on which the sweeper is working.
	Field() *field.Field

	// Timer returns a channel. When a value is returned through the channel the Sweep method should be called.
	Timer() <-chan time.Time

	// Start starts the sweeper. It should be called whenever the field is changed in such way that there is
	// a possibility that the sweeper may have a job to do (i.e. when blocks are added, but not when removed).
	Start()

	// Pause pauses the internal timer of the sweeper.
	Pause()

	// Unpause resumes the internal timer of the sweeper.
	Unpause()

	// Sweep removes blocks from the field.
	Sweep(p event.Pusher)
}
