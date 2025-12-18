// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

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
	"github.com/marko-gacesa/gamatet/internal/i18n"
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
	resultToken               message.Token

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

	loggerUDP := logger.With("component", "udp-service")
	udpService := udp.NewService(ctx, cfg.Network.Port,
		udp.WithLogger(loggerUDP),
		udp.WithServerStateCallback(func(state udp.ServerState, err error) {
			loggerUDP.Debug("state changed", "state", state)
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

func (app *App) LocalPlayerConfig(i byte) config.Player {
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

func (app *App) VideoConfig() config.Video {
	return app.cfg.Video
}

const (
	routeMain = "main"
	routeQuit = "quit"
	routeBack = "back"

	routeSinglePlayerPrefix      = "1p|"
	routeSinglePlayerMenu        = routeSinglePlayerPrefix + "menu"
	routeSinglePlayerPresetGameN = routeSinglePlayerPrefix + "preset-game:"
	routeSinglePlayerCustomSetup = routeSinglePlayerPrefix + "custom-setup"
	routeSinglePlayerCustomGame  = routeSinglePlayerPrefix + "custom-game"

	routeSinglePlayerPresetEditMenu = routeSinglePlayerPrefix + "preset-edit-menu"
	routeSinglePlayerPresetEditN    = routeSinglePlayerPrefix + "preset-edit:"

	routeMultiPlayerPrefix         = "mp|"
	routeMultiPlayerPresetEditMenu = routeMultiPlayerPrefix + "preset-edit-menu"
	routeMultiPlayerPresetEditN    = routeMultiPlayerPrefix + "preset-edit:"

	routeMultiPlayerLocalPrefix      = "mp-local|"
	routeMultiPlayerLocalMenu        = routeMultiPlayerLocalPrefix + "menu"
	routeMultiPlayerLocalPresetGameN = routeMultiPlayerLocalPrefix + "preset-game:"
	routeMultiPlayerLocalCustomSetup = routeMultiPlayerLocalPrefix + "custom-setup"
	routeMultiPlayerLocalCustomGame  = routeMultiPlayerLocalPrefix + "custom-game"

	routeMultiPlayerLANPrefix          = "mp-lan|"
	routeMultiPlayerLANMenu            = routeMultiPlayerLANPrefix + "menu"
	routeMultiPlayerLANHostMenu        = routeMultiPlayerLANPrefix + "host-menu"
	routeMultiPlayerLANHostPresetN     = routeMultiPlayerLANPrefix + "host-preset:"
	routeMultiPlayerLANHostCustomSetup = routeMultiPlayerLANPrefix + "host-custom-setup"
	routeMultiPlayerLANHostLobby       = routeMultiPlayerLANPrefix + "host-lobby"
	routeMultiPlayerLANJoinListen      = routeMultiPlayerLANPrefix + "join-listen"
	routeMultiPlayerLANJoinLobby       = routeMultiPlayerLANPrefix + "join-lobby"

	routeMultiPlayerDirectIPPrefix          = "mp-direct-ip|"
	routeMultiPlayerDirectIPMenu            = routeMultiPlayerDirectIPPrefix + "menu"
	routeMultiPlayerDirectIPHostMenu        = routeMultiPlayerDirectIPPrefix + "host-menu"
	routeMultiPlayerDirectIPHostPresetN     = routeMultiPlayerDirectIPPrefix + "host-preset:"
	routeMultiPlayerDirectIPHostCustomSetup = routeMultiPlayerDirectIPPrefix + "host-custom-setup"
	routeMultiPlayerDirectIPHostEnterIP     = routeMultiPlayerDirectIPPrefix + "host-enter-ip"
	routeMultiPlayerDirectIPHostLobby       = routeMultiPlayerDirectIPPrefix + "host-lobby"
	routeMultiPlayerDirectIPJoinEnterIP     = routeMultiPlayerDirectIPPrefix + "join-enter-ip"
	routeMultiPlayerDirectIPJoinLobby       = routeMultiPlayerDirectIPPrefix + "join-lobby"

	routeMultiPlayerUDPHostGame = "udp-host-game"
	routeMultiPlayerUDPJoinGame = "udp-join-game"

	routeConfigPrefix            = "config|"
	routeConfigMenu              = routeConfigPrefix + "menu"
	routeConfigLocalPlayerSetupN = routeConfigPrefix + "local-player-setup:"

	routeConfigLanguage = routeConfigPrefix + "lang"

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
		return screen.Screens(
			[]screen.Screen{
				app.screener.Screen(ctx, app.title(ctx)),
				app.screener.Screen(ctx, app.menuMain(ctx)),
			}), ctx.Done()

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
	case id == routeSinglePlayerCustomSetup:
		data = app.menuSinglePlayerSetup(ctx, -1, routeSinglePlayerCustomGame)
	case id == routeSinglePlayerCustomGame:
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
				data = app.menuErrorText(ctx, i18n.Tf(i18n.KeyErrorTooManyPlayers, setup.MaxLocalPlayers))
			} else {
				data = app.gameMultiPlayerLocal(ctx)
			}
		}
	case id == routeMultiPlayerLocalCustomSetup:
		data = app.menuMultiPlayerSetup(ctx, setup.MaxLocalPlayers, -1, routeMultiPlayerLocalCustomGame)
	case id == routeMultiPlayerLocalCustomGame:
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
	case id == routeMultiPlayerLANHostCustomSetup:
		data = app.menuMultiPlayerSetup(ctx, setup.MaxPlayers, -1, routeMultiPlayerLANHostLobby)
	case id == routeMultiPlayerLANHostLobby:
		data = app.menuMultiPlayerLANHostLobby(ctx)

	case id == routeMultiPlayerLANJoinListen:
		data = app.menuMultiPlayerLANJoinListen(ctx)
	case id == routeMultiPlayerLANJoinLobby:
		data = app.menuMultiPlayerLANJoinLobby(ctx)

	// Multi-player Direct IP

	case id == routeMultiPlayerDirectIPMenu:
		data = app.menuMultiPlayerDirectIPMain(ctx)
	case id == routeMultiPlayerDirectIPHostMenu:
		data = app.menuMultiPlayerDirectIPHostMenu(ctx)
	case strings.HasPrefix(string(id), routeMultiPlayerDirectIPHostPresetN):
		s := strings.TrimPrefix(string(id), routeMultiPlayerDirectIPHostPresetN)
		idx, err := strconv.Atoi(s)
		if err != nil {
			data = app.menuError(ctx, err)
		} else {
			app.loadPresetMulti(idx)
			if app.resultSetup != nil && app.resultSetup.PlayerCount() > 2 {
				data = app.menuErrorText(ctx, i18n.Tf(i18n.KeyErrorTooManyPlayers, 2))
			} else {
				data = app.menuMultiPlayerDirectIPHostEnterIP(ctx)
			}
		}
	case id == routeMultiPlayerDirectIPHostCustomSetup:
		data = app.menuMultiPlayerSetup(ctx, 2, -1, routeMultiPlayerDirectIPHostEnterIP)
	case id == routeMultiPlayerDirectIPHostEnterIP:
		data = app.menuMultiPlayerDirectIPHostEnterIP(ctx)
	case id == routeMultiPlayerDirectIPHostLobby:
		data = app.menuMultiPlayerDirectIPHostLobby(ctx)

	case id == routeMultiPlayerDirectIPJoinEnterIP:
		data = app.menuMultiPlayerDirectIPJoinEnterIP(ctx)
	case id == routeMultiPlayerDirectIPJoinLobby:
		data = app.menuMultiPlayerDirectIPJoinLobby(ctx)

	// Game

	case id == routeMultiPlayerUDPHostGame:
		data = app.gameUDPServer(ctx)
	case id == routeMultiPlayerUDPJoinGame:
		data = app.gameUDPClient(ctx)

	// Configure

	case id == routeConfigMenu:
		data = app.menuConfig(ctx)
	case strings.HasPrefix(string(id), routeConfigLocalPlayerSetupN):
		s := strings.TrimPrefix(string(id), routeConfigLocalPlayerSetupN)
		idx, err := strconv.Atoi(s)
		if err != nil {
			data = app.menuError(ctx, err)
		} else {
			data = app.menuConfigLocalPlayer(ctx, idx)
		}
	case id == routeConfigLanguage:
		data = app.menuConfigLanguage(ctx)
	case id == routeConfigVideoSetup:
		data = app.menuConfigVideo(ctx)

	// About

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

func (app *App) saveConfig() {
	_ = config.Save(app.logger, app.cfgPath, &app.cfg)
}

func latencyState(s udpstar.ClientState) string {
	switch s {
	case udpstar.ClientStateNew:
		return i18n.T(i18n.KeyFieldLatencyStateNew)
	case udpstar.ClientStateLocal:
		return i18n.T(i18n.KeyFieldLatencyStateLocal)
	case udpstar.ClientStateGood:
		return i18n.T(i18n.KeyFieldLatencyStateGood)
	case udpstar.ClientStateLagging:
		return i18n.T(i18n.KeyFieldLatencyStateLagging)
	case udpstar.ClientStateLost:
		return i18n.T(i18n.KeyFieldLatencyStateLost)
	default:
		return "?"
	}
}

func latenciesToString(l []udpstar.LatencyActor, names []string) string {
	if len(l) == 0 {
		return ""
	}
	sb := strings.Builder{}

	sb.WriteString(i18n.T(i18n.KeyFieldLatency))
	sb.WriteString(":\n")
	for i, v := range l {
		sb.WriteString(fmt.Sprintf("%d. %s [%s] %dms\n",
			i+1, names[i], latencyState(v.State), v.Latency.Milliseconds()))
	}
	return sb.String()
}
