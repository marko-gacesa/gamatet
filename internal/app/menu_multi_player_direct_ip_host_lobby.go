// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"fmt"
	"net"
	"time"

	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udpstar/message"
	"github.com/marko-gacesa/udpstar/udpstar/server"
)

func (app *App) menuMultiPlayerDirectIPHostLobby(ctx screen.Context) *menu.Menu {
	addr, err := net.ResolveIPAddr("ip", app.cfg.Network.DirectIPAddress)
	if err != nil || app.resultSetup == nil || app.resultSetup.GameOptions.PlayerCount() != 2 ||
		app.resultSetup.GameOptions.FieldCount > 2 || app.resultSetup.GameOptions.FieldCount < 1 ||
		app.resultToken == 0 {
		return app.menuError(ctx, errorInputMissing)
	}

	udpAddr := net.UDPAddr{
		IP:   addr.IP,
		Port: app.cfg.Network.Port,
		Zone: addr.Zone,
	}

	app.resultServerSession = nil
	app.resultClientMap = nil

	lobbyToken := app.resultToken

	gameDef := app.resultSetup.Def()
	gameStr := app.resultSetup.String()

	slotStories := make([]message.Token, 2)
	switch app.resultSetup.GameOptions.FieldCount {
	case 1:
		slotStories[0] = lobbyToken
		slotStories[1] = lobbyToken
	case 2:
		slotStories[0] = message.RandomToken()
		slotStories[1] = lobbyToken
	default:
		panic("should not happen: invalid field count")
	}

	var start int

	slots := makeLobbyEntries(slotStories, app.actorTokens[:], withHost(), withFixedSlots())
	blocker := makeStartBlocker()

	// Prepare menu

	items := make([]menu.Item, 0, 12)

	items = append(items, menu.NewStatic(
		gameStr, "", nil,
		menu.WithDisabled(func() bool { return true })))
	items = append(items, menu.NewStatic("", "", nil,
		menu.WithLabelFn(func() string { return slots.GetLabel(0) }),
		menu.WithDescriptionFn(func() string { return slots.GetDescription(0) })))
	items = append(items, menu.NewStatic("", "", nil,
		menu.WithLabelFn(func() string { return slots.GetLabel(1) }),
		menu.WithDescriptionFn(func() string { return slots.GetDescription(1) })))
	items = append(items, menu.NewCommand(&start, 1,
		T(KeyLobbyStartGame), T(KeyLobbyStartGameDesc),
		menu.WithVisible(blocker.CanStart),
		menu.WithDisabled(func() bool { return start > 0 })))
	items = append(items, menu.NewStatic(
		T(KeyLobbyIssueMissingPlayers), T(KeyLobbyIssueMissingPlayersDesc), nil,
		menu.WithVisible(blocker.NeedPlayers),
		menu.WithDisabled(func() bool { return true })))
	items = append(items, menu.NewStatic(
		T(KeyLobbyIssueNoRemotePlayers), T(KeyLobbyIssueNoRemotePlayersDesc), nil,
		menu.WithVisible(blocker.NeedRemotesProblem),
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
		app.menuStopper(ctx)(m)

		if start == 1 {
			start++
			go func() {
				app.resultServerSession, app.resultClientMap, err = app.gameServer.FinishLobby(ctx, lobbyToken)
				if err != nil {
					app.logger.Error("failed to finish lobby",
						"err", err,
						"lobbyToken", lobbyToken,
					)
					l, err := app.gameServer.GetLobby(lobbyToken, -1)
					if err != nil {
						app.logger.Error("failed to get lobby", "err", err, "lobbyToken", lobbyToken)
					}
					if l != nil {
						app.logger.Info("lobby", "state", l.State.String(), "slots",
							fmt.Sprintf("%+v", l.Slots))
					}
					app.screenIDNext = routeBack
					ctx.Stop()
					return
				}
				app.screenIDNext = routeMultiPlayerUDPHostGame
				ctx.Stop()
			}()
		}
	}, items...)

	// Start UDP server

	err = app.udpService.Handle(ctx, func(data []byte, addr net.UDPAddr) []byte {
		return app.gameServer.HandleIncomingMessages(data, addr)
	})
	if err != nil {
		return app.menuError(ctx, err)
	}

	// Start game server

	err = app.gameServer.StartLobbyNoBroadcast(ctx, &server.LobbySetup{
		Token:       lobbyToken,
		Name:        app.resultSetup.Name,
		Def:         gameDef,
		SlotStories: slotStories,
	})
	if err != nil {
		return app.menuError(ctx, err)
	}

	_ = app.gameServer.JoinLocal(lobbyToken, app.actorTokens[0], 0, 0, app.LocalPlayerName(0))

	// Start UI refresher

	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	go func() {
		var version int
		for {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				return
			}

			lobby, err := app.gameServer.GetLobby(lobbyToken, version)
			if err != nil {
				app.menuShowError(m, err)
				continue
			}
			if lobby == nil {
				continue
			}

			version = lobby.Version

			slots.setAll(lobby)
			blocker.update(lobby)

			if !blocker.CanStart() {
				// UDP hole punching
				app.udpService.Send([]byte("punch a hole"), udpAddr)
			}
		}
	}()

	return m
}
