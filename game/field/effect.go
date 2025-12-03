// Copyright (c) 2020 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package field

type Effect byte

const (
	EffectNone  Effect = 0
	EffectLid   Effect = 10
	EffectBigO  Effect = 11
	EffectRaise Effect = 12
	EffectPatch Effect = 20
)
