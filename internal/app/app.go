// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config"
	"github.com/marko-gacesa/gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udp"
	"github.com/marko-gacesa/udpstar/udpstar"
	"github.com/marko-gacesa/udpstar/udpstar/client"
	"github.com/marko-gacesa/udpstar/udpstar/message"
	"github.com/marko-gacesa/udpstar/udpstar/server"
)

type App struct {
	cfg     config.Config
	cfgPath string

	multicastAddress *net.UDPAddr

	actorTokens [setup.MaxLocalPlayers]message.Token
	clientToken message.Token

	resultSetup               *setup.Setup
	resultClientLobbySelected *udpstar.LobbyListenerInfo
	resultServerSession       *server.Session
	resultClientMap           map[message.Token]server.ClientData
	resultClientSession       *client.Session
	resultServerAddress       net.UDPAddr
	resultLock                sync.Mutex

	screenIDHistory *routes // screen history, the last entry is the id currently active screen
	screenIDNext    route

	screener screen.Screener

	udpService *udp.Service
	gameServer *server.Server

	logger *slog.Logger
	wg     *sync.WaitGroup
}

func NewApp(ctx context.Context, logger *slog.Logger, cfg config.Config, cfgPath string) *App {
	wg := &sync.WaitGroup{}

	udpService := udp.NewService(ctx, cfg.Network.Port,
		udp.WithLogger(logger),
		udp.WithServerStateCallback(func(state udp.ServerState, err error) {
			logger.Debug("server state changed", "state", state)
		}),
		udp.WithIdleTimeout(30*time.Second),
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		udpService.WaitDone()
	}()

	gameServer := server.NewServer(udpService,
		server.WithLogger(logger),
		server.WithBroadcastAddress(cfg.Network.GetMulticastAddress()))

	wg.Add(1)
	go func() {
		defer wg.Done()
		gameServer.Start(ctx)
	}()

	return &App{
		cfg:     cfg,
		cfgPath: cfgPath,
		actorTokens: [4]message.Token{
			message.RandomToken(), message.RandomToken(), message.RandomToken(), message.RandomToken(),
		},
		clientToken:     message.RandomToken(),
		screenIDHistory: (&routes{}).push(routeMain),
		udpService:      udpService,
		gameServer:      gameServer,
		logger:          logger,
		wg:              wg,
	}
}

func (app *App) WaitDone() {
	app.wg.Wait()
}

func (app *App) SetScreener(screener screen.Screener) {
	app.screener = screener
}

func (app *App) Log() *slog.Logger {
	return app.logger
}

func (app *App) LocalPlayerName(i byte) string {
	return app.cfg.LocalPlayers.Infos[i].Name
}

func (app *App) LocalPlayerConfig(i byte) config.PlayerConfig {
	return app.cfg.LocalPlayers.Infos[i].GameConfig
}

func (app *App) LocalPlayer(token message.Token) (config.PlayerInfo, int) {
	for i, t := range app.actorTokens {
		if t == token {
			return app.cfg.LocalPlayers.Infos[i], i
		}
	}

	return config.PlayerInfo{}, -1
}

const (
	routeMain = "main"
	routeQuit = "quit"
	routeBack = "back"

	routeSinglePlayerPrefix      = "1p|"
	routeSinglePlayerMenu        = routeSinglePlayerPrefix + "menu"
	routeSinglePlayerPresetGameN = routeSinglePlayerPrefix + "preset-game:"
	routeSinglePayerCustomSetup  = routeSinglePlayerPrefix + "custom-setup"
	routeSinglePayerCustomGame   = routeSinglePlayerPrefix + "custom-game"

	routeSinglePlayerPresetEditMenu = routeSinglePlayerPrefix + "preset-edit-menu"
	routeSinglePlayerPresetEditN    = routeSinglePlayerPrefix + "preset-edit:"

	routeMultiPlayerPrefix         = "mp|"
	routeMultiPlayerPresetEditMenu = routeMultiPlayerPrefix + "preset-edit-menu"
	routeMultiPlayerPresetEditN    = routeMultiPlayerPrefix + "preset-edit:"

	routeMultiPlayerLocalPrefix      = "mp-local|"
	routeMultiPlayerLocalMenu        = routeMultiPlayerLocalPrefix + "menu"
	routeMultiPlayerLocalPresetGameN = routeMultiPlayerLocalPrefix + "preset-game:"
	routeMultiPayerLocalCustomSetup  = routeMultiPlayerLocalPrefix + "custom-setup"
	routeMultiPayerLocalCustomGame   = routeMultiPlayerLocalPrefix + "custom-game"

	routeMultiPlayerLANPrefix         = "mp-lan|"
	routeMultiPlayerLANMenu           = routeMultiPlayerLANPrefix + "menu"
	routeMultiPlayerLANHostMenu       = routeMultiPlayerLANPrefix + "host-menu"
	routeMultiPlayerLANHostPresetN    = routeMultiPlayerLANPrefix + "host-preset:"
	routeMultiPayerLANHostCustomSetup = routeMultiPlayerLANPrefix + "host-custom-setup"
	routeMultiPayerLANHostLobby       = routeMultiPlayerLANPrefix + "host-lobby"
	routeMultiPlayerLANJoinListen     = routeMultiPlayerLANPrefix + "join-listen"
	routeMultiPlayerLANJoinLobby      = routeMultiPlayerLANPrefix + "join-lobby"

	routeMultiPlayerLANHostGame = routeMultiPlayerLANPrefix + "host-game"
	routeMultiPlayerLANJoinGame = routeMultiPlayerLANPrefix + "join-game"

	routeConfigPrefix           = "config|"
	routeConfigMenu             = routeConfigPrefix + "menu"
	routeConfigLocalPlayerN     = routeConfigPrefix + "local-player:"
	routeConfigLocalPlayerSetup = routeConfigPrefix + "local-player-setup"

	routeConfigVideoPrefix = "config-video|"
	routeConfigVideoSetup  = routeConfigVideoPrefix + "setup"

	routeAboutPrefix = "about|"
	routeAboutMenu   = routeAboutPrefix + "menu"
)

func (app *App) MakeScreen(parentCtx context.Context) (screen.Screen, <-chan struct{}) {
	id := app.screenIDHistory.curr()
	var data any

	ctx := screen.NewContext(parentCtx)

	switch {
	case id == "" || id == routeQuit:
		data = nil
	case id == routeMain:
		data = app.menuMain(ctx)

	// Single player

	case id == routeSinglePlayerMenu:
		data = app.menuSinglePlayer(ctx)
	case strings.HasPrefix(string(id), routeSinglePlayerPresetGameN):
		s := strings.TrimPrefix(string(id), routeSinglePlayerPresetGameN)
		idx, err := strconv.Atoi(s)
		if err != nil {
			data = app.menuError(ctx, err)
		} else {
			app.loadPresetSingle(idx)
			data = app.gameSinglePlayer(ctx)
		}
	case id == routeSinglePayerCustomSetup:
		data = app.menuSinglePlayerSetup(ctx, -1, routeSinglePayerCustomGame)
	case id == routeSinglePayerCustomGame:
		data = app.gameSinglePlayer(ctx)

	// Single player presets edit

	case id == routeSinglePlayerPresetEditMenu:
		data = app.menuSingleEditPresets(ctx)
	case strings.HasPrefix(string(id), routeSinglePlayerPresetEditN):
		s := strings.TrimPrefix(string(id), routeSinglePlayerPresetEditN)
		idx, err := strconv.Atoi(s)
		if err != nil {
			data = app.menuError(ctx, err)
		} else {
			data = app.menuSinglePlayerSetup(ctx, idx, routeBack)
		}

	// Multi-player presets edit

	case id == routeMultiPlayerPresetEditMenu:
		data = app.menuMultiPlayerEditPresets(ctx)
	case strings.HasPrefix(string(id), routeMultiPlayerPresetEditN):
		s := strings.TrimPrefix(string(id), routeMultiPlayerPresetEditN)
		idx, err := strconv.Atoi(s)
		if err != nil {
			data = app.menuError(ctx, err)
		} else {
			data = app.menuMultiPlayerSetup(ctx, setup.MaxPlayers, idx, routeBack)
		}

	// Multi-player local

	case id == routeMultiPlayerLocalMenu:
		data = app.menuMultiPlayerLocal(ctx)
	case strings.HasPrefix(string(id), routeMultiPlayerLocalPresetGameN):
		s := strings.TrimPrefix(string(id), routeMultiPlayerLocalPresetGameN)
		idx, err := strconv.Atoi(s)
		if err != nil {
			data = app.menuError(ctx, err)
		} else {
			app.loadPresetMulti(idx)
			if app.resultSetup != nil && app.resultSetup.PlayerCount() > setup.MaxLocalPlayers {
				data = app.menuErrorText(ctx, fmt.Sprintf("The preset defined too many players. Maximum is %d", setup.MaxLocalPlayers))
			} else {
				data = app.gameMultiPlayerLocal(ctx)
			}
		}
	case id == routeMultiPayerLocalCustomSetup:
		data = app.menuMultiPlayerSetup(ctx, setup.MaxLocalPlayers, -1, routeMultiPayerLocalCustomGame)
	case id == routeMultiPayerLocalCustomGame:
		data = app.gameMultiPlayerLocal(ctx)

	// Multi-player LAN

	case id == routeMultiPlayerLANMenu:
		data = app.menuMultiPlayerLANMain(ctx)
	case id == routeMultiPlayerLANHostMenu:
		data = app.menuMultiPlayerLANHostMenu(ctx)
	case strings.HasPrefix(string(id), routeMultiPlayerLANHostPresetN):
		s := strings.TrimPrefix(string(id), routeMultiPlayerLANHostPresetN)
		idx, err := strconv.Atoi(s)
		if err != nil {
			data = app.menuError(ctx, err)
		} else {
			app.loadPresetMulti(idx)
			data = app.menuMultiPlayerLANHostLobby(ctx)
		}
	case id == routeMultiPayerLANHostCustomSetup:
		data = app.menuMultiPlayerSetup(ctx, setup.MaxPlayers, -1, routeMultiPayerLANHostLobby)
	case id == routeMultiPayerLANHostLobby:
		data = app.menuMultiPlayerLANHostLobby(ctx)

	case id == routeMultiPlayerLANJoinListen:
		data = app.menuMultiPlayerLANJoinListen(ctx)
	case id == routeMultiPlayerLANJoinLobby:
		data = app.menuMultiPlayerLANJoinLobby(ctx)

	case id == routeMultiPlayerLANHostGame:
		data = app.gameUDPServer(ctx)
	case id == routeMultiPlayerLANJoinGame:
		data = app.gameUDPClient(ctx)

	case id == routeConfigMenu:
		data = app.menuConfig(ctx)
	case strings.HasPrefix(string(id), routeConfigLocalPlayerN):
		s := strings.TrimPrefix(string(id), routeConfigLocalPlayerN)
		idx, err := strconv.Atoi(s)
		if err != nil {
			data = app.menuError(ctx, err)
		} else {
			data = app.menuConfigLocalPlayer(ctx, idx)
		}
	case id == routeConfigVideoSetup:
		data = app.menuConfigVideo(ctx)

	case id == routeAboutMenu:
		data = app.menuAbout(ctx)
	}

	return app.screener.Screen(ctx, data), ctx.Done()
}

func (app *App) returnToMainScreen() {
	app.screenIDHistory.clear()
	app.screenIDNext = routeMain
}

func (app *App) ScreenFinish() {
	if app.screenIDNext == routeBack || app.screenIDNext == "" {
		app.screenIDHistory.pop()
		app.screenIDNext = ""
	} else if app.screenIDNext != "" {
		app.screenIDHistory.push(app.screenIDNext)
		app.screenIDNext = ""
	}
}

func (app *App) LogError(err error, msg string) {
	app.logger.Error(msg, "error", err)
}

func (app *App) loadPresetSingle(idx int) {
	s := app.cfg.Presets.Single[idx]
	if !s.MiscOptions.CustomSeed {
		s.MiscOptions.Seed = rand.Int64()
	}

	if s.SanitizeSingle() {
		app.logger.Warn("loading preset multi: sanitation is required")
	}

	app.resultSetup = &s
}

func (app *App) loadPresetMulti(idx int) {
	s := app.cfg.Presets.Multi[idx]
	if !s.MiscOptions.CustomSeed {
		s.MiscOptions.Seed = rand.Int64()
	}

	if s.SanitizeMulti() {
		app.logger.Warn("loading preset multi: sanitation is required")
	}

	app.resultSetup = &s
}
