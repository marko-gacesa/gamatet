// Copyright (c) 2020-2024 by Marko Gaćeša

package anim

import (
	"math"
	"time"
)

// Translation X, Y, Z - linear and quadratic

type animTransLin struct {
	animBase
	dx, dy, dz float32
}

func NewTransLin(now time.Time, duration time.Duration, dx, dy, dz float32) Anim {
	return &animTransLin{animBase: animBase{startedAt: now, duration: duration}, dx: -dx, dy: -dy, dz: -dz}
}

func (*animTransLin) Feature() Feature { return Translate }

func (a *animTransLin) Translate() (dx, dy, dz float32) {
	rt := 1.0 - a.T()
	dx = a.dx * rt
	dy = a.dy * rt
	dz = a.dz * rt
	return
}

func (*animTransLin) Rotate() (rx, ry, rz float32) { return }
func (*animTransLin) Scale() (sx, sy, sz float32)  { return }
func (*animTransLin) Color() (r, g, b, a float32)  { return }

type animTransQuad struct {
	animBase
	dx, dy, dz float32
}

func NewTransQuad(now time.Time, duration time.Duration, dx, dy, dz float32) Anim {
	return &animTransQuad{animBase: animBase{startedAt: now, duration: duration}, dx: -dx, dy: -dy, dz: -dz}
}

func (*animTransQuad) Feature() Feature { return Translate }

func (a *animTransQuad) Translate() (dx, dy, dz float32) {
	rt := 1.0 - a.T()
	rt = rt * rt
	dx = a.dx * rt
	dy = a.dy * rt
	dz = a.dz * rt
	return
}

func (*animTransQuad) Rotate() (rx, ry, rz float32) { return }
func (*animTransQuad) Scale() (sx, sy, sz float32)  { return }
func (*animTransQuad) Color() (r, g, b, a float32)  { return }

// Translation X - linear and quadratic

type animXLin struct {
	animBase
	dx float32
}

func NewXLin(now time.Time, duration time.Duration, dx float32) Anim {
	return &animTransLin{animBase: animBase{startedAt: now, duration: duration}, dx: -dx}
}

func (*animXLin) Feature() Feature { return Translate }

func (a *animXLin) Translate() (dx, dy, dz float32) {
	rt := 1.0 - a.T()
	dx = a.dx * rt
	return
}

func (*animXLin) Rotate() (rx, ry, rz float32) { return }
func (*animXLin) Scale() (sx, sy, sz float32)  { return }
func (*animXLin) Color() (r, g, b, a float32)  { return }

type animXQuad struct {
	animBase
	dx float32
}

func NewXQuad(now time.Time, duration time.Duration, dx float32) Anim {
	return &animXQuad{animBase: animBase{startedAt: now, duration: duration}, dx: -dx}
}

func (*animXQuad) Feature() Feature { return Translate }

func (a *animXQuad) Translate() (dx, dy, dz float32) {
	rt := 1.0 - a.T()
	dx = a.dx * rt * rt
	return
}

func (*animXQuad) Rotate() (rx, ry, rz float32) { return }
func (*animXQuad) Scale() (sx, sy, sz float32)  { return }
func (*animXQuad) Color() (r, g, b, a float32)  { return }

// Translation Y - linear and quadratic

type animYLin struct {
	animBase
	dy float32
}

func NewYLin(now time.Time, duration time.Duration, dy float32) Anim {
	return &animYLin{animBase: animBase{startedAt: now, duration: duration}, dy: -dy}
}

func (*animYLin) Feature() Feature { return Translate }

func (a *animYLin) Translate() (dx, dy, dz float32) {
	rt := 1.0 - a.T()
	dy = a.dy * rt
	return
}

func (*animYLin) Rotate() (rx, ry, rz float32) { return }
func (*animYLin) Scale() (sx, sy, sz float32)  { return }
func (*animYLin) Color() (r, g, b, a float32)  { return }

type animYQuad struct {
	animBase
	dy float32
}

func NewYQuad(now time.Time, duration time.Duration, dy float32) Anim {
	return &animYQuad{animBase: animBase{startedAt: now, duration: duration}, dy: -dy}
}

func (*animYQuad) Feature() Feature { return Translate }

func (a *animYQuad) Translate() (dx, dy, dz float32) {
	rt := 1.0 - a.T()
	dy = a.dy * rt * rt
	return
}

func (*animYQuad) Rotate() (rx, ry, rz float32) { return }
func (*animYQuad) Scale() (sx, sy, sz float32)  { return }
func (*animYQuad) Color() (r, g, b, a float32)  { return }

// Fall

type animFall struct {
	animBase
	dy float32
}

func NewFall(now time.Time, duration time.Duration, height float32) Anim {
	return &animFall{animBase: animBase{startedAt: now, duration: duration}, dy: height}
}

func (*animFall) Feature() Feature { return Translate }

func (a *animFall) Translate() (dx, dy, dz float32) {
	rt := a.T()
	dy = a.dy * (1.0 - rt*rt)
	return
}

func (*animFall) Rotate() (rx, ry, rz float32) { return }
func (*animFall) Scale() (sx, sy, sz float32)  { return }
func (*animFall) Color() (r, g, b, a float32)  { return }

// Quake

type animQuake struct {
	animBase
	dur float32
	mag float32
}

func NewQuake(now time.Time, intensity byte) Anim {
	intf := float32(intensity)
	dur := 1 / intf
	mag := 0.08 * intf

	d := time.Duration(50 * float32(time.Millisecond) * intf)

	return &animQuake{animBase: animBase{startedAt: now, duration: d}, dur: dur, mag: mag}
}

func (*animQuake) Feature() Feature { return Translate }

func (a *animQuake) Translate() (dx, dy, dz float32) {
	//                                 1
	// 0.08 * intensity * 2 * ( --------------- - 0.5 )^2 * SIN(...)
	//                          1 + t/intensity
	f := 1 / (1 + a.t*a.dur)
	amp := a.mag * 2 * (f*f - 0.5)
	dx = amp * float32(math.Sin(31*float64(a.t)))
	dy = amp * float32(math.Sin(27*float64(a.t)))
	return
}

func (*animQuake) Rotate() (rx, ry, rz float32) { return }
func (*animQuake) Scale() (sx, sy, sz float32)  { return }
func (*animQuake) Color() (r, g, b, a float32)  { return }

// Rotation Z - linear and quadratic

type animZRotLin struct {
	animBase
	rz float32
}

func NewZRotLin(now time.Time, duration time.Duration, rz float32) Anim {
	return &animZRotLin{animBase: animBase{startedAt: now, duration: duration}, rz: -rz}
}

func (*animZRotLin) Feature() Feature { return Rotate }

func (a *animZRotLin) Rotate() (rx, ry, rz float32) {
	rt := 1.0 - a.T()
	rz = a.rz * rt
	return
}

func (*animZRotLin) Translate() (dx, dy, dz float32) { return }
func (*animZRotLin) Scale() (sx, sy, sz float32)     { return }
func (*animZRotLin) Color() (r, g, b, a float32)     { return }

type animZRotQuad struct {
	animBase
	rz float32
}

func NewZRotQuad(now time.Time, duration time.Duration, rz float32) Anim {
	return &animZRotQuad{animBase: animBase{startedAt: now, duration: duration}, rz: -rz}
}

func (*animZRotQuad) Feature() Feature { return Rotate }

func (a *animZRotQuad) Rotate() (rx, ry, rz float32) {
	rt := 1.0 - a.T()
	rz = a.rz * rt * rt
	return
}

func (*animZRotQuad) Translate() (dx, dy, dz float32) { return }
func (*animZRotQuad) Scale() (sx, sy, sz float32)     { return }
func (*animZRotQuad) Color() (r, g, b, a float32)     { return }

// PopIn

type animPopIn struct {
	animBase
}

func NewPopIn(now time.Time, duration time.Duration) Anim {
	return &animPopIn{animBase{startedAt: now, duration: duration}}
}

func (*animPopIn) Feature() Feature { return Scale }

func (a *animPopIn) Scale() (sx, sy, sz float32) {
	t := 1 - a.T()
	t = 1 - (t * t)
	return t, t, t
}

func (*animPopIn) Translate() (dx, dy, dz float32) { return }
func (*animPopIn) Rotate() (rx, ry, rz float32)    { return }
func (*animPopIn) Color() (r, g, b, a float32)     { return }

// PopOut

type animPopOut struct {
	animBase
}

func NewPopOut(now time.Time, duration time.Duration) Anim {
	return &animPopOut{animBase{startedAt: now, duration: duration}}
}

func (*animPopOut) Feature() Feature { return Scale }

func (a *animPopOut) Scale() (sx, sy, sz float32) {
	t := 1 - a.T()
	t = t * t
	return t, t, t
}

func (*animPopOut) Translate() (dx, dy, dz float32) { return }
func (*animPopOut) Rotate() (rx, ry, rz float32)    { return }
func (*animPopOut) Color() (r, g, b, a float32)     { return }

// Spin

type animSpin struct {
	animCyclic
}

func NewSpin(now time.Time, period time.Duration) Anim {
	return &animSpin{animCyclic{startedAt: now, period: period}}
}

func (*animSpin) Feature() Feature { return Rotate }

func (a *animSpin) Rotate() (rx, ry, rz float32) {
	t := a.T()
	t *= 2 * math.Pi
	rx = 5.1 * t
	rz = 1.7 * t
	return
}

func (*animSpin) Translate() (dx, dy, dz float32) { return }
func (*animSpin) Scale() (sx, sy, sz float32)     { return }
func (*animSpin) Color() (r, g, b, a float32)     { return }

// SpinOnce

type animSpinOnce struct {
	animBase
}

func NewSpinOnce(now time.Time, duration time.Duration) Anim {
	return &animSpinOnce{animBase{startedAt: now, duration: duration}}
}

func (*animSpinOnce) Feature() Feature { return Rotate }

func (a *animSpinOnce) Rotate() (rx, ry, rz float32) {
	t := float64(a.T())
	t = t * t * 2 * math.Pi
	rx = float32(math.Sin(2 * t))
	ry = float32(math.Sin(4 * t))
	rz = float32(math.Sin(t))
	return
}

func (*animSpinOnce) Translate() (dx, dy, dz float32) { return }
func (*animSpinOnce) Scale() (sx, sy, sz float32)     { return }
func (*animSpinOnce) Color() (r, g, b, a float32)     { return }

// RotateZ

type animRotateZ struct {
	animCyclic
}

func NewRotateZ(now time.Time, period time.Duration) Anim {
	return &animRotateZ{animCyclic{startedAt: now, period: period}}
}

func (*animRotateZ) Feature() Feature { return Rotate }

func (a *animRotateZ) Rotate() (rx, ry, rz float32) {
	return 0, 0, 2 * math.Pi * a.T()
}

func (*animRotateZ) Translate() (dx, dy, dz float32) { return }
func (*animRotateZ) Scale() (sx, sy, sz float32)     { return }
func (*animRotateZ) Color() (r, g, b, a float32)     { return }

// ColorTrans

type animColorTrans struct {
	r0, g0, b0 float32
	r1, g1, b1 float32
	animBase
}

func NewColorTrans(now time.Time, duration time.Duration, color0, color1 uint32) Anim {
	return &animColorTrans{
		r0:       float32(color0>>24) / 255,
		g0:       float32(color0>>16&0xFF) / 255,
		b0:       float32(color0>>8&0xFF) / 255,
		r1:       float32(color1>>24) / 255,
		g1:       float32(color1>>16&0xFF) / 255,
		b1:       float32(color1>>8&0xFF) / 255,
		animBase: animBase{startedAt: now, duration: duration},
	}
}

func (*animColorTrans) Feature() Feature { return Color }

func (q *animColorTrans) Color() (r, g, b, a float32) {
	t := q.T()
	r = (1-t)*q.r0 + t*q.r1
	g = (1-t)*q.g0 + t*q.g1
	b = (1-t)*q.b0 + t*q.b1
	a = 1
	return
}

func (*animColorTrans) Translate() (dx, dy, dz float32) { return }
func (*animColorTrans) Scale() (sx, sy, sz float32)     { return }
func (*animColorTrans) Rotate() (rx, ry, rz float32)    { return }

// Flash

type animFlash struct {
	animBase
}

func NewFlash(now time.Time, duration time.Duration) Anim {
	return &animFlash{animBase{startedAt: now, duration: duration}}
}

func (*animFlash) Feature() Feature { return Color }

func (*animFlash) Scale() (sx, sy, sz float32)     { return }
func (*animFlash) Translate() (dx, dy, dz float32) { return }
func (*animFlash) Rotate() (rx, ry, rz float32)    { return }

func (q *animFlash) Color() (r, g, b, a float32) {
	i := float32(0.7 + 0.9*math.Sin(float64(q.T()*8*math.Pi)))
	r = i
	g = i
	b = i
	a = 1.0
	return
}

// Slide

type animSlide struct {
	animBase
}

func NewSlide(now time.Time, duration time.Duration) Anim {
	return &animSlide{animBase{startedAt: now, duration: duration}}
}

func (*animSlide) Feature() Feature { return Color }

func (*animSlide) Scale() (sx, sy, sz float32)     { return }
func (*animSlide) Translate() (dx, dy, dz float32) { return }
func (*animSlide) Rotate() (rx, ry, rz float32)    { return }

func (q *animSlide) Color() (r, g, b, a float32) {
	t := 1.0 + 0.4*float32(math.Sin(float64(q.T())*4*math.Pi))
	return t, t, t, 1
}

// Meld

const meldColorR = 0.4
const meldColorG = 1.5
const meldColorB = 1.2

type animMeld struct {
	animBase
}

func NewMeld(now time.Time, duration time.Duration) Anim {
	return &animMeld{animBase{startedAt: now, duration: duration}}
}

func (*animMeld) Feature() Feature { return Color }

func (*animMeld) Scale() (sx, sy, sz float32)     { return }
func (*animMeld) Translate() (dx, dy, dz float32) { return }
func (*animMeld) Rotate() (rx, ry, rz float32)    { return }

func (q *animMeld) Color() (r, g, b, a float32) {
	t := 1 - q.T()
	t = t * t
	r = meldColorR*t + (1 - t)
	g = meldColorG*t + (1 - t)
	b = meldColorB*t + (1 - t)
	a = 1.0
	return
}

// Rainbow

type animRainbow struct {
	animCyclic
}

func NewRainbow(now time.Time, period time.Duration) Anim {
	return &animRainbow{animCyclic{startedAt: now, period: period}}
}

func (*animRainbow) Feature() Feature { return Color }

func (*animRainbow) Scale() (sx, sy, sz float32)     { return }
func (*animRainbow) Translate() (dx, dy, dz float32) { return }
func (*animRainbow) Rotate() (rx, ry, rz float32)    { return }

func (q *animRainbow) Color() (r, g, b, a float32) {
	t := float64(q.T()) * 2 * math.Pi
	r = float32(math.Sin(t) + 0.5)
	g = float32(math.Sin(t+2*math.Pi/3) + 0.5)
	b = float32(math.Sin(t+4*math.Pi/3) + 0.5)
	a = 1.0
	return
}

// Pulse

type animPulse struct {
	animCyclic
}

func NewPulse(now time.Time, period time.Duration) Anim {
	return &animPulse{animCyclic{startedAt: now, period: period}}
}

func (*animPulse) Feature() Feature { return Scale }

func (q *animPulse) Scale() (sx, sy, sz float32) {
	a := float32(math.Sin(3 * float64(q.T())))
	a = 0.5 + 0.5*a*a
	return a, a, a
}

func (*animPulse) Translate() (dx, dy, dz float32) { return }
func (*animPulse) Rotate() (rx, ry, rz float32)    { return }
func (*animPulse) Color() (r, g, b, a float32)     { return }
