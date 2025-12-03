// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

// LevelUpBlocks returns number of destroyed blocks needed to reach the desired level of speed.
//
// Based on the code used in GMT1 30 years earlier:
// To reach the next level from current level it takes this many lines
// progression := []int{11, 13, 15, 17, 19, 21, 23, 25, 27, 29}
// neededLines := progression[current_level+1]
func LevelUpBlocks(l, w int) int {
	needed := (10*l + l*l) * w
	return needed
}
