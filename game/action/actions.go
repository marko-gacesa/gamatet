// Copyright (c) 2020, 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package action

import "fmt"

type Action byte

const (
	NoOp Action = 0

	Abort Action = 10
	Pause Action = 11

	SpeedUp   Action = 20
	SpeedDown Action = 21
	MoveLeft  Action = 30
	MoveRight Action = 31
	MoveDown  Action = 33
	Drop      Action = 40
	Activate  Action = 50
)

func (a Action) String() string {
	switch a {
	case NoOp:
		return "NoOp"
	case Abort:
		return "Abort"
	case Pause:
		return "Pause"
	case SpeedUp:
		return "SpeedUp"
	case SpeedDown:
		return "SpeedDown"
	case MoveLeft:
		return "MoveLeft"
	case MoveRight:
		return "MoveRight"
	case MoveDown:
		return "MoveDown"
	case Drop:
		return "Drop"
	case Activate:
		return "Activate"
	default:
		return fmt.Sprintf("Action(%X)", byte(a))
	}
}
