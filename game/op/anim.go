// Copyright (c) 2020, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package op

import (
	"math"
	"time"

	"github.com/marko-gacesa/gamatet/game/piece"
	"github.com/marko-gacesa/gamatet/logic/anim"
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

func animateFlipVPiece(ctrl *piece.Ctrl, isAnim bool) {
	if isAnim {
		ctrl.List.Add(anim.NewXRotLin(time.Now(), piece.DurationRotate, -math.Pi))
	}
}

func animateFlipHPiece(ctrl *piece.Ctrl, isAnim bool) {
	if isAnim {
		ctrl.List.Add(anim.NewYRotLin(time.Now(), piece.DurationRotate, -math.Pi))
	}
}

func animateRotatePiece(ctrl *piece.Ctrl, cw, inverted, isAnim bool) {
	if !isAnim {
		return
	}

	t := time.Now()

	if cw {
		ctrl.List.Add(anim.NewZRotQuad(t, piece.DurationRotate, -math.Pi/2))
		if inverted {
			ctrl.List.Add(anim.NewYQuad(t, piece.DurationRotate, -1))
		}
	} else {
		ctrl.List.Add(anim.NewZRotQuad(t, piece.DurationRotate, math.Pi/2))
		if inverted {
			ctrl.List.Add(anim.NewXQuad(t, piece.DurationRotate, 1))
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
