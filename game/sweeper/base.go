// Copyright (c) 2024 by Marko Gaćeša

package sweeper

import (
	"gamatet/game/field"
	"gamatet/game/op"
	"time"
)

func newBase(f *field.Field) *base {
	s := &base{}
	s.field = f
	s.timer = time.NewTimer(time.Second)
	s.timer.Stop()
	return s
}

type base struct {
	field  *field.Field
	timer  *time.Timer
	active bool
}

func (s *base) Timer() <-chan time.Time {
	return s.timer.C
}

func (s *base) Start(op.Analyzer) {
	if s.active {
		// timer is already active or nothing to do
		return
	}

	s.active = true
	s.timer.Reset(time.Microsecond)
}

func (s *base) Pause() {
	if !s.active {
		return
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

	s.timer.Reset(time.Millisecond)
}

func (s *base) endIteration() {
	s.active = false
}

func (s *base) reschedule(d time.Duration) {
	s.timer.Reset(d)
}
