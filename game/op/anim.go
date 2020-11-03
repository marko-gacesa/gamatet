// Copyright (c) 2020 by Marko Gaćeša

package op

import (
	"gamatet/game/piece"
	"gamatet/logic/anim"
	"math"
	"time"
)

func animateNewPiece(ctrl *piece.Ctrl, isAnim bool) {
	if !isAnim {
		return
	}

	t := time.Now()

	ctrl.List.Add(anim.NewPopIn(t, piece.DurationAnimNewPiece))
	ctrl.List.Add(anim.NewFall(t, piece.DurationAnimNewPiece, 1))
}

func animateMovePiece(ctrl *piece.Ctrl, dx, dy int, isAnim bool) {
	if !isAnim {
		return
	}

	t := time.Now()

	if dx != 0 {
		ctrl.List.Add(anim.NewXQuad(t, piece.DurationMove, float32(dx)))
	}
	if dy != 0 {
		ctrl.List.Add(anim.NewYQuad(t, piece.DurationMove, float32(dy)))
	}
}

func animateRotatePiece(ctrl *piece.Ctrl, cw, inverted, isAnim bool) {
	if !isAnim {
		return
	}

	if cw {
		ctrl.List.Add(anim.NewZRotQuad(time.Now(), piece.DurationRotate, -math.Pi/2))
		if inverted {
			ctrl.List.Add(anim.NewYQuad(time.Now(), piece.DurationRotate, -1))
		}
	} else {
		ctrl.List.Add(anim.NewZRotQuad(time.Now(), piece.DurationRotate, math.Pi/2))
		if inverted {
			ctrl.List.Add(anim.NewXQuad(time.Now(), piece.DurationRotate, 1))
		}
	}
}

func animateDropPiece(ctrl *piece.Ctrl, height int, isAnim bool) {
	if !isAnim {
		return
	}

	duration := piece.GetFallDuration(height)
	ctrl.List.Add(anim.NewFall(time.Now(), duration, float32(height)))
}

func animateSlidePiece(ctrl *piece.Ctrl, isAnim bool) {
	if !isAnim {
		return
	}

	ctrl.List.Add(anim.NewSlide(time.Now(), piece.GetSlideDuration(ctrl.Level)))
}

func animateBlinkPiece(ctrl *piece.Ctrl, isAnim bool) {
	if !isAnim {
		return
	}

	// TODO: Fix flash animation
	ctrl.List.Add(anim.NewFlash(time.Now(), 400*time.Millisecond))
}
