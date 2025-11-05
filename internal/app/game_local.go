// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/game/setup"
	"gamatet/internal/types"
	"gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/channel"
)

func (app *App) gameMultiPlayerLocal(ctx screen.Context) types.GameParams {
	var s setup.Setup
	if app.resultSetup != nil {
		s = *app.resultSetup
	}
	s.Sanitize()

	pieceFeed := Feed(s)

	var (
		playerCount    = s.PlayerCount()
		teamSize       = s.GameOptions.TeamSize
		fieldCount     = s.GameOptions.FieldCount
		zones          = s.GameOptions.PlayerZones
		pieceCollision = s.GameOptions.PieceCollision
		fieldW         = s.FieldOptions.WidthSingle
		fieldH         = s.FieldOptions.Height
		speed          = s.FieldOptions.Speed
		seed           = int(s.MiscOptions.Seed)
	)
	if teamSize > 1 {
		fieldW = s.WidthPerPlayer
	}

	fieldChs := make([]chan []byte, fieldCount)
	for i := range fieldChs {
		fieldChs[i] = make(chan []byte)
	}

	var playerInChs [setup.MaxLocalPlayers]chan<- []byte
	var playerOutChs [setup.MaxLocalPlayers]<-chan []byte
	for i := range playerCount {
		pipe := channel.MakePipe[[]byte]()
		playerInChs[i], playerOutChs[i] = pipe.In, pipe.Out
	}

	fields := make([]core.FieldSetup, fieldCount)
	for i := range fields {
		players := make([]core.PlayerSetup, teamSize)

		for j := range players {
			playerIdx := i*int(teamSize) + j
			players[j] = core.PlayerSetup{
				Name:    app.cfg.LocalPlayers.Infos[playerIdx].Name,
				IsLocal: true,
				Config:  piece.Config(app.cfg.LocalPlayers.Infos[playerIdx].PlayerConfig),
				InCh:    playerOutChs[playerIdx],
			}
		}

		fields[i] = core.FieldSetup{
			OutCh:   fieldChs[i],
			Players: players,
		}
	}

	actionCh := make(chan action.Action)

	setup := core.Setup{
		Name: s.Name,
		Config: core.GameConfig{
			WidthPerPlayer: int(fieldW),
			Height:         int(fieldH),
			Level:          int(speed),
			PlayerZones:    zones,
			FieldConfig: field.Config{
				PieceCollision: pieceCollision,
				Anim:           true,
			},
			RandomSeed: seed,
			PieceFeed:  pieceFeed,
			SamePieces: s.SamePiecesForAll,
		},
		Fields:   fields,
		ActionCh: actionCh,
	}

	g := core.MakeHost(setup, core.HostOptions{})

	// go-routine for processing events for the field
	go func() {
		defer ctx.Stop()

		g.Perform(ctx)
	}()

	// go-routine to consume all field events
	for _, fieldCh := range fieldChs {
		go channel.Drain(fieldCh)
	}

	app.returnToMainScreen()

	return types.GameParams{
		PlayerInCh: playerInChs,
		FieldCount: fieldCount,
		ActionCh:   actionCh,
		Game:       g,
		Done:       ctx.Done(),
	}
}
