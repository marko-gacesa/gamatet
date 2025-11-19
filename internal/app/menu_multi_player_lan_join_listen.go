// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/marko-gacesa/bitdata"
	"github.com/marko-gacesa/gamatet/game/setup"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
	"github.com/marko-gacesa/gamatet/logic/menu"
	"github.com/marko-gacesa/gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udp"
	"github.com/marko-gacesa/udpstar/udpstar"
	"github.com/marko-gacesa/udpstar/udpstar/client"
)

const maxClientLobbyEntries = 20

func (app *App) menuMultiPlayerLANJoinListen(ctx screen.Context) *menu.Menu {
	app.resultClientLobbySelected = nil

	list := makeClientLobbyList()

	// create UI interface

	items := make([]menu.Item, 0, 3+maxClientLobbyEntries)

	for i := range maxClientLobbyEntries {
		items = append(items, menu.NewCommand(&app.resultClientLobbySelected, list.ResultPtr(i), "", "",
			menu.WithLabelFn(func() string { return list.Label(i) }),
			menu.WithVisible(func() bool { return list.Exists(i) }),
		))
	}
	items = append(items, menu.NewStatic(T(KeyMenuJoinLANListening), T(KeyMenuJoinLANListeningDesc), nil,
		menu.WithVisible(list.Empty)))
	items = append(items, app.menuItemEscape())
	items = append(items, app.menuItemBack())

	m := menu.New(T(KeyMenuJoinLANTitle), func(*menu.Menu) {
		if app.screenIDNext != "" {
			ctx.Stop()
			return
		}

		if app.resultClientLobbySelected != nil {
			app.screenIDNext = routeMultiPlayerLANJoinLobby
			ctx.Stop()
			return
		}
	}, items...)

	// lobby listener with two processes: periodic refresher and multicast listener.

	lobbyListener := client.NewLobbyListener(app.clientToken, client.WithLobbyListenerLogger(app.logger))

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				lobbyListener.Refresh()
				lobbies, version := lobbyListener.List(list.Version())
				needsFocus := list.Update(lobbies, version)
				if needsFocus >= 0 {
					m.Focus(needsFocus)
				}
			}
		}
	}()

	go func(multicastAddr net.UDPAddr) {
		app.logger.Debug("listening for games",
			"addr", multicastAddr)

		if err := udp.ListenMulticast(ctx, multicastAddr, func(data []byte, addr net.UDPAddr) {
			lobbyListener.HandleBroadcastMessages(data, addr)
			lobbies, version := lobbyListener.List(list.Version())
			needsFocus := list.Update(lobbies, version)
			if needsFocus >= 0 {
				m.Focus(needsFocus)
			}

			app.logger.Debug("received multicast packet",
				"version", version,
				"lobbies_len", len(lobbies),
				"lobbies", lobbies,
				"addr", addr.String())
		}); err != nil {
			app.menuShowError(m, err)
		}
	}(app.cfg.Network.GetMulticastAddress())

	return m
}

type clientLobbyList struct {
	entryList   [maxClientLobbyEntries]udpstar.LobbyListenerInfo
	entryLabels [maxClientLobbyEntries]string
	entryCount  int
	version     int
	mx          sync.Mutex
}

func makeClientLobbyList() *clientLobbyList {
	return &clientLobbyList{
		entryList:   [maxClientLobbyEntries]udpstar.LobbyListenerInfo{},
		entryLabels: [maxClientLobbyEntries]string{},
		entryCount:  0,
		version:     0,
		mx:          sync.Mutex{},
	}
}

func (l *clientLobbyList) Label(idx int) string {
	l.mx.Lock()
	defer l.mx.Unlock()
	return l.entryLabels[idx]
}

func (l *clientLobbyList) Exists(idx int) bool {
	l.mx.Lock()
	defer l.mx.Unlock()
	return idx < l.entryCount
}

func (l *clientLobbyList) Empty() bool {
	l.mx.Lock()
	defer l.mx.Unlock()
	return l.entryCount == 0
}

func (l *clientLobbyList) ResultPtr(idx int) *udpstar.LobbyListenerInfo {
	return &l.entryList[idx]
}

func (l *clientLobbyList) Version() int {
	l.mx.Lock()
	defer l.mx.Unlock()
	return l.version
}

func (l *clientLobbyList) Update(newData []udpstar.LobbyListenerInfo, newVersion int) int {
	l.mx.Lock()
	defer l.mx.Unlock()

	needsFocus := -1

	if newVersion != l.version {
		oldCount := l.entryCount

		if len(newData) > maxClientLobbyEntries {
			newData = newData[:maxClientLobbyEntries]
		}

		l.version = newVersion
		l.entryCount = len(newData)
		copy(l.entryList[:], newData)

		if oldCount == 0 && l.entryCount > 0 {
			needsFocus = 0
		}
	}

	l.refresh()

	return needsFocus
}

func (l *clientLobbyList) Refresh() {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.refresh()
}

func (l *clientLobbyList) refresh() {
	for i := 0; i < l.entryCount; i++ {
		def := l.entryList[i].Lobby.Def
		defStr := "???"

		r := bitdata.NewReaderError(def)
		var o setup.Setup
		o.Read(r)
		if r.Error() == nil {
			defStr = o.String()
		}

		slots := l.entryList[i].Lobby.Slots

		available := 0
		for slotIdx := range slots {
			if slots[slotIdx].Availability == udpstar.SlotAvailable {
				available++
			}
		}

		label := fmt.Sprintf("%s (%s) - %d/%d - [%s] %s",
			l.entryList[i].Lobby.Name,
			defStr,
			available,
			len(slots),
			l.entryList[i].State.String(),
			l.entryList[i].Addr.IP.String())

		l.entryLabels[i] = label
	}
}
