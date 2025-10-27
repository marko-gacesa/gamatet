// Copyright (c) 2024, 2025 by Marko Gaćeša

package render

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/marko-gacesa/udpstar/udpstar"
	"slices"
	"strings"
	"time"
)

type Latencies struct {
	valueFn  func() []udpstar.LatencyActor
	valueOld []udpstar.LatencyActor
	valueStr string
}

func NewLatencies(fn func() []udpstar.LatencyActor) *Latencies {
	return &Latencies{valueFn: fn}
}

func (l *Latencies) Prepare() {
	if l.valueFn == nil {
		return
	}

	value := l.valueFn()
	for i := range value {
		value[i].Latency = value[i].Latency.Truncate(100 * time.Microsecond)
	}

	if len(l.valueOld) != len(value) {
		l.valueOld = make([]udpstar.LatencyActor, len(value))
	}

	equal := slices.EqualFunc(value, l.valueOld, func(l1 udpstar.LatencyActor, l2 udpstar.LatencyActor) bool {
		return l1.State == l2.State && l1.Latency == l2.Latency
	})

	if !equal {
		sb := strings.Builder{}
		for i, v := range value {
			l.valueOld[i] = value[i]

			sb.WriteString(fmt.Sprintf("%d. %s [%s] %s,",
				i+1, v.Name, v.State.String(), v.Latency.String()))
			sb.WriteString("\n")
		}
		l.valueStr = sb.String()
	}
}

func (l *Latencies) Render(r *Renderer, text *Text) {
	if l.valueStr == "" {
		return
	}

	//tw, th := text.Dim(l.valueStr)
	//_, _ = tw, th

	const contentW = 80
	const contentH = contentW * 9 / 16
	r.OrthogonalFull(contentW, contentH, contentW, contentH, 1)

	model := mgl32.Translate3D(float32(-contentW)/2, float32(contentH)/2-0.5, 0)

	text.String(r, model, mgl32.Vec4{0.5, 0.5, 0, 0.7}, l.valueStr)
}
