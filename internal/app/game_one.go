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

func (app *App) gameOne(ctx screen.Context) types.GameOneParams {
	const fieldW = 10
	const fieldH = 24

	const level = 8
	const seed = 101

	fieldCh := make(chan []byte)
	playerPipe := channel.MakePipe[[]byte]()
	playerInCh, playerOutCh := playerPipe.In, playerPipe.Out

	actionCh := make(chan action.Action)

	setup := core.Setup{
		Name: "",
		Config: core.GameConfig{
			WidthPerPlayer: fieldW,
			Height:         fieldH,
			Level:          level,
			PlayerZones:    false,
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
						Name: "Marko",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideDisabled:       false,
							WallKick:            piece.WallKickDefault,
						},
						IsLocal: true,
						InCh:    playerOutCh,
					},
				},
			},
		},
		ActionCh: actionCh,
	}

	g := core.MakeHost(setup)

	// go-routine for processing events for the field
	go func() {
		defer ctx.Stop()

		g.Perform(ctx)
	}()

	// go-routine to consume all field events
	go channel.Drain(fieldCh)

	app.returnToMainScreen()

	return types.GameOneParams{
		PlayerInCh: playerInCh,
		ActionCh:   actionCh,
		Game:       g,
		Done:       ctx.Done(),
	}
}
