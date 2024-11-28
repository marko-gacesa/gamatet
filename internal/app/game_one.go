// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
)

func (app *App) gameOne(ctx context.Context) (core.GameOneParams, context.Context) {
	const fieldW = 10
	const fieldH = 24

	var level = 0
	const seed = 101

	ctx, cancelCtx := context.WithCancel(ctx)

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
						Name: "marko",
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
	go func(ctx context.Context) {
		defer close(fieldCh)
		defer cancelCtx()

		g.Perform(ctx)
	}(ctx)

	// go-routine to consume all field events
	go func() {
		for range fieldCh {
		}
	}()

	return core.GameOneParams{
		PlayerInCh: playerInCh,
		Game:       g,
		Done:       ctx.Done(),
	}, ctx
}
