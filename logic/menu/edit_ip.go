// Copyright (c) 2024, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package menu

import (
	"net"
	"slices"
)

var _ Item = (*IP)(nil)

// IP is menu item that assigns a value to a net.IP variable, parsed as an IP address.
type IP struct {
	textBase
	ptr *string
}

// NewIP creates new IP menu item.
func NewIP(ptr *string, label, description string, options ...func(Item)) *IP {
	if ptr == nil {
		panic(strNilPointer)
	}
	ip := &IP{
		textBase: makeTextBase(20, 20, label, description),
		ptr:      ptr,
	}
	ip.textBase.converter = ip
	ip.fix()
	applyOptions(ip, options...)
	return ip
}

func (ip *IP) fix() {
	addr, err := net.ResolveIPAddr("ip", *ip.ptr)
	if err != nil || ip.isInvalid(addr.IP) {
		*ip.ptr = ""
		return
	}

	*ip.ptr = addr.String()
}

func (ip *IP) getValueAsStr() string {
	return *ip.ptr
}

func (ip *IP) setValueFromStr(s string) {
	*ip.ptr = s
	ip.fix()
}

func (*IP) allowed(r rune) bool {
	return r > 32 && r < 127
}

func (ip *IP) allowedInsert(r rune, s []rune, cursor int) bool {
	hasPercent := slices.Contains(s, '%')
	if hasPercent {
		return r > 32 && r < 127
	}
	return r == '.' || r == ':' || r == '%' || r >= '0' && r <= '9' || r >= 'A' && r <= 'F' || r >= 'a' && r <= 'f'
}

func (*IP) isInvalid(a net.IP) bool {
	invalid := a == nil || a.IsUnspecified() ||
		a.IsMulticast() || a.IsLinkLocalMulticast() || a.IsLinkLocalUnicast() ||
		a.IsLoopback()
	return invalid
}
