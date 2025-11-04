// Copyright (c) 2024, 2025 by Marko Gaćeša

package app

import (
	"context"
	"gamatet/game/setup"
	"gamatet/internal/config"
	"gamatet/logic/screen"
	"github.com/marko-gacesa/udpstar/udp"
	"github.com/marko-gacesa/udpstar/udpstar"
	"github.com/marko-gacesa/udpstar/udpstar/client"
	"github.com/marko-gacesa/udpstar/udpstar/message"
	"github.com/marko-gacesa/udpstar/udpstar/server"
	"log/slog"
	"net"
	"sync"
	"time"
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
	return app.cfg.LocalPlayers.Infos[i].PlayerConfig
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
	routeMain route = "main"
	routeQuit route = "quit"
	routeBack route = "back"

	routeMenuSinglePlayer route = "menu-single-player"

	routeMenuLocalCreate route = "menu-local-create"

	routeMenuLANMain         route = "menu-lan-main"
	routeMenuLANServerCreate route = "menu-lan-server-create"
	routeMenuLANServerLobby  route = "menu-lan-server-lobby"
	routeMenuLANClientJoin   route = "menu-lan-client-join"
	routeMenuLANClientLobby  route = "menu-lan-client-lobby"

	routeGame              route = "game"
	routeGameSinglePreset1 route = "game-single-preset-1"
	routeGameSinglePreset2 route = "game-single-preset-2"
	routeGameSinglePreset3 route = "game-single-preset-3"
	routeGameSinglePreset4 route = "game-single-preset-4"
	routeGameSinglePreset5 route = "game-single-preset-5"

	routeGameUDPServer route = "game-udp-server"
	routeGameUDPClient route = "game-udp-client"
)

func (app *App) MakeScreen(parentCtx context.Context) (screen.Screen, <-chan struct{}) {
	id := app.screenIDHistory.curr()
	var data any

	ctx := screen.NewContext(parentCtx)

	switch id {
	case "", routeQuit:
		data = nil
	case routeMain:
		data = app.menuMain(ctx)
	case routeMenuSinglePlayer:
		data = app.menuSinglePlayer(ctx)

	case routeMenuLocalCreate:
		data = app.menuLocalCreateGame(ctx, 0)

	case routeMenuLANMain:
		data = app.menuLANMain(ctx)
	case routeMenuLANServerCreate:
		data = app.menuLANServerCreate(ctx)
	case routeMenuLANServerLobby:
		data = app.menuLANServerLobby(ctx)
	case routeMenuLANClientJoin:
		data = app.menuLANClientJoin(ctx)
	case routeMenuLANClientLobby:
		data = app.menuLANClientLobby(ctx)

	case routeGameSinglePreset1:
		app.resultSetup = &app.cfg.Presets.Single[0]
		data = app.gameOne(ctx)
	case routeGameSinglePreset2:
		app.resultSetup = &app.cfg.Presets.Single[1]
		data = app.gameOne(ctx)
	case routeGameSinglePreset3:
		app.resultSetup = &app.cfg.Presets.Single[2]
		data = app.gameOne(ctx)
	case routeGameSinglePreset4:
		app.resultSetup = &app.cfg.Presets.Single[3]
		data = app.gameOne(ctx)
	case routeGameSinglePreset5:
		app.resultSetup = &app.cfg.Presets.Single[4]
		data = app.gameOne(ctx)
	case routeGame:
		data = app.game(ctx)

	case routeGameUDPServer:
		data = app.gameUDPServer(ctx)
	case routeGameUDPClient:
		data = app.gameUDPClient(ctx)
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
