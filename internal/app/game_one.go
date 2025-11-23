// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"github.com/marko-gacesa/channel"
	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/types"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func (app *App) gameSinglePlayer(ctx screen.Context) types.GameOneParams {
	var s setup.Setup
	if app.resultSetup != nil {
		s = *app.resultSetup
	}
	s.Sanitize()

	fieldCh := make(chan []byte)
	playerPipe := channel.MakePipe[[]byte]()
	playerInCh, playerOutCh := playerPipe.In, playerPipe.Out

	actionCh := make(chan action.Action)

	pieceFeed := Feed(s)

	setup := core.Setup{
		Name: "",
		Config: core.GameConfig{
			WidthPerPlayer: int(s.FieldOptions.WidthSingle),
			Height:         int(s.FieldOptions.Height),
			Level:          int(s.FieldOptions.Speed),
			PlayerZones:    false,
			FieldConfig: field.Config{
				PieceCollision: false,
				Anim:           true,
			},
			RandomSeed: int(s.MiscOptions.Seed),
			PieceFeed:  pieceFeed,
			SamePieces: true,
		},
		Fields: []core.FieldSetup{
			{
				OutCh: fieldCh,
				Players: []core.PlayerSetup{
					{
						Name:    app.cfg.LocalPlayers.Infos[0].Name,
						Config:  piece.Config(app.cfg.LocalPlayers.Infos[0].GameConfig),
						IsLocal: true,
						Index:   0,
						InCh:    playerOutCh,
					},
				},
			},
		},
		ActionCh: actionCh,
	}

	g := core.MakeHost(setup, core.HostOptions{})

	// go-routine for processing events for the field
	go func() {
		defer ctx.Stop()

		g.Perform(ctx)
	}()

	// go-routine to consume all field events
	go channel.Drain(fieldCh)

	app.returnToMainScreen()

	return types.GameOneParams{
		PlayerInCh:  playerInCh,
		PlayerInput: app.cfg.LocalPlayers.Infos[0].Input,
		ActionCh:    actionCh,
		Game:        g,
		Done:        ctx.Done(),
	}
}
