// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"errors"
	"fmt"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/game/piece"
	"gamatet/game/setup"
	"gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udpstar/client"
	"net"
)

// gameUDPClient function is a glue that links the game engine with the network layer.
// It needs to create channels that would link them (because lobby->session leaves channels unassigned).
// * The channel from which the network layer reads local player inputs: session.Actors[i].InputCh
// * [A] The channel from which the game engine will read field events: core.FieldSetup[i].InCh
// * [B] The previous should be linked to where the network layer would put field events: session.Stories[i].Channel
// * [C] The channel from which the network engine reads local player actions: session.Actor[i].Channel
//
// Server:                                                         |  Client:
// +-Game Engine: Server ---------+  +- Network Engine: Server -+  |  +- Network Engine: Client -----+  +-Game Engine: Client ---------+
// | Field                        |  |                          |  |  |                              |  | Field                        |
// |    InCh {must be nil}        |  |  Story                   |  |  |  Story              /->------|--|--> InCh [A]                  |
// |    OutCh [1] ->--------------|--|----> [2] Channel ->------|--|--|---> Channel [B] ->-/         |  |    OutCh {must be nil}       |
// |                              |  |                          |  |  |                              |  |                              |
// | Local player                 |  |  Client Actor            |  |  |  Local Actor                 |  | Local player                 |
// |    InCh [3] <-- keyboard     |  |     Channel <------------|--|--|--<- InputCh [C] <-- keyboard |  |    InCh {must be nil}        |
// |    OutCh {must be nil}       |  |      [5]|                |  |  |                              |  |    OutCh {must be nil}       |
// |                              |  +--------------------------+  |  +------------------------------+  |                              |
// | Remote player                |            |                   |                                    | Remote player                |
// |    InCh [4] <----------------|----------<-+                   |                                    |    InCh {must be nil}        |
// |    OutCh {must be nil}       |                                |                                    |    OutCh {must be nil}       |
// +------------------------------+                                |                                    +------------------------------+
func (app *App) gameUDPClient(ctx screen.Context) core.GameParams {
	session := app.resultClientSession
	serverAddr := app.resultServerAddress

	app.returnToMainScreen()

	gameParams, err := app._gameUDPClient(ctx, session, serverAddr)
	if err != nil {
		panic(err)
	}

	return gameParams
}

func (app *App) _gameUDPClient(ctx screen.Context, session *client.Session, serverAddr net.UDPAddr) (core.GameParams, error) {
	if session == nil {
		return core.GameParams{}, errors.New("input is missing")
	}

	var s setup.Setup

	if err := setup.Unpack(&s, session.Def); err != nil {
		return core.GameParams{}, fmt.Errorf("unable to unpack setup: %s", err)
	}

	if s.Sanitize() {
		return core.GameParams{}, errors.New("sanitize is required")
	}

	if int(s.GameOptions.FieldCount) != len(session.Stories) {
		return core.GameParams{}, fmt.Errorf("mismatch: field count=%d, story count=%d",
			s.GameOptions.FieldCount, len(session.Stories))
	}

	// Input channels for local players. Closed on the UI component. Elements can be nil.
	var playerInChs [setup.MaxLocalPlayers]chan<- []byte

	// Field pipes. Closed when the game engine completes.
	fieldPipes := make([]core.ChannelPipe[[]byte], s.GameOptions.FieldCount)

	fields := make([]core.FieldSetup, len(session.Stories))
	for fieldIdx := range fields {
		storyToken := session.Stories[fieldIdx].StoryInfo.Token

		actors, err := session.StoryActors(storyToken)
		if err != nil {
			return core.GameParams{}, fmt.Errorf("unable to get actors for story %x: %s", storyToken, err)
		}

		if int(s.GameOptions.TeamSize) != len(actors) {
			return core.GameParams{}, fmt.Errorf("mismatch: team size=%d, actor count=%d",
				s.GameOptions.TeamSize, len(actors))
		}

		fieldPlayers := make([]core.PlayerSetup, len(actors))
		for storyActorIdx, actor := range actors {
			if actor.Token == 0 {
				fieldPlayers[actor.ActorIdx] = core.PlayerSetup{
					Name:   actor.Name,
					Config: piece.Config{},
				}
				continue
			}

			localPlayerInfo, localPlayerIdx := app.LocalPlayer(actor.Token)
			if localPlayerIdx < 0 {
				return core.GameParams{}, fmt.Errorf("local player token=%x not found", actor.Token)
			}

			playerInputPipe := core.MakeChannelPipe[[]byte](ctx) // The "In" part of the pipe should be closed on UI component.
			playerInChs[localPlayerIdx] = playerInputPipe.In

			session.Actors[storyActorIdx].InputCh = playerInputPipe.Out // [C] The network layer reads player inputs from here.

			fieldPlayers[storyActorIdx] = core.PlayerSetup{
				Name:   localPlayerInfo.Name,
				Config: piece.Config(localPlayerInfo.PlayerConfig),
				InCh:   nil,
			}
		}

		fieldPipes[fieldIdx] = core.MakeChannelPipe[[]byte](ctx) // The "In" part of the pipe is closed when the game engine stops.

		session.Stories[fieldIdx].Channel = fieldPipes[fieldIdx].In // [B] The network engine puts field events to this channel.

		fields[fieldIdx] = core.FieldSetup{
			InCh:    fieldPipes[fieldIdx].Out, // [A] The game engine reads game events from this channel.
			Players: fieldPlayers,
		}
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
			PieceFeed:  piece.NewRotTetrominoFeed(4, seed),
		},
		Fields: fields,
	})

	udpSender := udpSender{
		addr: serverAddr,
		srv:  app.udpService,
	}

	cli, err := client.New(udpSender, *session, client.WithLogger(app.logger))
	if err != nil {
		return core.GameParams{}, fmt.Errorf("unable to create client: %s", err)
	}

	go func() {
		cli.Start(ctx)
	}()

	err = app.udpService.Handle(ctx, func(data []byte, addr net.UDPAddr) []byte {
		cli.HandleIncomingMessages(data)
		return nil
	})
	if err != nil {
		return core.GameParams{}, fmt.Errorf("failed to start udp packet handler: %s", err)
	}

	// Go-routine for processing events for the field
	go func() {
		defer func() {
			for _, fieldPipe := range fieldPipes {
				close(fieldPipe.In)
			}
		}()
		defer ctx.Stop()

		gameInterpreter.Perform(ctx)
	}()

	return core.GameParams{
		PlayerInCh: playerInChs,
		FieldCount: byte(len(fields)),
		Game:       gameInterpreter,
		Done:       ctx.Done(),
	}, nil
}
