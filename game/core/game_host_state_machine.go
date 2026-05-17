// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
)

type hostState byte

const (
	hostStateGetReady hostState = iota
	hostStateNormal
	hostStatePause
	hostStateSuspended
	hostStateFinish
)

func (g *GameHost) stateTransitionFinish() {
	g.state = hostStateFinish
}

func (g *GameHost) stateTransitionGetReady() {
	switch g.state {
	case hostStateNormal:
		g.state = hostStateGetReady
		g._pauseAllFields()
		g._publishNewMode(field.ModeGetReady)
	case hostStateGetReady, hostStatePause, hostStateSuspended, hostStateFinish:
	}
}

func (g *GameHost) stateTransitionPlay() {
	switch g.state {
	case hostStateGetReady:
		g.state = hostStateNormal
		g._unpauseAllFields()
		g._publishNewMode(field.ModeNormal)
	case hostStateNormal, hostStatePause, hostStateSuspended, hostStateFinish:
	}
}

func (g *GameHost) stateTransitionSuspend() {
	switch g.state {
	case hostStateGetReady, hostStateNormal:
		g.state = hostStateSuspended
		g._pauseAllFields()
		g._publishNewMode(field.ModeSuspended)
	case hostStatePause:
		g.state = hostStateSuspended
		g._publishNewMode(field.ModeSuspended)
	case hostStateSuspended, hostStateFinish:
	}
}

func (g *GameHost) stateTransitionUnsuspend() {
	switch g.state {
	case hostStateGetReady, hostStateNormal, hostStatePause, hostStateFinish:
	case hostStateSuspended:
		g.state = hostStatePause
		g._publishNewMode(field.ModePause)
	}
}

func (g *GameHost) stateTransitionPauseToggle() {
	switch g.state {
	case hostStateGetReady:
		g.state = hostStatePause
		g._publishNewMode(field.ModePause)
	case hostStateNormal:
		g.state = hostStatePause
		g._pauseAllFields()
		g._publishNewMode(field.ModePause)
	case hostStatePause:
		g.state = hostStateNormal
		g._unpauseAllFields()
		g._publishNewMode(field.ModeNormal)
	case hostStateSuspended, hostStateFinish:
	}
}

func (g *GameHost) _publishNewMode(newMode field.Mode) {
	for fIdx := range g.fields {
		f := g.fields[fIdx].Field
		g.fields[fIdx].events.Push(op.NewFieldMode(f, newMode, false))
	}
}

func (g *GameHost) _pauseAllFields() {
	for fIdx := range g.fields {
		g._pauseField(fIdx)
	}
}

func (g *GameHost) _unpauseAllFields() {
	for fIdx := range g.fields {
		g._unpauseField(fIdx)
	}
}

func (g *GameHost) _pauseField(fIdx int) {
	if g.fields[fIdx].Field.IsFinished() {
		return
	}
	g.fields[fIdx].Field.Pause()
	for _, s := range g.fields[fIdx].Sweepers {
		s.Pause()
	}
}

func (g *GameHost) _unpauseField(fIdx int) {
	if g.fields[fIdx].Field.IsFinished() {
		return
	}
	g.fields[fIdx].Field.Unpause()
	for _, s := range g.fields[fIdx].Sweepers {
		s.Unpause()
	}
}
