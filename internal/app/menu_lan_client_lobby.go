// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"gamatet/game/setup"
	"gamatet/internal/values"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
	"github.com/marko-gacesa/bitdata"
	"github.com/marko-gacesa/udpstar/udpstar/client"
	"github.com/marko-gacesa/udpstar/udpstar/message"
	"net"
	"slices"
	"time"
	"unicode"
)

func (app *App) menuLANClientLobby(ctx screen.Context) *menu.Menu {
	if app.resultClientLobbySelected == nil {
		return app.menuErrorText(ctx, "Input missing")
	}

	app.resultClientSession = nil

	lobbyToken := app.resultClientLobbySelected.Token
	serverAddr := app.resultClientLobbySelected.Addr
	slotCount := byte(len(app.resultClientLobbySelected.Lobby.Slots))

	slotStories := make([]message.Token, slotCount)
	for i := range app.resultClientLobbySelected.Lobby.Slots {
		slotStories[i] = app.resultClientLobbySelected.Lobby.Slots[i].StoryToken
	}

	r := bitdata.NewReaderError(app.resultClientLobbySelected.Lobby.Def)
	var o setup.Setup
	o.Read(r)
	if err := r.Error(); err != nil {
		return app.menuError(ctx, err)
	}

	gameStr := o.String()

	slots := makeLobbyEntries(slotStories, false, app.actorTokens[:])
	blocker := makeStartBlocker()

	// lobby

	udpSender := udpSender{
		addr: serverAddr,
		srv:  app.udpService,
	}

	lobbyClient := client.NewLobby(udpSender, lobbyToken, app.clientToken, client.WithLobbyLogger(app.logger))

	go func() {
		gameSession := lobbyClient.Start(ctx)

		if gameSession != nil {
			app.screenIDNext = routeGameUDPClient

			app.resultServerAddress = serverAddr
			app.resultClientSession = gameSession
		}

		ctx.Stop()
	}()

	var version int

	err := app.udpService.Handle(ctx, func(data []byte, addr net.UDPAddr) []byte {
		lobbyClient.HandleIncomingMessages(data)

		lobby, age := lobbyClient.Get(version)
		if lobby == nil {
			return nil
		}

		if age > time.Minute {
			ctx.Stop()
			return nil
		}

		version = lobby.Version

		slots.setAll(lobby)
		blocker.update(lobby)

		return nil
	})
	if err != nil {
		return app.menuError(ctx, err)
	}

	items := make([]menu.Item, 0, slotCount+2)

	items = append(items, menu.NewStatic(
		gameStr, "", nil,
		menu.WithDisabled(func() bool { return true })))
	for i := byte(0); i < slotCount; i++ {
		items = append(items, menu.NewStatic("", "",
			func(r rune) bool {
				switch unicode.ToLower(r) {
				case '\n':
					name := app.cfg.PlayerInfos[0].Name
					cfg := app.cfg.PlayerInfos[0].PlayerConfig.Serialize()
					lobbyClient.Join(app.actorTokens[0], i, name, cfg)
				case '1', '2', '3', '4':
					idx := byte(r - '1')
					name := app.cfg.PlayerInfos[idx].Name
					cfg := app.cfg.PlayerInfos[idx].PlayerConfig.Serialize()
					lobbyClient.Join(app.actorTokens[idx], i, name, cfg)
				case 'x':
					if idx := slices.Index(app.actorTokens[:], slots.GetActor(i)); idx >= 0 {
						lobbyClient.Leave(app.actorTokens[idx])
					}
				}
				return false
			},
			menu.WithLabelFn(func() string {
				return slots.GetLabel(i)
			}),
			menu.WithDescriptionFn(func() string {
				return slots.GetDescription(i)
			})))
	}
	items = append(items, menu.NewStatic(
		"Waiting for players to join...", "", nil,
		menu.WithVisible(blocker.NeedPlayers),
		menu.WithDisabled(func() bool { return true })))
	items = append(items, menu.NewStatic(
		"Waiting for the host to start...", "", nil,
		menu.WithVisible(blocker.HavePlayers),
		menu.WithDisabled(func() bool { return true })))
	items = append(items, menu.NewStatic(
		"Starting 3...", "", nil,
		menu.WithVisible(blocker.Starting3)))
	items = append(items, menu.NewStatic(
		"Starting 2...", "", nil,
		menu.WithVisible(blocker.Starting2)))
	items = append(items, menu.NewStatic(
		"Starting 1...", "", nil,
		menu.WithVisible(blocker.Starting1)))
	items = append(items, app.menuItemEscape())
	items = append(items, app.menuItemBack())

	m := menu.New(values.ProgramName, func(*menu.Menu) {
		if app.screenIDNext != "" {
			if app.screenIDNext == routeBack {
				lobbyClient.LeaveAll()
			}

			ctx.Stop()
			return
		}
	}, items...)

	return m
}
