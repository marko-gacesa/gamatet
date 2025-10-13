// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/game/setup"
	"gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/channel"
)

func (app *App) game(ctx screen.Context) core.GameParams {
	if app.resultSetup == nil {
		panic("Input missing")
	}

	s := app.resultSetup

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
		pipe := core.MakeChannelPipe[[]byte](ctx)
		playerInChs[i], playerOutChs[i] = pipe.In, pipe.Out
	}

	fields := make([]core.FieldSetup, fieldCount)
	for i := range fields {
		players := make([]core.PlayerSetup, teamSize)

		for j := range players {
			playerIdx := i*int(teamSize) + j
			players[j] = core.PlayerSetup{
				Name:   app.cfg.PlayerInfos[playerIdx].Name,
				Config: piece.Config(app.cfg.PlayerInfos[playerIdx].PlayerConfig),
				InCh:   playerOutChs[playerIdx],
			}
		}

		fields[i] = core.FieldSetup{
			OutCh:   fieldChs[i],
			Players: players,
		}
	}

	setup := core.Setup{
		Name: app.resultSetup.Name,
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
			PieceFeed:  piece.NewRotTetrominoFeed(4, seed),
		},
		Fields: fields,
	}

	g := core.MakeHost(setup)

	// go-routine for processing events for the field
	go func() {
		defer func() {
			for _, fieldCh := range fieldChs {
				close(fieldCh)
			}
		}()
		defer ctx.Stop()

		g.Perform(ctx)
	}()

	// go-routine to consume all field events
	for _, fieldCh := range fieldChs {
		go channel.Drain(fieldCh)
	}

	app.returnToMainScreen()

	return core.GameParams{
		PlayerInCh: playerInChs,
		FieldCount: fieldCount,
		Game:       g,
		Done:       ctx.Done(),
	}
}
