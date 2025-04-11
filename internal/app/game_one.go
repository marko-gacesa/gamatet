// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/channel"
)

func (app *App) gameOne(ctx screen.Context) core.GameOneParams {
	const fieldW = 10
	const fieldH = 24

	var level = 7
	const seed = 101

	fieldCh := make(chan []byte)
	playerInCh, playerOutCh := core.ChannelPipe[[]byte](ctx)

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
			PieceFeed:  piece.NewTetrominoFeed(4, seed),
		},
		Fields: []core.FieldSetup{
			{
				OutCh: fieldCh,
				Players: []core.PlayerSetup{
					{
						Name: "",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideEnabled:        true,
							MaxWallKick:         2,
						},
						InCh: playerOutCh,
					},
				},
			},
		},
	}

	g := core.MakeHost(setup)

	// go-routine for processing events for the field
	go func() {
		defer close(fieldCh)
		defer ctx.Stop()

		g.Perform(ctx)
	}()

	// go-routine to consume all field events
	go channel.Drain(fieldCh)

	return core.GameOneParams{
		PlayerInCh: playerInCh,
		Game:       g,
		Done:       ctx.Done(),
	}
}
