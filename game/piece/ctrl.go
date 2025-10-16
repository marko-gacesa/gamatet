// Copyright (c) 2020-2025 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
	"gamatet/logic/anim"
	"time"
)

const MaxLevel = 15

const MaxWallKick = 2

const NextBlockCount = 3

type State byte

const (
	StatePause State = iota
	StateInit
	StateNew
	StateDescend
	StateFall
	StateSlide
	StateWon
	StateLost
	StateStop
)

var StateName = map[State]string{
	StatePause:   "Pause",
	StateInit:    "Init",
	StateNew:     "New",
	StateDescend: "Descend",
	StateFall:    "Fall",
	StateSlide:   "Slide",
	StateWon:     "Won",
	StateLost:    "Lost",
	StateStop:    "Stop",
}

// IsPausable returns if the game can paused in the current state.
func (s State) IsPausable() bool {
	return s != StateLost && s != StateWon && s != StateStop
}

// IsAbortable returns if the game can be aborted in the current state.
func (s State) IsAbortable() bool {
	return s == StatePause || s == StateWon || s == StateLost || s == StateStop
}

// IsTerminal returns if the current state is the final state.
func (s State) IsTerminal() bool {
	return s == StateWon || s == StateLost || s == StateStop
}

type NextPieceInfo struct {
	Type   Type
	Blocks []block.XYB
}

type Ctrl struct {
	// Idx is an index of Ctrl in the field
	Idx int

	// Name of the player
	Name string

	// X and Y are current position of the piece on the board
	X, Y int

	// Feed is the next piece feed
	Feed Feed

	// Piece is the piece itself
	Piece Piece

	// Blocks is copy of individual piece blocks
	Blocks []block.XYB

	NextPieces [NextBlockCount]NextPieceInfo

	// IsShadowShown is true if the shadow should be rendered, or false is there is no shadow
	IsShadowShown bool
	// Shadow is used to render the piece shadow
	Shadow Shadow

	// IsColumnLimited turns on or off column range limitation for the piece
	IsColumnLimited bool
	// ColumnLimit holds column range to which the piece is limited
	ColumnLimit ColumnLimit

	// Score is player's score.
	Score    int
	ScoreStr string

	// PieceCount is the index number of the piece since the game start
	PieceCount    int
	PieceCountStr string

	// Level is the speed at the player plays
	Level    int
	LevelStr string

	// Config contains piece related configuration
	Config

	// List contains the piece animations
	List anim.List

	// State is used for the state machine
	State State
	// Timer is used for different purpose depending on State
	Timer *time.Timer

	// PausedState holds state when State=StatePaused
	PausedState State

	InfoPosition DisplayPosition
}

func NewCtrl(idx int) *Ctrl {
	const zeroStr = "0"

	c := &Ctrl{}
	c.Idx = idx
	c.Name = ""

	c.Score = 0
	c.ScoreStr = zeroStr
	c.PieceCount = 0
	c.PieceCountStr = zeroStr
	c.Level = 0
	c.LevelStr = zeroStr

	c.State = StateInit
	c.InfoPosition = DisplayPosition(idx + 1)
	return c
}

func (c *Ctrl) SetXYP(x, y int, p Piece) {
	c.X = x
	c.Y = y
	c.Piece = p
	c.List.Clear()
}

func (c *Ctrl) RestartTimer(param int) {
	var dur time.Duration

	switch c.State {
	case StateInit:
		dur = DurationInit
	case StateNew:
		dur = DurationNewPiece
	case StateDescend:
		dur = getDescendDuration(c.Level)
	case StateFall:
		if param > 0 {
			dur = GetFallDuration(param)
		} else {
			dur = DurationFall
		}
	case StateSlide:
		dur = GetSlideDuration(c.Level)
	case StateLost, StateWon:
		dur = DurationNewPiece
	case StatePause, StateStop:
		// no timer for these states
	default:
		panic("invalid state")
	}

	c.StopTimer()

	if dur == 0 {
		return
	}

	if c.Timer == nil {
		c.Timer = time.NewTimer(dur)
	} else {
		c.Timer.Reset(dur)
	}
}

func (c *Ctrl) StopTimer() {
	if c.Timer != nil && !c.Timer.Stop() {
		select {
		default:
		case <-c.Timer.C:
		}
	}
}

type Config struct {
	// RotationDirectionCW determines Piece rotation direction.
	// Value true means the rotation is clockwise and performed with RotateCW().
	// Value false means the rotation is counterclockwise and performed with RotateCCW().
	RotationDirectionCW bool

	// SlideDisabled enables slide mode after a piece drop.
	// It enables piece move, left or right, while the piece is falling
	SlideDisabled bool

	// WallKick enables wall kick feature.
	// On piece rotate command when a piece is next to a wall and rotation is not possible because of the wall,
	// the game will try to move the piece max WallKick places to the left or right to try to place the rotated piece.
	WallKick byte
}

type Shadow struct {
	// Blocks of the shadow
	Blocks []block.XYB

	// ColL is the left most column of the shadow
	ColL int

	// ColR is the right most column of the shadow
	ColR int
}

type ColumnLimit struct {
	// Min is the left most column allowed for the piece
	Min int

	// Max is the right most + 1 column allowed for the piece
	Max int
}
