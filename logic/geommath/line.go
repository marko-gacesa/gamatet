// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package geommath

func Line[T Integer](p0, p1 P[T], fn func(p P[T]) bool) {
	dx := Abs[T](p1.X - p0.X)
	dy := Abs[T](p1.Y - p0.Y)

	var sx, sy T

	sx = -1
	if p0.X < p1.X {
		sx = 1
	}

	sy = -1
	if p0.Y < p1.Y {
		sy = 1
	}

	diff := dx - dy

	for {
		if !fn(p0) {
			return
		}

		if p0.X == p1.X && p0.Y == p1.Y {
			break
		}

		e2 := 2 * diff

		if e2 > -dy {
			diff -= dy
			p0.X += sx
		}

		if e2 < dx {
			diff += dx
			p0.Y += sy
		}
	}
}
