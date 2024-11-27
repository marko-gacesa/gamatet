// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"sync"
)

func (app *App) gameOne(ctx context.Context) core.GameOneParams {
	const fieldW = 10
	const fieldH = 24

	var level = 0
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

	wg := &sync.WaitGroup{}

	// go-routine for processing events for the field
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		defer close(fieldCh)

		g.Perform(ctx)
	}(ctx)

	// go-routine to consume all field events
	wg.Add(1)
	go func() {
		defer wg.Done()

		for range fieldCh {
		}
	}()

	// go-routine to close the "done" channel when other go-routines finish.
	chDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(chDone)
	}()

	return core.GameOneParams{
		PlayerInCh: playerInCh,
		Game:       g,
		Done:       chDone,
	}
}
