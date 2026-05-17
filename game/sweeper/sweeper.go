// Copyright (c) 2020, 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"time"

	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
)

type Sweeper interface {
	// Timer returns a channel. When a value is returned through the channel the Sweep method should be called.
	Timer() <-chan time.Time

	// Analyze inspects a batch of events and conditionally starts the sweeper.
	Analyze(events event.Reader)

	// Pause pauses the internal timer of the sweeper.
	Pause()

	// Unpause resumes the internal timer of the sweeper.
	Unpause()

	// Sweep removes blocks from the field.
	Sweep(p event.Pusher)
}

type FieldPusher struct {
	Field  field.Reader
	Pusher event.Pusher
}
