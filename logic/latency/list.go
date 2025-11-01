// Copyright (c) 2025 by Marko Gaćeša

package latency

import (
	"fmt"
	"gamatet/logic/cache"
	"github.com/marko-gacesa/udpstar/udpstar"
	"slices"
	"strings"
	"time"
)

type List struct {
	cached *cache.String[[]udpstar.LatencyActor]
}

func NewList(fn func() []udpstar.LatencyActor) *List {
	return &List{
		cached: cache.NewString[[]udpstar.LatencyActor](
			fn, func(l1 []udpstar.LatencyActor, l2 []udpstar.LatencyActor) bool {
				return slices.EqualFunc(l1, l2, func(l1 udpstar.LatencyActor, l2 udpstar.LatencyActor) bool {
					return l1.State == l2.State && l1.Latency == l2.Latency
				})
			},
			func(l []udpstar.LatencyActor) string {
				sb := strings.Builder{}
				sb.WriteString("Latencies:\n")
				for i, v := range l {
					sb.WriteString(fmt.Sprintf("%d. %s [%s] %dms\n",
						i+1, v.Name, clientState(v.State), v.Latency.Milliseconds()))
				}
				return sb.String()
			},
			time.Second,
		),
	}
}

func (l *List) String() string { return l.cached.String() }

func clientState(s udpstar.ClientState) string {
	switch s {
	case udpstar.ClientStateNew:
		return "new"
	case udpstar.ClientStateLocal:
		return "local"
	case udpstar.ClientStateGood:
		return "good"
	case udpstar.ClientStateLagging:
		return "slow"
	case udpstar.ClientStateLost:
		return "LOST"
	default:
		return "?"
	}
}
