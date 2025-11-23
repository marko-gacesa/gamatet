// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package latency

import (
	"strconv"
	"strings"
	"time"

	"github.com/marko-gacesa/gamatet/logic/cache"
)

type Latency struct {
	cached *cache.String[time.Duration]
}

func NewLatency(fn func() time.Duration) *Latency {
	return &Latency{
		cached: cache.NewString[time.Duration](
			fn, func(prev *time.Duration, curr time.Duration) bool {
				equal := (*prev).Milliseconds() == curr.Milliseconds()
				if !equal {
					*prev = curr
				}
				return equal
			},
			func(l time.Duration) string {
				sb := strings.Builder{}
				sb.WriteString("latency=")
				sb.WriteString(strconv.FormatInt(l.Milliseconds(), 10))
				sb.WriteString("ms")
				return sb.String()
			},
			time.Second,
		),
	}
}

func (l *Latency) String() string { return l.cached.String() }

type LQ struct {
	cached *cache.String[LQValue]
}

type LQValue struct {
	Latency time.Duration
	Quality time.Duration
}

func NewLQCache(fn func() LQValue) *LQ {
	return &LQ{
		cached: cache.NewString[LQValue](
			fn, func(l1 *LQValue, l2 LQValue) bool {
				equal := *l1 == l2
				if !equal {
					*l1 = l2
				}
				return equal
			},
			func(l LQValue) string {
				sb := strings.Builder{}
				sb.WriteString("L=")
				sb.WriteString(strconv.FormatInt(l.Latency.Milliseconds(), 10))
				sb.WriteString("ms Q=")
				sb.WriteString(strconv.FormatInt(l.Quality.Milliseconds(), 10))
				sb.WriteString("ms")
				return sb.String()
			},
			time.Second,
		),
	}
}

func (l *LQ) String() string { return l.cached.String() }
