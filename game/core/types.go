// Copyright (c) 2020-2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package core

import (
	"context"
	"time"

	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/logic/anim"
)

type Performer interface {
	Perform(ctx context.Context) []PlayerResult
}

type RenderRequester interface {
	// RenderRequest is a method for requesting render info for a field.
	RenderRequest(fieldIdx int, t time.Time, info *field.RenderInfo, chDone chan<- bool)

	// GetSize return size of the field and number of players.
	GetSize(idx int) (int, int, int)

	// AddAnim adds an animation to all fields in the game. Not intended as a gameplay feature, but for the UI effects.
	AddAnim(anim anim.Anim)
}

type Setup struct {
	Name     string
	Config   GameConfig
	Fields   []FieldSetup
	ActionCh <-chan action.Action
}

type GameConfig struct {
	WidthPerPlayer int
	Height         int
	Level          int
	PlayerZones    bool
	FieldConfig    field.Config
	RandomSeed     int // used for random events
	PieceFeed      piece.Feed
	SamePieces     bool
	Shooters       bool
}

type FieldSetup struct {
	// InCh is user by clients: To read events coming from the server.
	InCh <-chan []byte

	// OutCh is used by server: To sent events to clients.
	OutCh chan<- []byte

	Players []PlayerSetup
}

type PlayerSetup struct {
	Name   string
	Config piece.Config

	IsLocal    bool
	LocalIndex int
	Index      int

	ControlsStr string

	// InCh is used for direct player input. Actions are read from the channel.
	InCh <-chan []byte
}

type PiecePlace struct {
	FieldIdx byte
	CtrlIdx  byte
}

type PlayerResult struct {
	FieldIdx      byte
	CtrlIdx       byte
	PlayerIndex   byte
	Outcome       field.Outcome
	BlocksRemoved int
	Score         uint
	PieceCount    uint
	Level         uint
}
