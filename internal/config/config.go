// Copyright (c) 2024 by Marko Gaćeša

package config

import (
	"encoding/json"
	"fmt"
	"gamatet/game/piece"
	"gamatet/game/setup"
	"net"
	"os"
	"path"
	"slices"
)

const (
	filename = ".gamatet.config.json"

	defaultLanguage = "en"

	//defaultMulticast = "239.255.231.79"
	defaultMulticast     = "224.0.0.79"
	defaultPort          = 64774
	defaultMulticastPort = 64775

	LocalPlayerLimit = 4
	LocalGameSetups  = 4
)

type Config struct {
	Language string `json:"language"`

	PlayerInfos []PlayerInfo `json:"player_names"`

	Network Network `json:"network"`

	LocalGameDefaults []setup.Setup `json:"local_game_defaults"`
	LANGameDefaults   setup.Setup   `json:"lan_game_defaults"`
}

type PlayerInfo struct {
	Name string `json:"name"`
	PlayerConfig
}

type PlayerConfig struct {
	RotationDirectionCW bool `json:"rotation_direction_cw"`
	SlideDisabled       bool `json:"slide_disabled"`
	WallKick            int  `json:"wall_kick"`
}

func (cfg PlayerConfig) Serialize() []byte {
	return setup.Pack((*setup.PlayerConfig)(&cfg))
}

func (cfg *Config) Sanitize() {
	if cfg.Language == "" {
		cfg.Language = defaultLanguage
	}

	if len(cfg.PlayerInfos) < LocalPlayerLimit {
		for i := len(cfg.PlayerInfos); i < LocalPlayerLimit; i++ {
			cfg.PlayerInfos = append(cfg.PlayerInfos, PlayerInfo{
				Name: "",
				PlayerConfig: PlayerConfig{
					RotationDirectionCW: false,
					SlideDisabled:       false,
					WallKick:            2,
				},
			})
		}
	} else if len(cfg.PlayerInfos) > LocalPlayerLimit {
		cfg.PlayerInfos = cfg.PlayerInfos[:LocalPlayerLimit]
	}

	cfg.PlayerInfos = slices.Clip(cfg.PlayerInfos)
	for i := range cfg.PlayerInfos {
		if cfg.PlayerInfos[i].Name == "" {
			cfg.PlayerInfos[i].Name = fmt.Sprintf("Player %d", i+1)
		}
		if cfg.PlayerInfos[i].WallKick < 0 || cfg.PlayerInfos[i].WallKick > piece.MaxWallKick {
			cfg.PlayerInfos[i].WallKick = piece.MaxWallKick
		}
	}

	cfg.Network.Sanitize()

	cfg.LANGameDefaults.Sanitize()

	if len(cfg.LocalGameDefaults) < LocalGameSetups {
		for i := len(cfg.LocalGameDefaults); i < LocalGameSetups; i++ {
			cfg.LocalGameDefaults = append(cfg.LocalGameDefaults, setup.MultiplayerPlayerSetupDefault())
		}
	} else if len(cfg.LocalGameDefaults) > LocalGameSetups {
		cfg.LocalGameDefaults = cfg.LocalGameDefaults[:LocalGameSetups]
	}
	for i := range cfg.LocalGameDefaults {
		cfg.LocalGameDefaults[i].Sanitize()
	}
}

type Network struct {
	Port             int    `json:"port"`
	MulticastPort    int    `json:"multicast_port"`
	MulticastAddress string `json:"multicast_address"`
}

func (cfg *Network) Sanitize() {
	if cfg.Port == 0 || cfg.Port > 65535 {
		cfg.Port = defaultPort
	}

	if cfg.MulticastPort == 0 || cfg.MulticastPort > 65535 {
		cfg.MulticastPort = defaultMulticastPort
	}

	if cfg.MulticastAddress != "" {
		a, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", cfg.MulticastAddress, cfg.MulticastPort))
		if err != nil || !a.IP.IsLinkLocalMulticast() {
			cfg.MulticastAddress = ""
		}

	}

	if cfg.MulticastAddress == "" {
		cfg.MulticastAddress = defaultMulticast
	}
}

func (cfg *Network) GetMulticastAddress() net.UDPAddr {
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", cfg.MulticastAddress, cfg.MulticastPort))
	return *addr
}

func Load() (Config, string) {
	var cfg Config

	dirs := getDirList()
	for _, dir := range dirs {
		fn := path.Join(dir, filename)
		f, err := os.Open(fn)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			fmt.Printf("failed to open config file %s: %s\n", fn, err.Error())
			continue
		}

		if err := func() error {
			defer f.Close()
			return json.NewDecoder(f).Decode(&cfg)
		}(); err != nil {
			fmt.Printf("failed to load config from %s: %s\n", fn, err.Error())
		}

		return cfg, dir
	}

	cfg.Sanitize()

	if len(dirs) > 0 {
		return cfg, path.Join(dirs[0], filename)
	}

	return cfg, filename
}

func getDirList() []string {
	var dirs []string

	dir, _ := os.UserHomeDir()
	if dir != "" {
		dirs = append(dirs, dir)
	}

	if len(os.Args) > 0 {
		dir = path.Dir(os.Args[0])
		if dir != "" {
			dirs = append(dirs, dir)
		}
	}

	return dirs
}
