// Copyright (c) 2020-2025 by Marko Gaćeša

package types

import (
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/game/setup"
	"github.com/marko-gacesa/udpstar/udpstar"
)

type GameOneParams struct {
	PlayerInCh chan<- []byte
	ActionCh   chan<- action.Action
	Game       core.RenderRequester
	Done       <-chan struct{}
}

type GameDoubleParams struct {
	Player1InCh chan<- []byte
	Player2InCh chan<- []byte
	ActionCh    chan<- action.Action
	Game        core.RenderRequester
	Done        <-chan struct{}
}

type GameParams struct {
	PlayerInCh  [setup.MaxLocalPlayers]chan<- []byte
	FieldCount  byte
	ActionCh    chan<- action.Action
	LatenciesFn func() []udpstar.LatencyActor
	Game        core.RenderRequester
	Done        <-chan struct{}
}
