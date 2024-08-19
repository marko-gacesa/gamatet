// Copyright (c) 2024 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"sync"
)

func (app *App) singleSimple(ctx context.Context) (*core.GameHost, chan<- []byte, *sync.WaitGroup) {
	const fieldW = 10
	const fieldH = 24

	var level = 0
	const seed = 101

	fieldCh := make(chan []byte)
	playerInCh, playerOutCh := core.ChPair[[]byte](ctx)

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
			RandomSeed:  seed,
			FeedBagSize: 2,
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

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		defer close(fieldCh)

		g.Perform(ctx)
	}(ctx)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for range fieldCh {
		}
	}()

	return g, playerInCh, wg
}
