// Copyright (c) 2020-2024 by Marko Gaćeša

package piece

import (
	"gamatet/game/block"
	"gamatet/logic/anim"
	"time"
)

const MaxLevel = 15

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

type Ctrl struct {
	// Idx is an index of Ctrl in the field
	Idx int

	// Name of the player
	Name string

	// Score is player's score.
	Score int

	// X and Y are current position of the piece on the board
	X, Y int

	// Feed is the next piece feed
	Feed Feed

	// Piece is the piece itself
	Piece Piece

	// Blocks is copy of individual piece blocks
	Blocks []block.XYB

	NextBlocks [NextBlockCount][]block.XYB

	// IsShadowShown is true if the shadow should be rendered, or false is there is no shadow
	IsShadowShown bool
	// Shadow is used to render the piece shadow
	Shadow Shadow

	// IsColumnLimited turns on or off column range limitation for the piece
	IsColumnLimited bool
	// ColumnLimit holds column range to which the piece is limited
	ColumnLimit ColumnLimit

	// PieceCount is the index number of the piece since the game start
	PieceCount int

	// Level is the speed at the player plays
	Level int

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

	// SlideEnabled enables slide mode after a piece drop.
	// It enables piece move, left or right, while the piece is falling
	SlideEnabled bool

	// MaxWallKick enables wall kick feature.
	// On piece rotate command when a piece is next to a wall and rotation is not possible because of the wall,
	// the game will try to move the piece max WallKick places to the left or right to try to place the rotated piece.
	MaxWallKick int
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
