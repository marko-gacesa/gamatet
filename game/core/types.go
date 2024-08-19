// Copyright (c) 2020 by Marko Gaćeša

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
	RenderRequest(ctx context.Context, fieldIdx int, t time.Time, ch chan<- *field.RenderInfo)
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
	RandomSeed     int // used for piece feed and random events
	FeedBagSize    int
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

func ChPair[T any](ctx context.Context) (in chan<- T, out <-chan T) {
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
