// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"gamatet/game/action"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/internal/types"
	"gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/channel"
)

func (app *App) gameDouble(ctx screen.Context) types.GameDoubleParams {
	const fieldW = 8
	const fieldH = 24

	var level = 0
	const seed = 101

	fieldCh := make(chan []byte)

	player1Pipe := channel.MakePipe[[]byte]()
	player2Pipe := channel.MakePipe[[]byte]()
	player1InCh, player1OutCh := player1Pipe.In, player1Pipe.Out
	player2InCh, player2OutCh := player2Pipe.In, player2Pipe.Out

	actionCh := make(chan action.Action)

	setup := core.Setup{
		Name: "",
		Config: core.GameConfig{
			WidthPerPlayer: fieldW,
			Height:         fieldH,
			Level:          level,
			PlayerZones:    true,
			FieldConfig: field.Config{
				PieceCollision: false,
				Anim:           true,
			},
			RandomSeed: seed,
			PieceFeed:  piece.NewRotTetrominoFeed(4, seed),
		},
		Fields: []core.FieldSetup{
			{
				OutCh: fieldCh,
				Players: []core.PlayerSetup{
					{
						Name: "marko",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideDisabled:       false,
							WallKick:            2,
						},
						IsLocal: true,
						InCh:    player1OutCh,
					},
					{
						Name: "ogi",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideDisabled:       false,
							WallKick:            2,
						},
						IsLocal: true,
						InCh:    player2OutCh,
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

	return types.GameDoubleParams{
		Player1InCh: player1InCh,
		Player2InCh: player2InCh,
		ActionCh:    actionCh,
		Game:        g,
		Done:        ctx.Done(),
	}
}
