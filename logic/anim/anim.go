// Copyright (c) 2020, 2025 by Marko Gaćeša

package anim

import "time"

type Feature byte

const (
	Translate Feature = 1 << iota
	Rotate
	Scale
	Color
)

type Anim interface {
	setNext(Anim)
	next() Anim
	Update(now time.Time) (isDone bool)
	T() float32
	Feature() Feature
	Translate() (dx, dy, dz float32)
	Rotate() (rx, ry, rz float32)
	Scale() (sx, sy, sz float32)
	Color() (r, g, b, a float32)
}

type animBase struct {
	chain     Anim
	startedAt time.Time
	duration  time.Duration
	t         float32
}

func (a *animBase) setNext(next Anim) { a.chain = next }
func (a *animBase) next() Anim        { return a.chain }
func (a *animBase) T() float32        { return a.t }

func (*animBase) Scale() (sx, sy, sz float32)     { return 1, 1, 1 }
func (*animBase) Translate() (dx, dy, dz float32) { return }
func (*animBase) Rotate() (rx, ry, rz float32)    { return }
func (*animBase) Color() (r, g, b, a float32)     { return 1, 1, 1, 1 }

func (a *animBase) Update(now time.Time) (done bool) {
	elapsed := now.Sub(a.startedAt)
	if elapsed < 0.0 {
		return
	}
	if elapsed >= a.duration {
		a.t = 1.0
		done = true
		return
	}

	a.t = float32(float64(elapsed) / float64(a.duration))
	return
}

type animCyclic struct {
	chain     Anim
	startedAt time.Time
	period    time.Duration
	t         float32
}

func (a *animCyclic) setNext(next Anim) { a.chain = next }
func (a *animCyclic) next() Anim        { return a.chain }
func (a *animCyclic) T() float32        { return a.t }

func (*animCyclic) Scale() (sx, sy, sz float32)     { return 1, 1, 1 }
func (*animCyclic) Translate() (dx, dy, dz float32) { return }
func (*animCyclic) Rotate() (rx, ry, rz float32)    { return }
func (*animCyclic) Color() (r, g, b, a float32)     { return 1, 1, 1, 1 }

func (a *animCyclic) Update(now time.Time) bool {
	elapsed := now.Sub(a.startedAt)
	if elapsed < 0.0 {
		return false
	}

	a.t = float32(float64(elapsed%a.period) / float64(a.period))
	return false
}
