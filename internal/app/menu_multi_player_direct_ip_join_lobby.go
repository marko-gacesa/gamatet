// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/marko-gacesa/bitdata"
	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udpstar/client"
	"github.com/marko-gacesa/udpstar/udpstar/message"
)

func (app *App) menuMultiPlayerDirectIPJoinLobby(ctx screen.Context) *menu.Menu {
	addr, err := net.ResolveIPAddr("ip", app.cfg.Network.DirectIPAddress)
	if err != nil || app.resultToken == 0 {
		return app.menuError(ctx, errorInputMissing)
	}

	app.resultClientSession = nil

	lobbyToken := app.resultToken

	serverAddr := net.UDPAddr{
		IP:   addr.IP,
		Port: app.cfg.Network.Port,
		Zone: addr.Zone,
	}

	playerName := app.LocalPlayerName(0)
	playerCfg := app.LocalPlayerConfig(0).Serialize()

	slots := makeLobbyEntries([]message.Token{lobbyToken, lobbyToken}, app.actorTokens[:], withFixedSlots())
	blocker := makeStartBlocker()

	// lobby

	udpSender := udpSender{
		addr: serverAddr,
		srv:  app.udpService,
	}

	lobbyClient := client.NewLobby(udpSender, lobbyToken, app.clientToken, client.WithLobbyLogger(app.logger))
	lobbyJoin := func() {
		lobbyClient.Join(app.actorTokens[0], 1, playerName, playerCfg)
	}

	ctxLobbyJoin, stopLobbyJoin := context.WithCancel(ctx)

	go func(ctx context.Context) {
		lobbyJoin()

		t := time.NewTicker(time.Second)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				lobbyJoin()
			case <-ctx.Done():
				return
			}
		}
	}(ctxLobbyJoin)

	go func() {
		gameSession := lobbyClient.Start(ctx)

		if gameSession != nil {
			app.screenIDNext = routeMultiPlayerUDPJoinGame

			app.resultServerAddress = serverAddr
			app.resultClientSession = gameSession
		}

		ctx.Stop()
	}()

	var version int
	var gameStr atomic.Value
	gameStr.Store("")

	err = app.udpService.Handle(ctx, func(data []byte, addr net.UDPAddr) []byte {
		lobbyClient.HandleIncomingMessages(data)

		lobby, age := lobbyClient.Get(version)
		if lobby == nil || len(lobby.Slots) != 2 {
			return nil
		}

		r := bitdata.NewReaderError(lobby.Def)
		var o setup.Setup
		o.Read(r)
		if err := r.Error(); err != nil {
			gameStr.Store("?")
		} else {
			gameStr.Store(o.String())
		}

		if age > time.Minute {
			ctx.Stop()
			return nil
		}

		version = lobby.Version

		storyTokens := []message.Token{lobby.Slots[0].StoryToken, lobby.Slots[1].StoryToken}

		slots.updateStoryTokens(storyTokens)
		slots.setAll(lobby)

		blocker.update(lobby)

		stopLobbyJoin()

		return nil
	})
	if err != nil {
		return app.menuError(ctx, err)
	}

	items := make([]menu.Item, 0, 12)

	items = append(items, menu.NewStatic("", "", nil,
		menu.WithLabelFn(func() string { return gameStr.Load().(string) }),
		menu.WithVisible(func() bool { return gameStr.Load().(string) != "" }),
		menu.WithDisabled(func() bool { return true })))
	items = append(items, menu.NewStatic("", "", nil,
		menu.WithLabelFn(func() string { return slots.GetLabel(0) }),
		menu.WithDescriptionFn(func() string { return slots.GetDescription(0) }),
		menu.WithVisible(func() bool { return gameStr.Load().(string) != "" })))
	items = append(items, menu.NewStatic("", "", nil,
		menu.WithLabelFn(func() string { return slots.GetLabel(1) }),
		menu.WithDescriptionFn(func() string { return slots.GetDescription(1) }),
		menu.WithVisible(func() bool { return gameStr.Load().(string) != "" })))
	items = append(items, menu.NewStatic(
		T(KeyLobbyIssueIncomplete), T(KeyLobbyIssueIncompleteDesc), nil,
		menu.WithVisible(blocker.NeedPlayers),
		menu.WithDisabled(func() bool { return true })))
	items = append(items, menu.NewStatic(
		T(KeyLobbyIssueWaitingForHost), T(KeyLobbyIssueWaitingForHostDesc), nil,
		menu.WithVisible(blocker.HavePlayers),
		menu.WithDisabled(func() bool { return true })))
	items = append(items, menu.NewStatic(
		T(KeyLobbyStarting3), "", nil,
		menu.WithVisible(blocker.Starting3)))
	items = append(items, menu.NewStatic(
		T(KeyLobbyStarting2), "", nil,
		menu.WithVisible(blocker.Starting2)))
	items = append(items, menu.NewStatic(
		T(KeyLobbyStarting1), "", nil,
		menu.WithVisible(blocker.Starting1)))
	items = append(items, app.menuItemEscape())
	items = append(items, app.menuItemBack())

	m := menu.New(T(KeyLobbyTitle), func(m *menu.Menu) {
		if app.screenIDNext == routeBack {
			lobbyClient.LeaveAll()
			time.Sleep(10 * time.Millisecond)
		}
		app.menuStopper(ctx)(m)
	}, items...)

	return m
}
