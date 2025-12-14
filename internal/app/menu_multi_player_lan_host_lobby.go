// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	"unicode"

	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udpstar"
	"github.com/marko-gacesa/udpstar/udpstar/message"
	"github.com/marko-gacesa/udpstar/udpstar/server"
)

func (app *App) menuMultiPlayerLANHostLobby(ctx screen.Context) *menu.Menu {
	if app.resultSetup == nil {
		return app.menuError(ctx, errorInputMissing)
	}

	app.resultServerSession = nil
	app.resultClientMap = nil

	lobbyToken := message.RandomToken()

	gameDef := app.resultSetup.Def()
	gameStr := app.resultSetup.String()
	playerCount := app.resultSetup.GameOptions.PlayerCount()
	slotStories := app.resultSetup.GameOptions.CreateSlotsStories()

	var start int

	slots := makeLobbyEntries(slotStories, app.actorTokens[:], withHost())
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
				case '\b', '\xFF':
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
	items = append(items, app.menuItemEscape(menu.WithDisabled(func() bool { return blocker.Starting() })))
	items = append(items, app.menuItemBack(menu.WithDisabled(func() bool { return blocker.Starting() })))

	m := menu.New(T(KeyLobbyTitle), func(m *menu.Menu) {
		app.menuStopper(ctx)(m)

		if start == 1 {
			start++
			go func() {
				var err error

				app.resultServerSession, app.resultClientMap, err = app.gameServer.FinishLobby(ctx, lobbyToken)
				if err != nil && !errors.Is(err, context.Canceled) {
					app.logger.Error("failed to finish lobby", "err", err)
				}

				if app.resultServerSession == nil || app.resultClientMap == nil {
					app.screenIDNext = routeBack
				} else {
					app.screenIDNext = routeMultiPlayerUDPHostGame
				}

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
	config lobbyEntriesConfig

	localActors []message.Token

	slotInfos []slotInfo
	entries   []lobbyEntry
	mx        sync.Mutex
}

type lobbyEntriesConfig struct {
	host       bool
	fixedSlots bool
}

type slotInfo struct {
	storyIndex      int
	storyEntryIndex int
	teamName        string
}

type lobbyEntry struct {
	text         string
	description  string
	availability udpstar.Availability
	actor        message.Token
}

func makeLobbyEntries(
	slotStories []message.Token,
	localActors []message.Token,
	options ...func(entries *lobbyEntries),
) *lobbyEntries {
	l := &lobbyEntries{
		config:      lobbyEntriesConfig{},
		localActors: localActors,
		slotInfos:   nil,
		entries:     make([]lobbyEntry, len(slotStories)),
		mx:          sync.Mutex{},
	}

	l.updateStoryTokens(slotStories)

	for _, option := range options {
		option(l)
	}

	return l
}

func withHost() func(entries *lobbyEntries) {
	return func(entries *lobbyEntries) {
		entries.config.host = true
	}
}

func withFixedSlots() func(entries *lobbyEntries) {
	return func(entries *lobbyEntries) {
		entries.config.fixedSlots = true
	}
}

func (l *lobbyEntries) updateStoryTokens(slotStories []message.Token) {
	teams := map[message.Token]string{}
	for _, storyToken := range slotStories {
		if _, ok := teams[storyToken]; !ok {
			teams[storyToken] = T(KeyLobbyTeam) + " " + strconv.Itoa(len(teams)+1)
		}
	}
	if len(teams) == 1 {
		for storyToken := range teams {
			teams[storyToken] = T(KeyLobbyTeam)
		}
	}

	slotInfos := make([]slotInfo, len(slotStories))
	for i := range slotStories {
		if i == 0 {
			slotInfos[0] = slotInfo{
				storyIndex:      0,
				storyEntryIndex: 0,
				teamName:        teams[slotStories[i]],
			}
		} else if slotStories[i] == slotStories[i-1] {
			slotInfos[i] = slotInfo{
				storyIndex:      slotInfos[i-1].storyIndex,
				storyEntryIndex: slotInfos[i-1].storyEntryIndex + 1,
				teamName:        teams[slotStories[i]],
			}
		} else {
			slotInfos[i] = slotInfo{
				storyIndex:      slotInfos[i-1].storyIndex + 1,
				storyEntryIndex: 0,
				teamName:        teams[slotStories[i]],
			}
		}
	}

	l.mx.Lock()
	l.slotInfos = slotInfos
	l.mx.Unlock()
}

func (l *lobbyEntries) setAll(lobby *udpstar.Lobby) {
	l.mx.Lock()
	for i := range lobby.Slots {
		avail := lobby.Slots[i].Availability
		actor := lobby.Slots[i].ActorToken

		l.entries[i].availability = avail
		l.entries[i].actor = actor

		info := l.slotInfos[i]

		name := playerName(lobby.Slots[i].Name, info.storyIndex, info.storyEntryIndex, i)
		team := info.teamName

		switch avail {
		case udpstar.SlotAvailable:
			if l.config.fixedSlots {
				l.entries[i].text = fmt.Sprintf("\t%s: <%s>",
					team, T(KeyLobbyFixedSlotWaiting))
				l.entries[i].description = T(KeyLobbyFixedSlotWaitingDesc)
			} else {
				l.entries[i].text = fmt.Sprintf("\t%s: <%s>",
					team, T(KeyLobbySlotAvailable))
				l.entries[i].description = T(KeyLobbySlotAvailableDesc)
			}
		case udpstar.SlotLocal0, udpstar.SlotLocal1, udpstar.SlotLocal2, udpstar.SlotLocal3:
			const fmtLocal = "\t%s: %s [%s %d]"

			if l.config.fixedSlots {
				l.entries[i].text = fmt.Sprintf(fmtLocal, team, name, T(KeyLobbyFixedSlotLocal), avail-udpstar.SlotLocal0+1)
				if l.config.host {
					l.entries[i].description = T(KeyLobbyFixedSlotLocalDesc)
				} else {
					l.entries[i].description = T(KeyLobbyFixedSlotRemoteDesc)
				}
			} else if l.config.host {
				l.entries[i].text = fmt.Sprintf(fmtLocal, team, name, T(KeyLobbySlotLocal), avail-udpstar.SlotLocal0+1)
				l.entries[i].description = T(KeyLobbySlotLocalDesc)
			} else {
				l.entries[i].text = fmt.Sprintf(fmtLocal, team, name, T(KeyLobbySlotLocal), avail-udpstar.SlotLocal0+1)
				l.entries[i].description = ""
			}
		case udpstar.SlotRemote:
			const fmtRemote = "\t%s: %s [%s, %s=%dms]"
			latency := lobby.Slots[i].Latency

			if l.config.fixedSlots {
				l.entries[i].text = fmt.Sprintf(fmtRemote, team, name, T(KeyLobbyFixedSlotRemote), T(KeyLobbyLatency), latency.Milliseconds())
				if l.config.host {
					l.entries[i].description = T(KeyLobbyFixedSlotRemoteDesc)
				} else {
					l.entries[i].description = T(KeyLobbyFixedSlotLocalDesc)
				}
			} else if l.config.host {
				l.entries[i].text = fmt.Sprintf(fmtRemote,
					team, name, T(KeyLobbySlotRemote), T(KeyLobbyLatency), latency.Milliseconds())
				l.entries[i].description = T(KeyLobbySlotRemoteDesc)
			} else {
				l.entries[i].text = fmt.Sprintf(fmtRemote,
					team, name, T(KeyLobbySlotRemote), T(KeyLobbyLatency), latency.Milliseconds())
				l.entries[i].description = T(KeyLobbySlotLocalDesc)
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

func (b *startBlocker) Starting() bool {
	b.mx.Lock()
	q := b.lobbyState >= udpstar.LobbyStateStarting
	b.mx.Unlock()
	return q
}
