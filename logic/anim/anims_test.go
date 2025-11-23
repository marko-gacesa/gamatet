// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package anim

import (
	"testing"
	"time"
)

func TestAnim(t *testing.T) {
	startTime := time.Date(2020, 11, 01, 12, 30, 45, 0, time.UTC)

	type expect struct {
		delay      time.Duration
		dx, dy, dz float32
		isDone     bool
	}

	tests := []struct {
		name     string
		anim     Anim
		expected []expect
	}{
		{
			name: "linear translation",
			anim: NewTransLin(startTime, time.Minute, 100, 50, 40),
			expected: []expect{
				{delay: -1 * time.Second, dx: -100, dy: -50, dz: -40},
				{dx: -100, dy: -50, dz: -40},
				{delay: 15 * time.Second, dx: -75, dy: -37.5, dz: -30},
				{delay: 30 * time.Second, dx: -50, dy: -25, dz: -20},
				{delay: 45 * time.Second, dx: -25, dy: -12.5, dz: -10},
				{delay: 60 * time.Second, isDone: true},
			},
		},
		{
			name: "quadratic translation",
			anim: NewTransQuad(startTime, time.Minute, 100, 50, 40),
			expected: []expect{
				{delay: -1 * time.Second, dx: -100, dy: -50, dz: -40},
				{dx: -100, dy: -50, dz: -40},
				{delay: 15 * time.Second, dx: -56.25, dy: -28.125, dz: -22.5},
				{delay: 30 * time.Second, dx: -25, dy: -12.5, dz: -10},
				{delay: 45 * time.Second, dx: -6.25, dy: -3.125, dz: -2.5},
				{delay: 60 * time.Second, isDone: true},
			},
		},
	}

	for _, test := range tests {
		for _, e := range test.expected {
			now := startTime
			now = now.Add(e.delay)
			isDone := test.anim.Update(now)
			dx, dy, dz := test.anim.Translate()

			if isDone != e.isDone {
				t.Errorf("test %q failed; after %q expected done=%t but got %t", test.name, e.delay, e.isDone, isDone)
				continue
			}
			if dx != e.dx {
				t.Errorf("test %q failed; after %q expected dx=%f but got %f", test.name, e.delay, e.dx, dx)
				continue
			}
			if dy != e.dy {
				t.Errorf("test %q failed; after %q expected dy=%f but got %f", test.name, e.delay, e.dy, dy)
				continue
			}
			if dz != e.dz {
				t.Errorf("test %q failed; after %q expected dz=%f but got %f", test.name, e.delay, e.dz, dz)
				continue
			}
		}
	}
}
