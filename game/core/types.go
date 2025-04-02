// Copyright (c) 2020-2025 by Marko Gaćeša

package core

import (
	"context"
	"gamatet/game/field"
	"gamatet/game/piece"
	"time"
)

type Performer interface {
	Perform(ctx context.Context)
}

type RenderRequester interface {
	// RenderRequest is a method for requesting render info for a field. Once ready, it will be placed on the channel.
	// When rendering is completed the RenderInfo should be returned with a call to field.ReturnRenderInfo(renderInfo).
	RenderRequest(fieldIdx int, t time.Time, ch chan<- *field.RenderInfo)

	// GetSize return size of the field
	GetSize(idx int) (int, int, []piece.DisplayPosition)
}

type Setup struct {
	Name   string
	Config GameConfig
	Fields []FieldSetup
}

type GameConfig struct {
	WidthPerPlayer int
	Height         int
	Level          int
	PlayerZones    bool
	FieldConfig    field.Config
	RandomSeed     int // used for random events
	PieceFeed      piece.Feed
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

	// InCh is used for direct player input. Actions are read from the channel.
	InCh <-chan []byte

	// OutCh is used for clients to send actions to the server. Servers don't use this.
	OutCh chan<- []byte
}

type GameOneParams struct {
	PlayerInCh chan<- []byte
	Game       RenderRequester
	Done       <-chan struct{}
}

type GameDoubleParams struct {
	Player1InCh chan<- []byte
	Player2InCh chan<- []byte
	Game        RenderRequester
	Done        <-chan struct{}
}

func ChannelPipe[T any](ctx context.Context) (in chan<- T, out <-chan T) {
	chIn := make(chan T)
	chOut := make(chan T)
	go func() {
		defer close(chOut)

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-chIn:
				chOut <- data
			}
		}
	}()

	return chIn, chOut
}
