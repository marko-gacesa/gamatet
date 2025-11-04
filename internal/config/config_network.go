// Copyright (c) 2024, 2025 by Marko Gaćeša

package config

import (
	"fmt"
	"net"
)

const (
	//defaultMulticast = "239.255.231.79"
	defaultMulticast     = "224.0.0.79"
	defaultPort          = 64774
	defaultMulticastPort = 64775
)

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
