// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package menu

type Iter struct {
	Title       string
	Description string
	Current     int
	Items       []string
	IsDisabled  []bool
}
