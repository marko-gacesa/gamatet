// Copyright (c) 2020, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/marko-gacesa/gamatet/game/block"
)

func TestGenericFeed(t *testing.T) {
	tests := []struct {
		bagSize int
		seed    int
		bags    []int
	}{
		{bagSize: 1, seed: 45, bags: []int{0, 34, 54, 134, 1223}},
		{bagSize: 2, seed: 23, bags: []int{1, 4, 14, 7223}},
		{bagSize: 3, seed: 79, bags: []int{5, 637, 17223}},
		{bagSize: 4, seed: 34, bags: []int{3, 73, 953, 64531}},
	}

	printBag := func(pieces []polyominoRot) string {
		m := make([][]string, len(pieces))
		for i := range pieces {
			b := strings.Split(pieces[i].String(), "\n")
			m[i] = b[0 : len(b)-1]
			for len(m[i]) <= 4 {
				m[i] = append(m[i], strings.Repeat("_", len(m[i][0])))
			}
		}
		sb := strings.Builder{}
		for col := range 4 {
			for i := range pieces {
				sb.WriteString(m[i][col])
				sb.WriteString(" | ")
			}
			sb.WriteByte('\n')
		}
		return sb.String()
	}

	for testIdx, test := range tests {
		t.Run(fmt.Sprintf("seed:%d-bagSize:%d", test.seed, test.bagSize), func(t *testing.T) {
			shapes := shapesRotTetrominoes
			shapeCount := len(shapes)

			f := NewGenericFeed(test.bagSize, test.seed, shapeCount, func(idx, shapeIdx uint, playerIndex byte) Piece {
				return &polyominoRot{
					shapeSquare: shapes[shapeIdx],
					block:       block.Rock,
				}
			})

			pieceBagCount := test.bagSize * shapeCount

			// init

			type bagCheck struct {
				pieceCount map[polyominoRot]int
				pieces     []polyominoRot
			}
			bagTest := make(map[int]bagCheck, len(test.bags))
			for _, bag := range test.bags {
				bagTest[bag] = bagCheck{
					pieceCount: make(map[polyominoRot]int, shapeCount),
					pieces:     make([]polyominoRot, pieceBagCount),
				}
			}

			// the first test: uniqueness of pieces in a single bag (there must be bagSize of each piece in each bag)

			for _, bag := range test.bags {
				idx := bag * pieceBagCount
				for bagIdx := range pieceBagCount {
					p := f.Get(uint(idx+bagIdx), 0).(*polyominoRot)
					bagTest[bag].pieceCount[*p]++
					bagTest[bag].pieces[bagIdx] = *p
				}
			}

			for bag, bagCheckData := range bagTest {
				for p, count := range bagCheckData.pieceCount {
					if count != test.bagSize {
						t.Errorf("uniqueness test failed: test#=%d in bag=%d expected count=%d, got=%d for piece:\n%s\n",
							testIdx, bag, test.bagSize, count, p.String())
					}
				}
			}

			// the second test: order of pieces in each bag must be the same

			for _, bag := range test.bags {
				if len(bagTest[bag].pieceCount) != shapeCount {
					t.Errorf("piece type count failed: test#=%d in bag=%d expected piece count=%d got piece count=%d",
						testIdx, bag, shapeCount, len(bagTest[bag].pieceCount))
				}

				idx := bag * pieceBagCount
				for bagIdx := range pieceBagCount {
					p := f.Get(uint(idx+bagIdx), 0).(*polyominoRot)
					if bagTest[bag].pieces[bagIdx] != *p {
						t.Errorf("piece order test failed: test#=%d in bag=%d bagIdx=%d expected piece:\n%s\ngot piece:\n%s\n",
							testIdx, bag, bagIdx, bagTest[bag].pieces[bagIdx].String(), p.String())
					}
				}
			}

			// the third test: make sure order of pieces is different in each bag

			for i := 0; i < len(test.bags)-1; i++ {
				for j := i + 1; j < len(test.bags); j++ {
					bag1 := bagTest[test.bags[i]].pieces
					bag2 := bagTest[test.bags[j]].pieces
					if reflect.DeepEqual(bag1, bag2) {
						t.Errorf("bag piece order uniqueness test failed: test#=%d in bag1=%d bag2=%d\nbag1 pieces:\n%sbag2 pieces:\n%s",
							testIdx, i, j, printBag(bag1), printBag(bag2))
					}
				}
			}
		})
	}
}
