// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package app

import (
	"net"

	"github.com/marko-gacesa/udpstar/udp"
)

type udpSender struct {
	addr net.UDPAddr
	srv  *udp.Service
}

func (s udpSender) Send(data []byte) error {
	return s.srv.Send(data, s.addr)
}
