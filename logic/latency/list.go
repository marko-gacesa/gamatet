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
			fn, func(prev *[]udpstar.LatencyActor, curr []udpstar.LatencyActor) bool {
				equal := slices.EqualFunc(*prev, curr, func(l1 udpstar.LatencyActor, l2 udpstar.LatencyActor) bool {
					return l1.State == l2.State && l1.Latency.Milliseconds() == l2.Latency.Milliseconds()
				})
				if !equal {
					if len(*prev) != len(curr) {
						*prev = make([]udpstar.LatencyActor, len(curr))
					}
					for i := range curr {
						(*prev)[i] = curr[i]
					}
				}
				return equal
			},
			func(l []udpstar.LatencyActor) string {
				if len(l) == 0 {
					return ""
				}
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
