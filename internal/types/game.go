// Copyright (c) 2020-2025 by Marko Gaćeša

package types

import (
	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config/key"
	"github.com/marko-gacesa/gamatet/logic/latency"
)

type GameOneParams struct {
	PlayerInCh  chan<- []byte
	PlayerInput key.Input
	ActionCh    chan<- action.Action
	Game        core.RenderRequester
	Done        <-chan struct{}
}

type GameParams struct {
	PlayerInCh   [setup.MaxLocalPlayers]chan<- []byte
	PlayerInputs [setup.MaxLocalPlayers]key.Input
	FieldCount   byte
	ActionCh     chan<- action.Action
	Latencies    *latency.List
	Game         core.RenderRequester
	Done         <-chan struct{}
}
