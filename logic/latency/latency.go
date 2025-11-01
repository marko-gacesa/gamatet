// Copyright (c) 2025 by Marko Gaćeša

package latency

import (
	"gamatet/logic/cache"
	"strconv"
	"strings"
	"time"
)

type Latency struct {
	cached *cache.String[time.Duration]
}

func NewLatency(fn func() time.Duration) *Latency {
	return &Latency{
		cached: cache.NewString[time.Duration](
			fn, func(l1 time.Duration, l2 time.Duration) bool {
				return l1.Milliseconds() == l2.Milliseconds()
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
			fn, func(l1 LQValue, l2 LQValue) bool {
				return l1 == l2
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
