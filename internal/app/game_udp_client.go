// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"errors"
	"fmt"
	"net"

	"github.com/marko-gacesa/channel"
	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/core"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/types"
	"github.com/marko-gacesa/gamatet/logic/latency"
	"github.com/marko-gacesa/gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udpstar"
	"github.com/marko-gacesa/udpstar/udpstar/client"
)

// gameUDPClient function is a glue that links the game engine with the network layer.
// It needs to create channels that would link them (because lobby->session leaves channels unassigned).
// * The channel from which the network layer reads local player inputs: session.Actors[i].InputCh
// * [A] The channel from which the game engine will read field events: core.FieldSetup[i].InCh
// * [B] The previous should be linked to where the network layer would put field events: session.Stories[i].Channel
// * [C] The channel from which the network engine reads local player actions: session.Actor[i].Channel
//
// Server:                                                     |  Client:
// +-Game Engine: Server -----+  +- Network Engine: Server -+  |  +- Network Engine: Client -----+  +-Game Engine: Client ---+
// | Field                    |  |                          |  |  |                              |  | Field                  |
// |    InCh {must be nil}    |  |  Story                   |  |  |  Story              /->------|--|--> InCh [A]            |
// |    OutCh [1] ->----------|--|----> [2] Channel ->------|--|--|---> Channel [B] ->-/         |  |    OutCh {must be nil} |
// |                          |  |                          |  |  |                              |  |                        |
// | Local player             |  |  Client Actor            |  |  |  Local Actor                 |  | Local player           |
// |    InCh [3] <-- keyboard |  |     Channel <------------|--|--|--<- InputCh [C] <-- keyboard |  |    InCh {must be nil}  |
// |                          |  |      [5]|                |  |  |                              |  |                        |
// | Remote player            |  +---------|----------------+  |  +------------------------------+  | Remote player          |
// |    InCh [4] <------------|----------<-+                   |                                    |    InCh {must be nil}  |
// +--------------------------+                                |                                    +------------------------+
func (app *App) gameUDPClient(ctx screen.Context) types.GameParams {
	session := app.resultClientSession
	serverAddr := app.resultServerAddress

	app.returnToMainScreen()

	gameParams, err := app._gameUDPClient(ctx, session, serverAddr)
	if err != nil {
		panic(err)
	}

	return gameParams
}

func (app *App) _gameUDPClient(ctx screen.Context, session *client.Session, serverAddr net.UDPAddr) (types.GameParams, error) {
	if session == nil {
		return types.GameParams{}, errors.New("input is missing")
	}

	var s setup.Setup

	if err := setup.Unpack(&s, session.Def); err != nil {
		return types.GameParams{}, fmt.Errorf("unable to unpack setup: %s", err)
	}

	if s.SanitizeMulti() {
		return types.GameParams{}, errors.New("sanitize is required")
	}

	if int(s.GameOptions.FieldCount) != len(session.Stories) {
		return types.GameParams{}, fmt.Errorf("mismatch: field count=%d, story count=%d",
			s.GameOptions.FieldCount, len(session.Stories))
	}

	pieceFeed := Feed(s)

	fieldHasLocalPlayers := make(map[int]struct{})

	// Input channels for local players. Closed on the UI component. Elements can be nil.
	var playerInChs [setup.MaxLocalPlayers]chan<- []byte

	var playerIndex int

	fields := make([]core.FieldSetup, len(session.Stories))
	for fieldIdx := range fields {
		storyToken := session.Stories[fieldIdx].StoryInfo.Token

		actors, err := session.StoryActors(storyToken)
		if err != nil {
			return types.GameParams{}, fmt.Errorf("unable to get actors for story %x: %s", storyToken, err)
		}

		if int(s.GameOptions.TeamSize) != len(actors) {
			return types.GameParams{}, fmt.Errorf("mismatch: team size=%d, actor count=%d",
				s.GameOptions.TeamSize, len(actors))
		}

		fieldPlayers := make([]core.PlayerSetup, len(actors))
		for storyActorIdx, actor := range actors {
			if actor.Token == 0 {
				fieldPlayers[storyActorIdx] = core.PlayerSetup{
					Name:    playerName(actor.Name, fieldIdx, storyActorIdx, playerIndex),
					Config:  piece.Config{},
					IsLocal: false,
					Index:   playerIndex,
				}
				playerIndex++
				continue
			}

			fieldHasLocalPlayers[fieldIdx] = struct{}{}

			localPlayerInfo, localPlayerIdx := app.LocalPlayer(actor.Token)
			if localPlayerIdx < 0 {
				return types.GameParams{}, fmt.Errorf("local player token=%x not found", actor.Token)
			}

			playerInputPipe := channel.MakePipe[[]byte]() // The "In" part of the pipe should be closed on UI component.
			playerInChs[localPlayerIdx] = playerInputPipe.In

			session.Actors[actor.ActorIdx].InputCh = playerInputPipe.Out // [C] The network layer reads player inputs from here.

			fieldPlayers[storyActorIdx] = core.PlayerSetup{
				Name:    playerName(localPlayerInfo.Name, fieldIdx, storyActorIdx, playerIndex),
				Config:  piece.Config(localPlayerInfo.GameConfig),
				IsLocal: true,
				Index:   playerIndex,
			}
			playerIndex++
		}

		fieldPipe := channel.MakePipe[[]byte]() // The "In" part of the pipe is closed by the network layer, in the client.

		session.Stories[fieldIdx].Channel = fieldPipe.In // [B] The network engine puts field events to this channel.

		fields[fieldIdx] = core.FieldSetup{
			InCh:    fieldPipe.Out, // [A] The game engine reads game events from this channel.
			Players: fieldPlayers,
		}
	}

	udpSender := udpSender{
		addr: serverAddr,
		srv:  app.udpService,
	}

	cli, err := client.New(udpSender, *session, client.WithLogger(app.logger))
	if err != nil {
		return types.GameParams{}, fmt.Errorf("unable to create client: %s", err)
	}

	go func() {
		cli.Start(ctx)
	}()

	err = app.udpService.Handle(ctx, func(data []byte, addr net.UDPAddr) []byte {
		cli.HandleIncomingMessages(data)
		return nil
	})
	if err != nil {
		return types.GameParams{}, fmt.Errorf("failed to start udp packet handler: %s", err)
	}

	var (
		zones          = s.GameOptions.PlayerZones
		pieceCollision = s.GameOptions.PieceCollision
		fieldW         = s.FieldOptions.WidthSingle
		fieldH         = s.FieldOptions.Height
		speed          = s.FieldOptions.Speed
		seed           = int(s.MiscOptions.Seed)
	)

	if playerCount := s.PlayerCount(); playerCount > 1 {
		fieldW = s.WidthPerPlayer
	}

	actionCh := make(chan action.Action)

	var localPlayerCh chan<- []byte
	for i := range playerInChs {
		if playerInChs[i] != nil {
			localPlayerCh = playerInChs[i]
			break
		}
	}

	playerNames := func() []string {
		names := make([]string, 0)
		for fIdx := range fields {
			for fieldPlayerIdx := range fields[fIdx].Players {
				names = append(names, fields[fIdx].Players[fieldPlayerIdx].Name)
			}
		}
		return names
	}()

	latencies := latency.NewList(func() []udpstar.LatencyActor {
		return cli.Latencies().Latencies
	}, func(l []udpstar.LatencyActor) string {
		return latenciesToString(l, playerNames)
	})

	gameInterpreter := core.MakeInterpreter(core.Setup{
		Name: session.Name,
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
			PieceFeed:  pieceFeed,
			SamePieces: s.GameOptions.SamePiecesForAll,
			Shooters:   s.PieceOptions.Shooters,
		},
		Fields:   fields,
		ActionCh: actionCh,
	}, core.InterpreterOptions{
		LocalPlayerActionCh: localPlayerCh,
		SinceLastContactFn:  cli.SinceLastServerMessage,
		Latencies:           latencies,
	})

	// Go-routine for processing events for the field
	go func() {
		defer ctx.Stop()

		gameInterpreter.Perform(ctx)
	}()

	_ = cli.Latency
	_ = cli.Quality

	return types.GameParams{
		PlayerInCh:           playerInChs,
		PlayerInputs:         app.cfg.LocalPlayers.Inputs(),
		FieldHasLocalPlayers: fieldHasLocalPlayers,
		FieldCount:           byte(len(fields)),
		ActionCh:             actionCh,
		Latencies:            latencies,
		Game:                 gameInterpreter,
		Done:                 ctx.Done(),
	}, nil
}
