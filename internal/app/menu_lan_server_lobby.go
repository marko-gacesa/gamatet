// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"
	"gamatet/internal/values"
	"gamatet/logic/menu"
	"gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udpstar"
	"github.com/marko-gacesa/udpstar/udpstar/message"
	"github.com/marko-gacesa/udpstar/udpstar/server"
	"net"
	"slices"
	"sync"
	"time"
	"unicode"
)

func (app *App) menuLANServerLobby(ctx screen.Context) *menu.Menu {
	if app.resultSetup == nil {
		return app.menuErrorText(ctx, "Input missing")
	}

	app.resultServerSession = nil
	app.resultClientMap = nil

	lobbyToken := message.RandomToken()

	app.cfg.Presets.Multi[0] = *app.resultSetup

	gameDef := app.resultSetup.Def()
	gameStr := app.resultSetup.String()
	playerCount := app.resultSetup.GameOptions.PlayerCount()
	slotStories := app.resultSetup.GameOptions.CreateSlotsStories()

	var start int

	slots := makeLobbyEntries(slotStories, true, app.actorTokens[:])
	blocker := makeStartBlocker()

	// Prepare menu

	items := make([]menu.Item, 0, 10+playerCount)

	items = append(items, menu.NewStatic(
		gameStr, "", nil,
		menu.WithDisabled(func() bool { return true })))
	for i := range playerCount {
		items = append(items, menu.NewStatic("", "",
			func(r rune) bool {
				switch unicode.ToLower(r) {
				case '1', '2', '3', '4':
					idx := byte(r - '1')
					app.gameServer.JoinLocal(lobbyToken, app.actorTokens[idx], i, idx, app.LocalPlayerName(idx))
				case 'x', '\b', '\xFF':
					app.gameServer.EvictIdx(lobbyToken, i)
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
	items = append(items, menu.NewCommand(&start, 1,
		itemTextPrefixForward+"Start game", "",
		menu.WithVisible(blocker.CanStart),
		menu.WithDisabled(func() bool { return start > 0 })))
	items = append(items, menu.NewStatic(
		"Can't start - players are missing", "", nil,
		menu.WithVisible(blocker.NeedPlayers),
		menu.WithDisabled(func() bool { return true })))
	items = append(items, menu.NewStatic(
		"Can't start - no remote players", "", nil,
		menu.WithVisible(blocker.NeedRemotesProblem),
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
			ctx.Stop()
			return
		}

		if start == 1 {
			start++
			go func() {
				app.resultServerSession, app.resultClientMap, _ = app.gameServer.FinishLobby(ctx, lobbyToken)
				app.screenIDNext = routeGameUDPServer
				ctx.Stop()
			}()
		}
	}, items...)

	// Start UDP server

	err := app.udpService.Handle(ctx, func(data []byte, addr net.UDPAddr) []byte {
		return app.gameServer.HandleIncomingMessages(data, addr)
	})
	if err != nil {
		return app.menuError(ctx, err)
	}

	// Start game server

	err = app.gameServer.StartLobby(ctx, &server.LobbySetup{
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
		}
	}()

	return m
}

type lobbyEntries struct {
	isHost      bool
	localActors []message.Token

	teams   map[message.Token]string
	entries []lobbyEntry
	mx      sync.Mutex
}

type lobbyEntry struct {
	text         string
	description  string
	availability udpstar.Availability
	actor        message.Token
}

func makeLobbyEntries(slotStories []message.Token, isHost bool, localActors []message.Token) *lobbyEntries {
	teams := map[message.Token]string{}
	for _, storyToken := range slotStories {
		if _, ok := teams[storyToken]; !ok {
			teams[storyToken] = fmt.Sprintf("Team %d", len(teams)+1)
		}
	}
	if len(teams) == 1 {
		for storyToken := range teams {
			teams[storyToken] = "Team"
		}
	}

	return &lobbyEntries{
		isHost:      isHost,
		localActors: localActors,
		teams:       teams,
		entries:     make([]lobbyEntry, len(slotStories)),
		mx:          sync.Mutex{},
	}
}

func (l *lobbyEntries) setAll(lobby *udpstar.Lobby) {
	l.mx.Lock()
	for i := range lobby.Slots {
		avail := lobby.Slots[i].Availability
		actor := lobby.Slots[i].ActorToken

		l.entries[i].availability = avail
		l.entries[i].actor = actor

		name := lobby.Slots[i].Name
		team := l.teams[lobby.Slots[i].StoryToken]
		switch avail {
		case udpstar.SlotAvailable:
			l.entries[i].text = fmt.Sprintf("\t%s: <Available>", team)
			l.entries[i].description = "Press '1', '2', '3' or '4' to join game as a local player"
		case udpstar.SlotLocal0, udpstar.SlotLocal1, udpstar.SlotLocal2, udpstar.SlotLocal3:
			l.entries[i].text = fmt.Sprintf("\t%s: %s [local %d]", team, name, avail-udpstar.SlotLocal0+1)
			if l.isHost {
				l.entries[i].description = "Press X to leave and make the place available to other players"
			}
		case udpstar.SlotRemote:
			latency := lobby.Slots[i].Latency
			l.entries[i].text = fmt.Sprintf("\t%s: %s [latency %dms]", team, name, latency.Milliseconds())
			if l.isHost {
				l.entries[i].description = "Press X to evict the remote player"
			} else if slices.Contains(l.localActors, actor) {
				l.entries[i].description = "Press X to leave and make the place available to other players"
			}
		}
	}
	l.mx.Unlock()
}

func (l *lobbyEntries) GetActor(i byte) message.Token {
	l.mx.Lock()
	s := l.entries[i].actor
	l.mx.Unlock()
	return s
}

func (l *lobbyEntries) GetLabel(i byte) string {
	l.mx.Lock()
	s := l.entries[i].text
	l.mx.Unlock()
	return s
}

func (l *lobbyEntries) GetDescription(i byte) string {
	l.mx.Lock()
	s := l.entries[i].description
	l.mx.Unlock()
	return s
}

type startBlocker struct {
	lobbyState udpstar.LobbyState
	remotes    byte
	mx         sync.Mutex
}

func makeStartBlocker() *startBlocker {
	return &startBlocker{
		lobbyState: udpstar.LobbyStateActive,
		remotes:    0,
		mx:         sync.Mutex{},
	}
}

func (b *startBlocker) update(lobby *udpstar.Lobby) {
	b.mx.Lock()
	b.lobbyState = lobby.State
	b.remotes = 0
	for i := range lobby.Slots {
		avail := lobby.Slots[i].Availability
		if avail == udpstar.SlotRemote {
			b.remotes++
		}
	}
	b.mx.Unlock()
}

func (b *startBlocker) CanCancel() bool {
	b.mx.Lock()
	q := b.lobbyState == udpstar.LobbyStateActive || b.lobbyState == udpstar.LobbyStateReady
	b.mx.Unlock()
	return q
}

func (b *startBlocker) CanStart() bool {
	b.mx.Lock()
	q := b.lobbyState == udpstar.LobbyStateReady && b.remotes > 0
	b.mx.Unlock()
	return q
}

func (b *startBlocker) NeedPlayers() bool {
	b.mx.Lock()
	q := b.lobbyState == udpstar.LobbyStateActive
	b.mx.Unlock()
	return q
}

func (b *startBlocker) HavePlayers() bool {
	b.mx.Lock()
	q := b.lobbyState == udpstar.LobbyStateReady
	b.mx.Unlock()
	return q
}

func (b *startBlocker) NeedRemotesProblem() bool {
	b.mx.Lock()
	q := b.lobbyState == udpstar.LobbyStateReady && b.remotes == 0
	b.mx.Unlock()
	return q
}

func (b *startBlocker) Starting1() bool {
	b.mx.Lock()
	q := b.lobbyState == udpstar.LobbyStateStarting1
	b.mx.Unlock()
	return q
}

func (b *startBlocker) Starting2() bool {
	b.mx.Lock()
	q := b.lobbyState == udpstar.LobbyStateStarting2
	b.mx.Unlock()
	return q
}

func (b *startBlocker) Starting3() bool {
	b.mx.Lock()
	q := b.lobbyState == udpstar.LobbyStateStarting3
	b.mx.Unlock()
	return q
}
