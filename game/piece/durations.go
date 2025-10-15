// Copyright (c) 2020-2024 by Marko Gaćeša

package piece

import (
	"math"
	"time"
)

const (
	// DurationInit is initial pause before the game starts.
	DurationInit = 100 * time.Millisecond

	// DurationNewPiece is time to produce the next piece,
	// also time to recheck failed piece placement if piece collision is turned on.
	DurationNewPiece = 15 * time.Millisecond

	// DurationMove is time to move a piece left/right/down.
	DurationMove = 50 * time.Millisecond

	// DurationRotate time to rotate a piece.
	DurationRotate = 100 * time.Millisecond

	// DurationFall is base time for fall animation:
	// The final duration will be multiplied by square root of height.
	DurationFall = 40 * time.Millisecond

	// DurationFullLine is duration of the pause after a full line is cleared.
	DurationFullLine = 75 * time.Millisecond

	// DurationAnimNewPiece is time to animate appearance of a new piece.
	DurationAnimNewPiece = 100 * time.Millisecond

	// DurationAnimBlockChange is time to animate change of a block.
	DurationAnimBlockChange = 750 * time.Millisecond
)

var _durationFall []time.Duration
var _durationDescend []time.Duration
var _durationSlide []time.Duration

func init() {
	_durationFall = make([]time.Duration, 65)
	for height := 0; height < len(_durationFall); height++ {
		_durationFall[height] = time.Duration(float64(DurationFall) * math.Sqrt(float64(height)))
	}

	_durationDescend = make([]time.Duration, MaxLevel+1)
	_durationSlide = make([]time.Duration, MaxLevel+1)

	for level := 0; level <= MaxLevel; level++ {
		_durationDescend[level] = time.Duration(float64(time.Second) * math.Pow(4.0/3.0, float64(-level)))
		_durationSlide[level] = time.Duration(float64(time.Second) * 0.5 * math.Pow(1.125, float64(-level)))
		//fmt.Println(level, _durationDescend[level], _durationSlide[level])
	}
}

func getDescendDuration(level int) time.Duration {
	return _durationDescend[level]
}

func GetFallDuration(height int) time.Duration {
	return _durationFall[height]
}

func GetSlideDuration(level int) time.Duration {
	return _durationSlide[level]
}
