// Copyright (c) 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

import (
	"fmt"
	"slices"
	"testing"

	"github.com/marko-gacesa/gamatet/game/block"
)

func TestField_HasLOS(t *testing.T) {
	const dim = 8
	f := Make(dim, dim, 0)
	f.SetXY(3, 1, AnimNo, 0, block.Wall)
	f.SetXY(3, 2, AnimNo, 0, block.Wall)
	f.SetXY(3, 3, AnimNo, 0, block.Wall)
	f.SetXY(3, 4, AnimNo, 0, block.Wall)

	// 7         2
	// 6       . .
	// 5     . .
	// 4   . . X
	// 3   0   X
	// 2       X
	// 1       X     1
	// 0
	//   0 1 2 3 4 5 6 7

	p0 := block.XY{X: 1, Y: 3}
	p1 := block.XY{X: 6, Y: 1}
	p2 := block.XY{X: 4, Y: 7}

	tests := []struct {
		a, b block.XY
		exp  bool
	}{
		{p0, p1, false},
		{p0, p2, true},
		{p2, p1, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test), func(t *testing.T) {
			want := test.exp
			got := f.HasLOS(test.a, test.b)
			if got != want {
				t.Errorf("HasLOS(%v, %v) = %v, want %v", test.a, test.b, got, want)
			}
		})
	}
}

func TestField_Neighbors8(t *testing.T) {
	const dim = 4
	f := Make(dim, dim, 0)
	f.SetXY(3, 1, AnimNo, 0, block.Wall)
	f.SetXY(1, 3, AnimNo, 0, block.Wall)

	// 3   X 1
	// 2
	// 1       X
	// 0 0     2
	//   0 1 2 3

	tests := []struct {
		name string
		p    block.XY
		exp  Neighbors8
	}{
		{
			name: "x=0,y=0",
			p:    block.XY{X: 0, Y: 0},
			exp: Neighbors8([8]bool{
				false, // (-1, -1)
				false, // (0, -1)
				false, // (1, -1)
				false, // (-1, 0)
				true,  // (1, 0)
				false, // (-1, 1)
				true,  // (0, 1)
				true,  // (1, 1)
			}),
		},
		{
			name: "x=2,y=3",
			p:    block.XY{X: 2, Y: 3},
			exp: Neighbors8([8]bool{
				true,  // (1, 2)
				true,  // (2, 2)
				true,  // (3, 2)
				false, // (1, 3)
				true,  // (3, 3)
				false, // (1, 4)
				false, // (2, 4)
				false, // (3, 4)
			}),
		},
		{
			name: "x=3,y=0",
			p:    block.XY{X: 3, Y: 0},
			exp: Neighbors8([8]bool{
				false, // (2, -1)
				false, // (3, -1)
				false, // (4, -1)
				true,  // (2, 0)
				false, // (4, 0)
				true,  // (2, 1)
				false, // (3, 1)
				false, // (4, 1)
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := test.exp
			got := f.Neighbors8(test.p, validTarget)
			if got != want {
				t.Errorf("Neighbors8(%v) = %v, want %v", test.p, got, want)
			}
		})
	}
}

func TestField_Neighbors4(t *testing.T) {
	const dim = 4
	f := Make(dim, dim, 0)
	f.SetXY(3, 1, AnimNo, 0, block.Wall)
	f.SetXY(1, 3, AnimNo, 0, block.Wall)

	// 3   X 1
	// 2
	// 1       X
	// 0 0     2
	//   0 1 2 3

	tests := []struct {
		name string
		p    block.XY
		exp  Neighbors4
	}{
		{
			name: "x=0,y=0",
			p:    block.XY{X: 0, Y: 0},
			exp: Neighbors4([4]bool{
				false, // (0, -1)
				false, // (-1, 0)
				true,  // (1, 0)
				true,  // (0, 1)
			}),
		},
		{
			name: "x=2,y=3",
			p:    block.XY{X: 2, Y: 3},
			exp: Neighbors4([4]bool{
				true,  // (2, 2)
				false, // (1, 3)
				true,  // (3, 3)
				false, // (2, 4)
			}),
		},
		{
			name: "x=3,y=0",
			p:    block.XY{X: 3, Y: 0},
			exp: Neighbors4([4]bool{
				false, // (3, -1)
				true,  // (2, 0)
				false, // (4, 0)
				false, // (3, 1)
			}),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test.p), func(t *testing.T) {
			want := test.exp
			got := f.Neighbors4(test.p, validTarget)
			if got != want {
				t.Errorf("Neighbors4(%v) = %v, want %v", test.p, got, want)
			}
		})
	}
}

func TestField_Path4(t *testing.T) {
	tests := []struct {
		name  string
		dim   int
		walls []block.XY
		start block.XY
		end   block.XY
		exp   []block.XY
	}{
		{
			name:  "already_there",
			dim:   4,
			start: block.XY{X: 2, Y: 2},
			end:   block.XY{X: 2, Y: 2},
			exp:   []block.XY{{X: 2, Y: 2}},
		},
		{
			name:  "straight_path",
			dim:   4,
			start: block.XY{X: 0, Y: 2},
			end:   block.XY{X: 3, Y: 2},
			exp:   []block.XY{{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}},
		},
		{
			name: "maze_1",
			dim:  6,
			walls: []block.XY{
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 1},
				{X: 3, Y: 1},
				{X: 4, Y: 1},
				{X: 3, Y: 2},
				{X: 1, Y: 3},
				{X: 3, Y: 3},
				{X: 4, Y: 3},
				{X: 1, Y: 4},
			},
			// 5 . . . . . .
			// 4 . X o o o o
			// 3 . X o X X o
			// 2 o o o X . o
			// 1 o X X X X o
			// 0 o X . o o o
			//   0 1 2 3 4 5
			start: block.XY{X: 0, Y: 0},
			end:   block.XY{X: 3, Y: 0},
			exp: []block.XY{
				{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 1, Y: 2},
				{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}, {X: 3, Y: 4},
				{X: 4, Y: 4}, {X: 5, Y: 4}, {X: 5, Y: 3}, {X: 5, Y: 2},
				{X: 5, Y: 1}, {X: 5, Y: 0}, {X: 4, Y: 0}, {X: 3, Y: 0},
			},
		},
		{
			name: "maze_2",
			dim:  6,
			walls: []block.XY{
				{X: 3, Y: 1},
				{X: 4, Y: 1},
				{X: 3, Y: 2},
				{X: 0, Y: 3},
				{X: 1, Y: 3},
				{X: 2, Y: 3},
				{X: 3, Y: 3},
				{X: 5, Y: 3},
				{X: 3, Y: 4},
				{X: 1, Y: 5},
			},
			// 5 o X o o o .
			// 4 o o o X o .
			// 3 X X X X o X
			// 2 . . . X o o
			// 1 . o . X X o
			// 0 . o o o o o
			//   0 1 2 3 4 5
			start: block.XY{X: 1, Y: 1},
			end:   block.XY{X: 0, Y: 5},
			exp: []block.XY{
				{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0},
				{X: 4, Y: 0}, {X: 5, Y: 0}, {X: 5, Y: 1}, {X: 5, Y: 2},
				{X: 4, Y: 2}, {X: 4, Y: 3}, {X: 4, Y: 4}, {X: 4, Y: 5},
				{X: 3, Y: 5}, {X: 2, Y: 5}, {X: 2, Y: 4}, {X: 1, Y: 4},
				{X: 0, Y: 4}, {X: 0, Y: 5},
			},
		},
		{
			name: "no_way",
			dim:  4,
			walls: []block.XY{
				{X: 1, Y: 0},
				{X: 1, Y: 1},
				{X: 2, Y: 2},
				{X: 2, Y: 3},
			},
			// 3 o . X .
			// 2 . . X .
			// 1 . X . .
			// 0 . X . o
			//   0 1 2 3
			start: block.XY{X: 3, Y: 0},
			end:   block.XY{X: 0, Y: 3},
			exp:   nil,
		},
		{
			name:  "big",
			dim:   24,
			walls: nil,
			start: block.XY{X: 10, Y: 10},
			end:   block.XY{X: 23, Y: 2},
			exp: []block.XY{
				{X: 10, Y: 10}, {X: 10, Y: 9}, {X: 10, Y: 8}, {X: 10, Y: 7},
				{X: 10, Y: 6}, {X: 11, Y: 6}, {X: 12, Y: 6}, {X: 13, Y: 6},
				{X: 14, Y: 6}, {X: 15, Y: 6}, {X: 16, Y: 6}, {X: 17, Y: 6},
				{X: 18, Y: 6}, {X: 18, Y: 5}, {X: 19, Y: 5}, {X: 19, Y: 4},
				{X: 19, Y: 3}, {X: 20, Y: 3}, {X: 20, Y: 2}, {X: 21, Y: 2},
				{X: 22, Y: 2}, {X: 23, Y: 2},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(test.dim, test.dim, 0)
			for _, xy := range test.walls {
				f.SetXY(xy.X, xy.Y, AnimNo, 0, block.Wall)
			}

			if want, got := test.exp, f.Path4(test.start, test.end, validTarget); !slices.Equal(got, want) {
				t.Errorf("Path4(%v->%v)\nwant = %v\n got = %v\n", test.start, test.end, want, got)
			}
		})
	}
}

func TestField_Path8(t *testing.T) {
	tests := []struct {
		name  string
		dim   int
		walls []block.XY
		start block.XY
		end   block.XY
		exp   []block.XY
	}{
		{
			name:  "diagonal",
			dim:   4,
			start: block.XY{X: 0, Y: 2},
			end:   block.XY{X: 2, Y: 0},
			exp:   []block.XY{{X: 0, Y: 2}, {X: 1, Y: 1}, {X: 2, Y: 0}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(test.dim, test.dim, 0)
			for _, xy := range test.walls {
				f.SetXY(xy.X, xy.Y, AnimNo, 0, block.Wall)
			}

			if want, got := test.exp, f.Path8(test.start, test.end, validTarget); !slices.Equal(got, want) {
				t.Errorf("Path8(%v->%v)\nwant = %v\n got = %v\n", test.start, test.end, want, got)
			}
		})
	}
}

func TestField_FindNearest8(t *testing.T) {
	tests := []struct {
		name   string
		dx, dy int
		pos    block.XY
		r      int
		exp    []block.XY
	}{
		{
			name: "x=1,y=0",
			dx:   4,
			dy:   6,
			pos:  block.XY{X: 1, Y: 0},
			r:    3,
			// 5
			// 4
			// 3 c c c c
			// 2 b b b b
			// 1 a a a b
			// 0 a o a b
			//   0 1 2 3
			exp: []block.XY{
				{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, // d = 1
				{X: 3, Y: 0}, {X: 3, Y: 1}, {X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, // d = 2
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, // d = 3
			},
		},
		{
			name: "x=1,y=2",
			dx:   4,
			dy:   6,
			pos:  block.XY{X: 1, Y: 2},
			r:    3,
			// 5 c c c c
			// 4 b b b b
			// 3 a a a b
			// 2 a o a b
			// 1 a a a b
			// 0 b b b b
			//   0 1 2 3
			exp: []block.XY{
				// d = 1
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1},
				{X: 0, Y: 2}, {X: 2, Y: 2},
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3},
				// d = 2
				{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0},
				{X: 3, Y: 1}, {X: 3, Y: 2}, {X: 3, Y: 3},
				{X: 0, Y: 4}, {X: 1, Y: 4}, {X: 2, Y: 4}, {X: 3, Y: 4},
				// d = 3
				{X: 0, Y: 5}, {X: 1, Y: 5}, {X: 2, Y: 5}, {X: 3, Y: 5},
			},
		},
		{
			name: "x=2,y=5",
			dx:   4,
			dy:   6,
			pos:  block.XY{X: 2, Y: 5},
			r:    3,
			// 5 b a o a
			// 4 b a a a
			// 3 b b b b
			// 2 c c c c
			// 1
			// 0
			//   0 1 2 3
			exp: []block.XY{
				{X: 1, Y: 4}, {X: 2, Y: 4}, {X: 3, Y: 4}, {X: 1, Y: 5}, {X: 3, Y: 5}, // d = 1
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, // d = 2
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, // d = 3
			},
		},
		{
			name: "x=0,y=4",
			dx:   4,
			dy:   6,
			pos:  block.XY{X: 0, Y: 4},
			r:    3,
			// 5 a a b c
			// 4 o a b c
			// 3 a a b c
			// 2 b b b c
			// 1 c c c c
			// 0
			//   0 1 2 3
			exp: []block.XY{
				{X: 0, Y: 3}, {X: 1, Y: 3}, {X: 1, Y: 4}, {X: 0, Y: 5}, {X: 1, Y: 5}, // d = 1
				{X: 0, Y: 2}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}, {X: 2, Y: 5}, // d = 2
				{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 3, Y: 2}, {X: 3, Y: 3}, {X: 3, Y: 4}, {X: 3, Y: 5}, // d = 3
			},
		},
		{
			name: "x=3,y=0",
			dx:   4,
			dy:   6,
			pos:  block.XY{X: 3, Y: 0},
			r:    3,
			// 5
			// 4
			// 3 c c c c
			// 2 c b b b
			// 1 c b a a
			// 0 c b a o
			//   0 1 2 3
			exp: []block.XY{
				{X: 2, Y: 0}, {X: 2, Y: 1}, {X: 3, Y: 1}, // d = 1
				{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 3, Y: 2}, // d = 2
				{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 1, Y: 3}, {X: 2, Y: 3}, {X: 3, Y: 3}, // d = 3
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(test.dx, test.dy, 0)

			var got []block.XY
			f.FindNearest8(test.pos, test.r, func(xyb block.XYB, d int) bool {
				got = append(got, xyb.XY)
				return false
			})

			want := test.exp

			if !slices.Equal(got, want) {
				t.Errorf("want = %v\n got = %v\n", want, got)
			}
		})
	}
}

func TestField_FindNearest8LOS(t *testing.T) {
	const dim = 4
	tests := []struct {
		name   string
		setup  []block.XYB
		pos    block.XY
		expXYB block.XYB
		expOK  bool
	}{
		{
			name: "can-see-1",
			setup: []block.XYB{
				{XY: block.XY{X: 0, Y: 3}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 1, Y: 2}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 3, Y: 0}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 3, Y: 2}, Block: block.Block{Type: block.TypeRock, Color: 1}},
			},
			pos: block.XY{X: 1, Y: 0},
			// 3 W
			// 2   W   g
			// 1
			// 0 p     W
			//   0 1 2 3
			expXYB: block.XYB{
				XY:    block.XY{X: 3, Y: 2},
				Block: block.Block{Type: block.TypeRock, Color: 1},
			},
			expOK: true,
		},
		{
			name: "can-see-2",
			setup: []block.XYB{
				{XY: block.XY{X: 0, Y: 3}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 1, Y: 1}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 2, Y: 2}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 2, Y: 0}, Block: block.Block{Type: block.TypeRock, Color: 2}},
			},
			pos: block.XY{X: 1, Y: 3},
			// 3 W p
			// 2     W
			// 1   W
			// 0     g
			//   0 1 2 3
			expXYB: block.XYB{
				XY:    block.XY{X: 2, Y: 0},
				Block: block.Block{Type: block.TypeRock, Color: 2},
			},
			expOK: true,
		},
		{
			name: "can-see-3",
			setup: []block.XYB{
				{XY: block.XY{X: 1, Y: 3}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 1, Y: 2}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 1, Y: 0}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 0, Y: 1}, Block: block.Block{Type: block.TypeRock, Color: 3}},
			},
			pos: block.XY{X: 2, Y: 1},
			// 3   W
			// 2   W
			// 1 g   p
			// 0   W
			//   0 1 2 3
			expXYB: block.XYB{
				XY:    block.XY{X: 0, Y: 1},
				Block: block.Block{Type: block.TypeRock, Color: 3},
			},
			expOK: true,
		},
		{
			name: "can-see-4",
			setup: []block.XYB{
				{XY: block.XY{X: 2, Y: 1}, Block: block.Block{Type: block.TypeRock, Color: 4}},
			},
			pos: block.XY{X: 1, Y: 1},
			// 3
			// 2
			// 1   p g
			// 0
			//   0 1 2 3
			expXYB: block.XYB{
				XY:    block.XY{X: 2, Y: 1},
				Block: block.Block{Type: block.TypeRock, Color: 4},
			},
			expOK: true,
		},
		{
			name: "can-not-see-1",
			setup: []block.XYB{
				{XY: block.XY{X: 1, Y: 1}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 0, Y: 1}, Block: block.Block{Type: block.TypeRock}},
			},
			pos: block.XY{X: 2, Y: 1},
			// 3
			// 2
			// 1 g W p
			// 0
			//   0 1 2 3
			expXYB: block.XYB{
				XY:    block.XY{X: 2, Y: 1},
				Block: block.Block{Type: block.TypeEmpty},
			},
			expOK: false,
		},
		{
			name: "can-not-see-2",
			setup: []block.XYB{
				{XY: block.XY{X: 2, Y: 3}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 2, Y: 2}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 2, Y: 1}, Block: block.Block{Type: block.TypeWall}},
				{XY: block.XY{X: 1, Y: 3}, Block: block.Block{Type: block.TypeRock}},
			},
			pos: block.XY{X: 3, Y: 0},
			// 3   g W
			// 2     W
			// 1     W
			// 0       p
			//   0 1 2 3
			expXYB: block.XYB{
				XY:    block.XY{X: 3, Y: 0},
				Block: block.Block{Type: block.TypeEmpty},
			},
			expOK: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(dim, dim, 0)

			for _, xyb := range test.setup {
				f.setXY(xyb.X, xyb.Y, xyb.Block)
			}

			gotXYB, gotOK := f.FindNearest8(test.pos, dim, func(xyb block.XYB, d int) bool {
				return xyb.Block.Type == block.TypeRock && f.HasLOS(test.pos, xyb.XY)
			})

			wantXYB, wantOK := test.expXYB, test.expOK

			if wantXYB != gotXYB || wantOK != gotOK {
				t.Errorf("want = %v %t\n got = %v %t\n", wantXYB, wantOK, gotXYB, gotOK)
			}
		})
	}
}

func TestField_FindNearest4(t *testing.T) {
	tests := []struct {
		name   string
		dx, dy int
		pos    block.XY
		r      int
		exp    []block.XY
	}{
		{
			name: "x=1,y=1",
			dx:   4,
			dy:   6,
			pos:  block.XY{X: 1, Y: 1},
			r:    3,
			// 5
			// 4   c
			// 3 c b c
			// 2 b a b c
			// 1 a o a b
			// 0 b a b c
			//   0 1 2 3
			exp: []block.XY{
				{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 2, Y: 1}, {X: 1, Y: 2}, // d = 1
				{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 1}, {X: 0, Y: 2}, {X: 2, Y: 2}, {X: 1, Y: 3}, // d = 2
				{X: 3, Y: 0}, {X: 3, Y: 2}, {X: 0, Y: 3}, {X: 2, Y: 3}, {X: 1, Y: 4}, // d = 3
			},
		},
		{
			name: "x=2,y=5",
			dx:   4,
			dy:   6,
			pos:  block.XY{X: 2, Y: 5},
			r:    3,
			// 5 b a o a
			// 4 c b a b
			// 3   c b c
			// 2     c
			// 1
			// 0
			//   0 1 2 3
			exp: []block.XY{
				{X: 2, Y: 4}, {X: 1, Y: 5}, {X: 3, Y: 5}, // d = 1
				{X: 2, Y: 3}, {X: 1, Y: 4}, {X: 3, Y: 4}, {X: 0, Y: 5}, // d = 2
				{X: 2, Y: 2}, {X: 1, Y: 3}, {X: 3, Y: 3}, {X: 0, Y: 4}, // d = 3
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := Make(test.dx, test.dy, 0)

			var got []block.XY
			f.FindNearest4(test.pos, test.r, func(xyb block.XYB, d int) bool {
				got = append(got, xyb.XY)
				return false
			})

			want := test.exp

			if !slices.Equal(got, want) {
				t.Errorf("want = %v\n got = %v\n", want, got)
			}
		})
	}
}

func TestNeighbors8_ForEach(t *testing.T) {
	const dim = 4
	f := Make(dim, dim, 0)

	tests := []struct {
		name string
		p    block.XY
		exp  []block.XY
	}{
		{
			name: "bottom-left",
			p:    block.XY{X: 0, Y: 0},
			exp: []block.XY{
				{X: 1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1},
			},
		},
		{
			name: "bottom-right",
			p:    block.XY{X: dim - 1, Y: 0},
			exp: []block.XY{
				{X: dim - 2, Y: 0}, {X: dim - 2, Y: 1}, {X: dim - 1, Y: 1},
			},
		},
		{
			name: "top-right",
			p:    block.XY{X: 0, Y: dim - 1},
			exp: []block.XY{
				{X: 0, Y: dim - 2}, {X: 1, Y: dim - 2}, {X: 1, Y: dim - 1},
			},
		},
		{
			name: "top-right",
			p:    block.XY{X: dim - 1, Y: dim - 1},
			exp: []block.XY{
				{X: dim - 2, Y: dim - 2}, {X: dim - 1, Y: dim - 2}, {X: dim - 2, Y: dim - 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := test.exp
			n := f.Neighbors8(test.p, validTarget)
			var got []block.XY
			n.ForEach(f, test.p, func(xyb block.XYB) {
				got = append(got, xyb.XY)
			})

			if !slices.Equal(want, got) {
				t.Errorf("")
			}
		})
	}
}

func TestNeighbors4_ForEach(t *testing.T) {
	const dim = 4
	f := Make(dim, dim, 0)

	tests := []struct {
		name string
		p    block.XY
		exp  []block.XY
	}{
		{
			name: "bottom-left",
			p:    block.XY{X: 0, Y: 0},
			exp: []block.XY{
				{X: 1, Y: 0}, {X: 0, Y: 1},
			},
		},
		{
			name: "bottom-right",
			p:    block.XY{X: dim - 1, Y: 0},
			exp: []block.XY{
				{X: dim - 2, Y: 0}, {X: dim - 1, Y: 1},
			},
		},
		{
			name: "top-left",
			p:    block.XY{X: 0, Y: dim - 1},
			exp: []block.XY{
				{X: 0, Y: dim - 2}, {X: 1, Y: dim - 1},
			},
		},
		{
			name: "top-right",
			p:    block.XY{X: dim - 1, Y: dim - 1},
			exp: []block.XY{
				{X: dim - 1, Y: dim - 2}, {X: dim - 2, Y: dim - 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			want := test.exp
			n := f.Neighbors4(test.p, validTarget)
			var got []block.XY
			n.ForEach(f, test.p, func(xyb block.XYB) {
				got = append(got, xyb.XY)
			})

			if !slices.Equal(want, got) {
				t.Errorf("")
			}
		})
	}
}

func validTarget(t block.Type) bool { return t == block.TypeEmpty || t == block.TypeRock }
