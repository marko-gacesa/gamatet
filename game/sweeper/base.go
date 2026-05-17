// Copyright (c) 2024, 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package sweeper

import (
	"time"

	"github.com/marko-gacesa/gamatet/game/field"
)

func newBase(f field.Reader) *base {
	s := &base{}
	s.field = f
	s.timer = time.NewTimer(time.Second)
	s.timer.Stop()
	return s
}

type base struct {
	field     field.Reader
	timer     *time.Timer
	startedAt time.Time     // When the timer was started
	remaining time.Duration // Duration remaining if paused
	active    bool
}

func (s *base) Timer() <-chan time.Time {
	return s.timer.C
}

func (s *base) Pause() {
	if !s.active {
		return
	}

	s.remaining = time.Since(s.startedAt)
	if s.remaining <= 0 {
		s.remaining = time.Nanosecond
	}
	s.timer.Stop()
	select {
	default:
	case <-s.timer.C:
	}
}

func (s *base) Unpause() {
	if !s.active {
		return
	}

	s.reschedule(s.remaining)
}

func (s *base) start() bool {
	result := !s.active
	s.active = true
	s.reschedule(time.Microsecond)
	return result
}

// endIteration should be called during Sweep to stop the current iteration.
// Basically, during Sweep one of the two methods should be called: endIteration or reschedule.
func (s *base) endIteration() {
	s.active = false
}

// reschedule can be called during Sweep to trigger another iteration of Sweep.
// Basically, during Sweep one of the two methods should be called: endIteration or reschedule.
func (s *base) reschedule(d time.Duration) {
	s.startedAt = time.Now()
	s.timer.Reset(d)
}
