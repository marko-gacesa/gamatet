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
	"github.com/marko-gacesa/udpstar/udpstar/message"
	"github.com/marko-gacesa/udpstar/udpstar/server"
	"net"
)

// gameUDPServer function is a glue that links the game engine with the network layer.
// It needs to create channels that would link them (because lobby->session leaves channels unassigned).
// * [1] The channel on which the game engine would put game's field events: core.FieldSetup.Out
// * [2] The previous should be linked to where the network layer would read field events: session.Stories[i].Channel
// * [3] The channel from which game engine reads local player actions: core.FieldSetup[i].Players[j].InCh
// * [4] The channel from which game engine reads remote player actions: core.FieldSetup[i].Players[j].InCh
// * [5] The previous should be linked to the channel on which network puts remote actor's action: session.Clients[actor.ClientIdx].Actors[actorIdx].Channel
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
func (app *App) gameUDPServer(ctx screen.Context) core.GameParams {
	session := app.resultServerSession
	clientMap := app.resultClientMap

	app.returnToMainScreen()

	gameParams, err := app._gameUDPServer(ctx, session, clientMap)
	if err != nil {
		panic(err)
	}

	return gameParams
}

func (app *App) _gameUDPServer(ctx screen.Context, session *server.Session, clientMap map[message.Token]server.ClientData) (core.GameParams, error) {
	if session == nil || clientMap == nil {
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

	// Input channel pipes for remote players. Closed when the ctx closes.
	playerRemoteInputPipeMap := map[message.Token]core.ChannelPipe[[]byte]{}

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

		teamSize := len(actors)

		fieldPlayers := make([]core.PlayerSetup, teamSize)
		for storyActorIdx, actor := range actors {
			if !actor.IsLocal {
				var playerConfig setup.PlayerConfig
				if err := setup.Unpack(&playerConfig, actor.Config); err != nil {
					return core.GameParams{},
						fmt.Errorf("unable to unpack player config for actor %x for client %x: %s", actor.Token, storyToken, err)
				}

				pipe := core.MakeChannelPipe[[]byte](ctx) // The "In" part of the pipe is closed in this function.
				playerRemoteInputPipeMap[actor.Token] = pipe

				session.Clients[actor.ClientIdx].Actors[actor.ClientActorIdx].Channel = pipe.In // [5] The network layer accepts remote player inputs here.

				fieldPlayers[storyActorIdx] = core.PlayerSetup{
					Name:   actor.Name,
					Config: piece.Config(playerConfig),
					InCh:   pipe.Out, // [4] The game engine reads remote actors actions from here.
				}

				continue
			}

			localPlayerInfo, localPlayerIdx := app.LocalPlayer(actor.Token)
			if localPlayerIdx < 0 {
				return core.GameParams{}, fmt.Errorf("local player token=%x not found", actor.Token)
			}

			playerInputPipe := core.MakeChannelPipe[[]byte](ctx) // The "In" part of the pipe should be closed on UI component.
			playerInChs[localPlayerIdx] = playerInputPipe.In

			fieldPlayers[storyActorIdx] = core.PlayerSetup{
				Name:   localPlayerInfo.Name,
				Config: piece.Config(localPlayerInfo.PlayerConfig),
				InCh:   playerInputPipe.Out, // [3] The game engine reads local player actions from here (directly from the input device - keyboard).
			}
		}

		fieldPipes[fieldIdx] = core.MakeChannelPipe[[]byte](ctx)     // The "In" part of the pipe is closed when game host finishes.
		session.Stories[fieldIdx].Channel = fieldPipes[fieldIdx].Out // [2] The network layer reads events from here.

		fields[fieldIdx] = core.FieldSetup{
			OutCh:   fieldPipes[fieldIdx].In, // [1] The game engine puts field events here.
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

	gameHost := core.MakeHost(core.Setup{
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
			PieceFeed:  piece.NewTetrominoFeed(4, seed),
		},
		Fields: fields,
	})

	// Start the network engine component.
	if err := app.gameServer.StartSession(ctx, session, clientMap, gameHost); err != nil {
		return core.GameParams{}, fmt.Errorf("unable to start session: %s", err)
	}

	if err := app.udpService.Handle(ctx, func(data []byte, addr net.UDPAddr) []byte {
		return app.gameServer.HandleIncomingMessages(data, addr)
	}); err != nil {
		return core.GameParams{}, fmt.Errorf("unable to handle udp message: %s", err)
	}

	// Go-routine for closing remote player input pipes.
	go func() {
		<-ctx.Done()

		for _, pipe := range playerRemoteInputPipeMap {
			close(pipe.In)
		}
	}()

	// Go-routine for processing events for the field.
	go func() {
		defer func() {
			for _, fieldPipe := range fieldPipes {
				close(fieldPipe.In)
			}
		}()
		defer ctx.Stop()

		gameHost.Perform(ctx)
	}()

	return core.GameParams{
		PlayerInCh: playerInChs,
		FieldCount: byte(len(fields)),
		Game:       gameHost,
		Done:       ctx.Done(),
	}, nil
}
