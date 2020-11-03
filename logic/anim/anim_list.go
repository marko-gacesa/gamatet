// Copyright (c) 2020 by Marko Gaćeša

package anim

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Feature    Feature
	DX, DY, DZ float32
	RX, RY, RZ float32
	SX, SY, SZ float32
	R, G, B, A float32
}

type List struct {
	first Anim
}

type Processor interface {
	Process(now time.Time) Result
}

func (l *List) Clear() {
	l.first = nil
}

func (l *List) Add(anim Anim) {
	if anim == nil {
		return
	}

	if l.first != nil {
		anim.setNext(l.first)
	}
	l.first = anim
}

func (l *List) AddAll(anims ...Anim) {
	for _, anim := range anims {
		if l.first != nil {
			anim.setNext(l.first)
		}
		l.first = anim
	}
}

func (l *List) GetAll() (anims []Anim) {
	for anim := l.first; anim != nil; anim = anim.next() {
		anims = append(anims, anim)
	}
	return
}

func (l *List) Process(now time.Time) (r Result) {
	r.SX = 1.0
	r.SY = 1.0
	r.SZ = 1.0
	r.R = 1.0
	r.G = 1.0
	r.B = 1.0
	r.A = 1.0

	var prev, curr Anim
	curr = l.first

	for curr != nil {
		next := curr.next()
		isDone := curr.Update(now)

		if isDone {
			if prev == nil {
				l.first = next
			} else {
				prev.setNext(next)
			}

			curr.setNext(nil) // remove curr anim from the list
			curr = next

			continue
		}

		f := curr.Feature()

		if f&Translate > 0 {
			_dx, _dy, _dz := curr.Translate()
			r.DX += _dx
			r.DY += _dy
			r.DZ += _dz
		}

		if f&Rotate > 0 {
			_rx, _ry, _rz := curr.Rotate()
			r.RX += _rx
			r.RY += _ry
			r.RZ += _rz
		}

		if f&Scale > 0 {
			_sx, _sy, _sz := curr.Scale()
			r.SX *= _sx
			r.SY *= _sy
			r.SZ *= _sz
		}

		if f&Color > 0 {
			_r, _g, _b, _a := curr.Color()
			r.R *= _r
			r.G *= _g
			r.B *= _b
			r.A *= _a
		}

		r.Feature |= f

		prev = curr
		curr = next
	}

	return
}

func (l *List) Count() (n int) {
	for curr := l.first; curr != nil; curr = curr.next() {
		n++
	}
	return
}

func (l *List) String() string {
	sb := strings.Builder{}

	sb.WriteByte('[')
	for curr := l.first; curr != nil; curr = curr.next() {
		if curr != l.first {
			sb.WriteString(", ")
		}

		sb.WriteString(reflect.TypeOf(curr).String())
		sb.WriteByte('=')
		sb.WriteString(strconv.FormatFloat(float64(curr.T()), 'f', 8, 32))

	}
	sb.WriteByte(']')

	return sb.String()
}
