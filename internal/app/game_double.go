// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"github.com/marko-gacesa/udpstar/channel"
)

func (app *App) gameDouble(ctx context.Context) (core.GameDoubleParams, context.Context) {
	const fieldW = 8
	const fieldH = 24

	var level = 0
	const seed = 101

	ctx, cancelCtx := context.WithCancel(ctx)

	fieldCh := make(chan []byte)
	player1InCh, player1OutCh := core.ChannelPipe[[]byte](ctx)
	player2InCh, player2OutCh := core.ChannelPipe[[]byte](ctx)

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
						InCh: player1OutCh,
					},
					{
						Name: "ogi",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideEnabled:        true,
							MaxWallKick:         2,
						},
						InCh: player2OutCh,
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
	go channel.Drain(fieldCh)

	return core.GameDoubleParams{
		Player1InCh: player1InCh,
		Player2InCh: player2InCh,
		Game:        g,
		Done:        ctx.Done(),
	}, ctx
}
