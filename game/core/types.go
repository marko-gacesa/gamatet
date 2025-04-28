// Copyright (c) 2020-2025 by Marko Gaćeša

package core

import (
	"context"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/game/setup"
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

type GameParams struct {
	PlayerInCh [setup.MaxLocalPlayers]chan<- []byte
	FieldCount byte
	Game       RenderRequester
	Done       <-chan struct{}
}

type ChannelPipe[T any] struct {
	In  chan<- T
	Out <-chan T
}

func MakeChannelPipe[T any](ctx context.Context) ChannelPipe[T] {
	chIn := make(chan T)
	chOut := make(chan T)
	go func() {
		defer close(chOut)

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-chIn:
				if !ok {
					return
				}
				chOut <- data
			}
		}
	}()

	return ChannelPipe[T]{
		In:  chIn,
		Out: chOut,
	}
}
