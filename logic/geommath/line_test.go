// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package geommath

import (
	"slices"
	"testing"
)

func Test_Line(t *testing.T) {
	tests := []struct {
		name string
		p0   P[int]
		p1   P[int]
		exp  []P[int]
	}{
		{
			name: "dot",
			p0:   P[int]{X: 50, Y: 50},
			p1:   P[int]{X: 50, Y: 50},
			exp:  []P[int]{{X: 50, Y: 50}},
		},
		{
			name: "up-right",
			p0:   P[int]{X: 1, Y: 2},
			p1:   P[int]{X: 7, Y: 5},
			exp: []P[int]{
				{X: 1, Y: 2}, {X: 2, Y: 2},
				{X: 3, Y: 3}, {X: 4, Y: 3},
				{X: 5, Y: 4}, {X: 6, Y: 4},
				{X: 7, Y: 5},
			},
		},
		{
			name: "up-left",
			p0:   P[int]{X: 8, Y: 1},
			p1:   P[int]{X: 6, Y: 8},
			exp: []P[int]{
				{X: 8, Y: 1}, {X: 8, Y: 2},
				{X: 7, Y: 3}, {X: 7, Y: 4}, {X: 7, Y: 5}, {X: 7, Y: 6},
				{X: 6, Y: 7}, {X: 6, Y: 8},
			},
		},
		{
			name: "down-left",
			p0:   P[int]{X: 6, Y: 8},
			p1:   P[int]{X: 4, Y: 1},
			exp: []P[int]{
				{X: 6, Y: 8}, {X: 6, Y: 7},
				{X: 5, Y: 6}, {X: 5, Y: 5}, {X: 5, Y: 4}, {X: 5, Y: 3},
				{X: 4, Y: 2}, {X: 4, Y: 1},
			},
		},
		{
			name: "down-right",
			p0:   P[int]{X: 3, Y: 8},
			p1:   P[int]{X: 4, Y: 1},
			exp: []P[int]{
				{X: 3, Y: 8}, {X: 3, Y: 7}, {X: 3, Y: 6}, {X: 3, Y: 5},
				{X: 4, Y: 4}, {X: 4, Y: 3}, {X: 4, Y: 2}, {X: 4, Y: 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := make([]P[int], 0, len(test.exp))
			Line(test.p0, test.p1,
				func(p P[int]) bool {
					got = append(got, p)
					return true
				},
			)

			if !slices.Equal(got, test.exp) {
				t.Errorf("got %v, want %v", got, test.exp)
			}
		})
	}
}
