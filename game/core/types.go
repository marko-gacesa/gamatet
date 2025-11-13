// Copyright (c) 2020-2025 by Marko Gaćeša

package core

import (
	"context"
	"time"

	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
)

type Performer interface {
	Perform(ctx context.Context)
}

type RenderRequester interface {
	// RenderRequest is a method for requesting render info for a field. Once ready, it will be placed on the channel.
	// When rendering is completed the RenderInfo should be returned with a call to field.ReturnRenderInfo(renderInfo).
	RenderRequest(fieldIdx int, t time.Time, ch chan<- *field.RenderInfo)

	// GetSize return size of the field and number of players.
	GetSize(idx int) (int, int, int)
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

	IsLocal bool
	Index   int

	// InCh is used for direct player input. Actions are read from the channel.
	InCh <-chan []byte
}
