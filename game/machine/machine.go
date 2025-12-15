// Copyright (c) 2020-2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package machine

import (
	"reflect"
	"time"

	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/game/block"
	"github.com/marko-gacesa/gamatet/game/event"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/op"
	"github.com/marko-gacesa/gamatet/game/piece"
)

// Description for player action state machine:
// * Each state transition resets the timer (even if the transition is to the same state)
//
// Init: Initial state; Player has no piece
// * on timeout: ChangeState(New)
//
// New: Player has no piece; A new piece should be generated
// * on timeout: Take next piece (the piece at index PieceCount)...
// * ... if it can't be introduced the because of some other player's piece: ChangeState(New)
// * ... if it can't be introduced the because of a block on the field, the game is over: ChangeState(End)
// * ... if it can: ChangeState(Descend), PieceCount++, SetPiece (and it's location)
//
// Descend: Player has a piece, and it descends slowly
// * on timeout: Move piece down by one
// * ... successful: ChangeState(Descend)
// * ... failed: ChangeState(Slide)
// * on Drop: Move piece down by fall height
// * ... successful (height>0): ChangeState(Fall)
// * ... failed: ChangeState(Slide)
// * on MoveDown:
// * ... successful: ChangeState(Descend)
// * ... failed: ChangeState(Slide)
// * on MoveLeft, MoveRight, Rotate:
// * ... success: just perform the operation, the timer untouched, the state unchanged
//
// Fall: Player has a piece and has dropped it:
// * on timeout:
// * ... ChangeState(Slide)
// * on MoveLeft, MoveRight:
// * ... success+slide enabled: just perform the operation, the timer untouched, the state unchanged
//
// Slide: Player has a piece, and it has just hit the bottom. Player has a little time to react and adjust its position
// * on timeout: Move piece down by one
// * ... successful: ChangeState(Descend)
// * ... failed: MeldPiece, ClearPiece, ChangeState(New)
// * on Drop: Move piece down by fall height
// * ... successful: ChangeState(Fall)
// * ... failed: do nothing
// * on MoveDown:
// * ... successful: ChangeState(Descend)
// * ... failed: do nothing
// * on MoveLeft, MoveRight, Rotate:
// * ... success: just perform the operation, the timer untouched, the state unchanged
//
// Lost, Won:
// * on timeout: ChangeState(Stop)
//
// Stop:
// * note: this is the final state
//

func HandleActionInput(f *field.Field, ctrl *piece.Ctrl, p event.Pusher, a action.Action) {
	switch a {
	case action.MoveLeft, action.MoveRight:
		if ctrl.State != piece.StateDescend &&
			(ctrl.State != piece.StateSlide && ctrl.State != piece.StateFall || ctrl.SlideDisabled) {
			return
		}

		var dx int
		if a == action.MoveLeft {
			dx = -1
		} else {
			dx = 1
		}

		if success := f.CanMovePiece(dx, 0, ctrl.Idx, !f.PieceCollision); success {
			p.Push(op.NewPieceMove(ctrl.Idx, dx, 0))
		}

	case action.Activate:
		if ctrl.State != piece.StateDescend && ctrl.State != piece.StateSlide {
			return
		}

		switch ctrl.Piece.Type() {
		case piece.TypeRotation:
			if success, _, dx, dy, _ := f.CanRotatePiece(ctrl.Idx, !f.PieceCollision); success {
				if dx != 0 || dy != 0 {
					p.Push(op.NewPieceMove(ctrl.Idx, dx, dy))
				}
				p.Push(op.NewPieceRotate(ctrl.Idx, ctrl.RotationDirectionCW))
			}
		case piece.TypeFlipV:
			if success, _ := f.CanFlipVPiece(ctrl.Idx, !f.PieceCollision); success {
				p.Push(op.NewPieceFlip(ctrl.Idx))
			}
		case piece.TypeFlipH:
			if success, _ := f.CanFlipHPiece(ctrl.Idx, !f.PieceCollision); success {
				p.Push(op.NewPieceFlip(ctrl.Idx))
			}
		case piece.TypeShooter:
			_shoot(f, ctrl, p)
		}

	case action.MoveDown:
		if ctrl.State != piece.StateDescend && ctrl.State != piece.StateSlide {
			return
		}

		success := f.CanMovePiece(0, -1, ctrl.Idx, !f.PieceCollision)
		if !success {
			if ctrl.State != piece.StateSlide {
				_changeState(ctrl, p, piece.StateSlide)
			} else {
				ctrl.Timer.Reset(time.Nanosecond)
			}
			return
		}

		p.Push(op.NewPieceMove(ctrl.Idx, 0, -1))

		_changeState(ctrl, p, piece.StateDescend)

	case action.Drop:
		if ctrl.State != piece.StateDescend && ctrl.State != piece.StateSlide {
			return
		}

		switch ctrl.Piece.Type() {
		case piece.TypeDumb, piece.TypeFlipV, piece.TypeFlipH, piece.TypeRotation:
			if t := ctrl.Blocks[0].Type; t.NoSlide() {
				_meldPiece(f, ctrl, p)
				_clearPiece(ctrl, p)
				return
			}
		case piece.TypeShooter:
			_shoot(f, ctrl, p)
			return
		}

		height := f.GetDropHeight(ctrl.Idx, !f.PieceCollision)
		if height == 0 {
			if ctrl.State != piece.StateSlide {
				_changeState(ctrl, p, piece.StateSlide)
			}
			return
		}

		p.Push(op.NewPieceFall(ctrl.Idx, height))
		p.Push(op.NewPieceScore(ctrl.Idx, 2*height*int(ctrl.Level)))

		_changeStateWithParam(ctrl, p, piece.StateFall, height)

	case action.SpeedUp:
		p.Push(op.NewPieceLevelBoost(ctrl.Idx, true))
	case action.SpeedDown:
		p.Push(op.NewPieceLevelBoost(ctrl.Idx, false))
	}
}

func HandleTimeout(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) bool {
	switch ctrl.State {
	case piece.StateInit:
		_changeState(ctrl, p, piece.StateNew)

	case piece.StateNew:
		if ctrl.Piece != nil {
			// should not happen: player should not have a piece in this state
			panic("player already has a piece in state New")
		}

		pieceCount := ctrl.PieceCount

		newPiece := ctrl.Feed.Get(pieceCount, ctrl.PlayerIndex)
		success, x, y := f.GetPieceStartPosition(ctrl.Idx, ctrl, newPiece, !f.PieceCollision)
		if !success && f.PieceCollision {
			success, x, y = f.GetPieceStartPosition(ctrl.Idx, ctrl, newPiece, true)
			if success {
				_changeState(ctrl, p, piece.StateNew) // wait awhile, there is a piece in the way
				break
			}
		}
		if !success {
			// can't position piece: end game
			return true
		}

		p.Push(op.NewPieceSet(ctrl.Idx, op.TypeSet, x, y, newPiece, pieceCount+1))

		_changeState(ctrl, p, piece.StateDescend)

	case piece.StateDescend, piece.StateSlide:
		if f.CanMovePiece(0, -1, ctrl.Idx, !f.PieceCollision) {
			p.Push(op.NewPieceFall(ctrl.Idx, 1))
			p.Push(op.NewPieceScore(ctrl.Idx, int(ctrl.ApparentLevel())))
			_changeState(ctrl, p, piece.StateDescend)
			break
		}

		if ctrl.State == piece.StateDescend {
			_changeState(ctrl, p, piece.StateSlide)
			break
		}

		if ctrl.Piece.Type() == piece.TypeShooter {
			_clearPiece(ctrl, p)
			break
		}

		_meldPiece(f, ctrl, p)
		_clearPiece(ctrl, p)

	case piece.StateFall:
		_changeState(ctrl, p, piece.StateSlide)

	case piece.StatePause:

	case piece.StateStop:
		if ctrl.Piece != nil {
			p.Push(op.NewPieceSet(ctrl.Idx, op.TypeClear, ctrl.X, ctrl.Y, ctrl.Piece, ctrl.PieceCount))
		}

	default:
		// should not happen, unknown state
		panic("timer timeout for unknown state: " + ctrl.State.String())
	}

	return false
}

func _changeState(ctrl *piece.Ctrl, p event.Pusher, newState piece.State) {
	p.Push(op.NewPieceState(ctrl.Idx, ctrl.State, newState, 0, 0))
}

func _changeStateWithParam(ctrl *piece.Ctrl, p event.Pusher, newState piece.State, param int) {
	p.Push(op.NewPieceState(ctrl.Idx, ctrl.State, newState, 0, param))
}

func _clearPiece(ctrl *piece.Ctrl, p event.Pusher) {
	p.Push(op.NewPieceSet(ctrl.Idx, op.TypeClear, ctrl.X, ctrl.Y, ctrl.Piece, ctrl.PieceCount))
	_changeState(ctrl, p, piece.StateNew)
}

func _meldPiece(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	switch ctrl.Piece.Type() {
	case piece.TypeDumb, piece.TypeFlipV, piece.TypeFlipH, piece.TypeRotation:
		switch ctrl.Blocks[0].Type {
		case block.TypeLava, block.TypeAcid, block.TypeBomb, block.TypeCurl, block.TypeWave:
			_meldLiquidPiece(f, ctrl, p)
		default:
			_meldSolidPiece(f, ctrl, p)
		}
	case piece.TypeShooter:
	}
}

func _meldSolidPiece(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	// blocks are returned from the lowest row, up to the topmost row
	blocks := f.GetPieceBlockLocations(ctrl.X, ctrl.Y, ctrl.Piece)

	for i := len(blocks) - 1; i >= 0; i-- {
		_setBlock(f, p, blocks[i].X, blocks[i].Y, blocks[i].Block)
	}
}

func _meldLiquidPiece(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	// blocks are returned from the lowest row, up to the topmost row
	blocks := f.GetPieceBlockLocations(ctrl.X, ctrl.Y, ctrl.Piece)

	n := ctrl.Piece.BlockCount()
	events := make([]event.Event, 0, n)

	// it's important to start from the slice's end
	for i := len(blocks) - 1; i >= 0; i-- {
		xyb := blocks[i]
		var evs event.Slice

		switch xyb.Type {
		case block.TypeRock:
			_dropBlock(f, &evs, xyb.X, xyb.Y, xyb.Block)
		case block.TypeLava:
			_dropLava(f, &evs, xyb.X, xyb.Y)
		case block.TypeAcid:
			_dropAcid(f, &evs, xyb.X, xyb.Y)
		case block.TypeCurl:
			_dropCurl(f, &evs, xyb.X, xyb.Y)
		case block.TypeWave:
			_dropWave(f, &evs, xyb.X, xyb.Y)
		default:
			panic("unsupported block type found in piece")
		}

		if len(evs) == 0 {
			continue
		}

		for _, e := range evs {
			// push the event
			p.Push(e)

			// apply the event to the field
			switch v := e.(type) {
			case *op.FieldBlockSet:
				if v.Op == op.TypeSet {
					f.SetXY(int(v.Col), int(v.Row), field.AnimNo, 0, v.Block)
				} else {
					f.ClearXY(int(v.Col), int(v.Row), field.AnimNo, 0)
				}
			case *op.FieldBlockHardness:
				f.HardnessXY(int(v.Col), int(v.Row), int(v.Hardness), field.AnimNo, 0)
			case *op.FieldExBlock:
			default:
				panic("unsupported event type=" + reflect.TypeOf(e).String() + " for piece meld")
			}

			events = append(events, e)
		}
	}

	// undo the changes to the field
	for i := len(events) - 1; i >= 0; i-- {
		switch v := events[i].(type) {
		case *op.FieldBlockSet:
			if v.Op == op.TypeSet {
				f.ClearXY(int(v.Col), int(v.Row), field.AnimNo, 0)
			} else {
				f.SetXY(int(v.Col), int(v.Row), field.AnimNo, 0, v.Block)
			}
		case *op.FieldBlockHardness:
			f.HardnessXY(int(v.Col), int(v.Row), -int(v.Hardness), field.AnimNo, 0)
		}
	}
}

func _setBlock(f *field.Field, p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldBlockSet(x, y, op.TypeSet, field.AnimMeld, 0, b))
}

func _dropEx(f *field.Field, p event.Pusher, x, y int, b block.Block) {
	height, y0, ok := _dropHeight(f, x, y)
	if ok {
		p.Push(op.NewFieldExBlock(x, y0, field.AnimFall, height, b))
	}
}

func _destroyEx(f *field.Field, p event.Pusher, x, y int, b block.Block) {
	p.Push(op.NewFieldExBlock(x, y, field.AnimDestroy, 0, b))
}

func _dropBlock(f *field.Field, p event.Pusher, x, y int, b block.Block) bool {
	height, y0, ok := _dropHeight(f, x, y)
	if !ok {
		return false
	}

	p.Push(op.NewFieldBlockSet(x, y0, op.TypeSet, field.AnimFall, height, b))
	return true
}

func _dropLava(f *field.Field, p event.Pusher, x, y int) bool {
	height, y0, ok := _dropHeight(f, x, y)
	if !ok {
		return false
	}

	b := block.Block{Type: block.TypeRock, Color: block.Lava.Color}
	p.Push(op.NewFieldBlockSet(x, y0, op.TypeSet, field.AnimFall, height, b))
	return true
}

func _dropAcid(f *field.Field, p event.Pusher, x, y int) bool {
	height, y0, b, ok := _dropHeightToFull(f, x, y)
	if !ok {
		_dropEx(f, p, x, y, block.Acid)
		return false
	}

	if h := y - y0; h > 0 {
		p.Push(op.NewFieldExBlock(x, y0, field.AnimFall, h, block.Acid))
	}

	if b.Hardness == block.HardnessMax || !b.Type.Shootable() {
		return true
	}

	if b.Hardness > 0 {
		p.Push(op.NewFieldBlockHardness(x, y0, -1, field.AnimSpin, height))
	} else {
		p.Push(op.NewFieldBlockSet(x, y0, op.TypeClear, field.AnimPop, 0, b))
	}

	return true
}

func _dropCurl(f *field.Field, p event.Pusher, x, y int) bool {
	b := block.Block{Type: block.TypeRock, Color: block.Curl.Color}

	height, y0, ok := _dropHeightToHighestHole(f, x, y)
	if !ok {
		_destroyEx(f, p, x, y, b)
		return false
	}

	p.Push(op.NewFieldBlockSet(x, y0, op.TypeSet, field.AnimFall, height, b))
	return true
}

func _dropWave(f *field.Field, p event.Pusher, x, y int) bool {
	b := block.Block{Type: block.TypeRock, Color: block.Wave.Color}

	height, y0, ok := _dropHeightToLowestHole(f, x, y)
	if !ok {
		_destroyEx(f, p, x, y, b)
		return false
	}

	p.Push(op.NewFieldBlockSet(x, y0, op.TypeSet, field.AnimFall, height, b))
	return true
}

func _shoot(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	ammo := ctrl.Piece.ActivationCount()
	if ammo == 0 {
		return
	}

	b := ctrl.Piece.Get(0, 0)
	var hit bool

	switch b.Type {
	case block.TypeRock:
		hit = _shootBlock(f, ctrl, p, b)
	case block.TypeAcid:
		hit = _shootAcid(f, ctrl, p)
	case block.TypeLava:
		hit = _shootLava(f, ctrl, p)
	case block.TypeCurl:
		hit = _shootCurl(f, ctrl, p)
	case block.TypeWave:
		hit = _shootWave(f, ctrl, p)
	}

	p.Push(op.NewPieceShoot(ctrl.Idx, hit, b.Type))

	if ammo == 1 {
		_clearPiece(ctrl, p)
	}
}

func _shootBlock(f *field.Field, ctrl *piece.Ctrl, p event.Pusher, b block.Block) bool {
	return _dropBlock(f, p, ctrl.X, ctrl.Y, b)
}

func _shootAcid(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) bool {
	return _dropAcid(f, p, ctrl.X, ctrl.Y)
}

func _shootLava(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) bool {
	return _dropLava(f, p, ctrl.X, ctrl.Y)
}

func _shootCurl(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) bool {
	return _dropCurl(f, p, ctrl.X, ctrl.Y)
}

func _shootWave(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) bool {
	return _dropWave(f, p, ctrl.X, ctrl.Y)
}

func _dropHeight(f *field.Field, x, y int) (height, y0 int, ok bool) {
	height = f.GetHeightToTopmostEmpty(x, y)

	y0 = y - height
	b := f.GetXY(x, y0)

	if b.Type != block.TypeEmpty {
		return
	}

	ok = true

	return
}

func _dropHeightToFull(f *field.Field, x, y int) (height, y0 int, b block.Block, ok bool) {
	height = f.GetHeightToTopmostFull(x, y)
	if height == 0 {
		return
	}

	y0 = y - height
	b = f.GetXY(x, y0)
	ok = true

	return
}

func _dropHeightToHighestHole(f *field.Field, x, y int) (height, y0 int, ok bool) {
	height = f.GetHeightToHighestHole(x, y)
	if height == 0 {
		return
	}

	y0 = y - height
	ok = true

	return
}

func _dropHeightToLowestHole(f *field.Field, x, y int) (height, y0 int, ok bool) {
	height = f.GetHeightToLowestHole(x, y)
	if height == 0 {
		return
	}

	y0 = y - height
	ok = true

	return
}
