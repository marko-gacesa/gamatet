// Copyright (c) 2020 by Marko Gaćeša

package action

type Action byte

const (
	NoOp      Action = 0
	Abort     Action = 10
	Pause     Action = 11
	SpeedUp   Action = 20
	MoveLeft  Action = 30
	MoveRight Action = 31
	MoveUp    Action = 32
	MoveDown  Action = 33
	Drop      Action = 40
	Rotate    Action = 50
	RotateCW  Action = 51
	RotateCCW Action = 52
)
