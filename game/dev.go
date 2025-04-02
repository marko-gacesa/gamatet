// Copyright (c) 2024, 2025 by Marko Gaćeša

package game

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"io"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"
)

func NewFieldTest(
	ctx context.Context,
	fieldW int,
	fieldH int,
	stopFn func(),
) (*core.GameHost, *core.GameInterpreter, chan<- []byte, <-chan struct{}) {
	wgFields := &sync.WaitGroup{}
	wgEvents := &sync.WaitGroup{}
	chStopEvents := make(chan struct{})

	// event transfer between the host and the client

	fieldServerCh := make(chan []byte, 100) // capacity because of the artificial delay to simulate latency
	fieldClientCh := make(chan []byte, 100) // ... otherwise, the capacity isn't needed.

	r := rand.New(rand.NewSource(123))
	_ = r

	wgEvents.Add(1)
	go func() {
		defer wgEvents.Done()

		for e := range fieldServerCh {
			time.Sleep(time.Millisecond * time.Duration(30+r.Intn(100)))
			select {
			case <-chStopEvents:
				return

			case fieldClientCh <- e:
				if len(e) == 0 {
					fmt.Printf("event empty\n")
					continue
				}

				if e[0] == 'Z' {
					n, err := func() (int, error) {
						gzReader, err := gzip.NewReader(bytes.NewReader(e[1:]))
						if err != nil {
							return 0, err
						}

						defer gzReader.Close()

						n, err := io.Copy(io.Discard, gzReader)
						if err != nil {
							return 0, err
						}

						return int(n), nil
					}()
					if err != nil {
						fmt.Printf("events   compressed: gz-len=%-3d; ERROR=%s\n", len(e), err.Error())
						continue
					}

					fmt.Printf("events   compressed: gz-len=%-3d raw-len=%-3d savings=%5.1f%%\n",
						len(e), n, 100*(float64(n)/float64(len(e))))
				} else {
					n := func() int {
						b := bytes.NewBuffer(nil)
						w := gzip.NewWriter(b)
						w.Write(e)
						w.Close()
						return b.Len()
					}()
					fmt.Printf("events uncompressed: gz-len=%-3d raw-len=%-3d savings=%5.1f%%\n",
						n, len(e), 100*(float64(n)/float64(len(e))))
				}
			}
		}

		fmt.Println("GAME: Event transfer stopped")
	}()

	// game setup

	playerInCh, playerOutCh := core.ChannelPipe[[]byte](ctx)

	const seed = 101
	const level = 2

	pieceFeed := piece.NewDebugFeed(seed)

	setupHost := core.Setup{
		Name: "test game",
		Config: core.GameConfig{
			WidthPerPlayer: fieldW,
			Height:         fieldH,
			Level:          level,
			PlayerZones:    true,
			FieldConfig: field.Config{
				PieceCollision: false,
				Anim:           false,
			},
			RandomSeed: seed,
			PieceFeed:  pieceFeed,
		},
		Fields: []core.FieldSetup{
			{
				OutCh: fieldServerCh,
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

	gameHost := core.MakeHost(setupHost)

	setupClient := core.Setup{
		Name: "test game",
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
			PieceFeed:  pieceFeed,
		},
		Fields: []core.FieldSetup{
			{
				InCh: fieldClientCh,
				Players: []core.PlayerSetup{
					{
						Name: "marko",
						Config: piece.Config{
							RotationDirectionCW: false,
							SlideEnabled:        true,
							MaxWallKick:         2,
						},
					},
				},
			},
		},
	}

	gameInterpreter := core.MakeInterpreter(setupClient)

	wgFields.Add(1)
	go func(ctx context.Context) {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Printf("PANIC: %v\n%s\n", r, debug.Stack())
			}
		}()
		defer close(fieldServerCh)
		defer wgFields.Done()

		gameHost.Perform(ctx)

		fmt.Println("GAME: Host stopped")
	}(ctx)

	wgFields.Add(1)
	go func(ctx context.Context) {
		defer func() {
			r := recover()
			if r != nil {
				fmt.Printf("PANIC: %v\n%s\n", r, debug.Stack())
			}
		}()
		defer wgFields.Done()

		gameInterpreter.Perform(ctx)

		fmt.Println("GAME: Interpreter stopped")
	}(ctx)

	waitDoneCh := make(chan struct{})
	go func() {
		wgFields.Wait()
		close(chStopEvents)
		wgEvents.Wait()

		close(waitDoneCh)

		stopFn()

		fmt.Println("GAME: Wait channel closed")
	}()

	return gameHost, gameInterpreter, playerInCh, waitDoneCh
}
