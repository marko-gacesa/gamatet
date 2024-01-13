// Copyright (c) 2020-2024 by Marko Gaćeša

package machine

import (
	"gamatet/game/action"
	"gamatet/game/block"
	"gamatet/game/event"
	"gamatet/game/field"
	"gamatet/game/op"
	"gamatet/game/piece"
	"reflect"
	"time"
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
			(ctrl.State != piece.StateSlide && ctrl.State != piece.StateFall || !ctrl.SlideEnabled) {
			return
		}

		var dx int
		if a == action.MoveLeft {
			dx = -1
		} else {
			dx = 1
		}

		success := f.CanMovePiece(dx, 0, ctrl.Idx, !f.PieceCollision)
		if !success {
			break
		}

		p.Push(op.NewPieceMove(ctrl.Idx, dx, 0))

	case action.RotateCW, action.RotateCCW:
		if ctrl.State != piece.StateDescend && ctrl.State != piece.StateSlide {
			break
		}

		pieceType := ctrl.Piece.Type()
		if pieceType == piece.TypeShooter {
			_shoot(f, ctrl, p)
			break // shooters can't rotate
		}

		dirCW := a == action.RotateCW

		success, _, dx, _ := f.CanRotatePiece(dirCW, ctrl.Idx, !f.PieceCollision)
		if !success {
			break
		}

		if dx != 0 {
			p.Push(op.NewPieceMove(ctrl.Idx, dx, 0))
		}

		p.Push(op.NewPieceRotate(ctrl.Idx, dirCW))

	case action.MoveDown:
		if ctrl.State != piece.StateDescend && ctrl.State != piece.StateSlide {
			break
		}

		success := f.CanMovePiece(0, -1, ctrl.Idx, !f.PieceCollision)
		if !success {
			if ctrl.State != piece.StateSlide {
				_changeState(ctrl, p, piece.StateSlide)
			} else {
				ctrl.Timer.Reset(time.Nanosecond)
			}
			break
		}

		p.Push(op.NewPieceMove(ctrl.Idx, 0, -1))

		_changeState(ctrl, p, piece.StateDescend)

	case action.Drop:
		if ctrl.State != piece.StateDescend && ctrl.State != piece.StateSlide {
			break
		}

		if ctrl.Piece.Type() == piece.TypeShooter {
			_shoot(f, ctrl, p)
			break
		}

		if ctrl.Piece.Type() == piece.TypeStandard {
			if t := ctrl.Blocks[0].Type; t == block.TypeLava || t == block.TypeAcid || t == block.TypeWave {
				_meldPiece(f, ctrl, p)
				_clearPiece(ctrl, p)
				break
			}
		}

		height := f.GetDropHeight(ctrl.Idx, !f.PieceCollision)
		if height == 0 {
			if ctrl.State != piece.StateSlide {
				_changeState(ctrl, p, piece.StateSlide)
			}
			break
		}

		p.Push(op.NewPieceFall(ctrl.Idx, height))

		_changeStateWithParam(ctrl, p, piece.StateFall, height)
	}
}

func HandleTimeout(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	switch ctrl.State {
	case piece.StateInit:
		_changeState(ctrl, p, piece.StateNew)

	case piece.StateNew:
		if ctrl.Piece != nil {
			// should not happen: player should not have a piece in this state
			panic("player already has a piece in state New")
		}

		pieceCount := ctrl.PieceCount

		newPiece := ctrl.Feed.Get(pieceCount)
		success, x, y := f.GetPieceStartPosition(ctrl.Idx, ctrl, newPiece, !f.PieceCollision)
		if !success && f.PieceCollision {
			success, x, y = f.GetPieceStartPosition(ctrl.Idx, ctrl, newPiece, true)
			if success {
				_changeState(ctrl, p, piece.StateNew) // wait awhile, there is a piece in the way
				break
			}
		}
		if !success {
			_changeState(ctrl, p, piece.StateLost) // can't position piece: end game
			break
		}

		p.Push(op.NewPieceSet(ctrl.Idx, op.OpSet, x, y, newPiece, pieceCount+1))

		_changeState(ctrl, p, piece.StateDescend)

	case piece.StateDescend, piece.StateSlide:
		if f.CanMovePiece(0, -1, ctrl.Idx, !f.PieceCollision) {
			p.Push(op.NewPieceFall(ctrl.Idx, 1))
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

	case piece.StateLost, piece.StateWon:
		_changeState(ctrl, p, piece.StateStop)

	case piece.StateStop:
		// should not happen, this state doesn't use timer
		panic("timer timeout for game state with no timer")

	default:
		// should not happen, unknown state
		panic("timer timeout for unknown state")
	}
}

func _changeState(ctrl *piece.Ctrl, p event.Pusher, newState piece.State) {
	p.Push(op.NewPieceState(ctrl.Idx, ctrl.State, newState, 0, 0))
}

func _changeStateWithParam(ctrl *piece.Ctrl, p event.Pusher, newState piece.State, param int) {
	p.Push(op.NewPieceState(ctrl.Idx, ctrl.State, newState, 0, param))
}

func _clearPiece(ctrl *piece.Ctrl, p event.Pusher) {
	p.Push(op.NewPieceSet(ctrl.Idx, op.OpClear, ctrl.X, ctrl.Y, ctrl.Piece, ctrl.PieceCount))
	_changeState(ctrl, p, piece.StateNew)
}

func _meldPiece(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	if ctrl.Piece.Type() == piece.TypeShooter {
		return
	}

	if ctrl.Piece.Type() == piece.TypeStandard {
		switch ctrl.Blocks[0].Type {
		case block.TypeLava, block.TypeAcid, block.TypeWave:
			_meldLiquidPiece(f, ctrl, p)
		default:
			_meldSolidPiece(f, ctrl, p)
		}
	}
}

func _meldSolidPiece(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	// blocks are returned from the lowest row, up to the topmost row
	blocks := f.GetPieceBlockLocations(ctrl.X, ctrl.Y, ctrl.Piece)

	for i := len(blocks) - 1; i >= 0; i-- {
		xyb := blocks[i]
		if e := _setBlockEvent(f, xyb.X, xyb.Y, xyb.Block); e != nil {
			p.Push(e)
		}
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
		var e event.Event
		switch xyb.Type {
		case block.TypeRock:
			e = _setBlockEvent(f, xyb.X, xyb.Y, xyb.Block)
		case block.TypeLava:
			e = _dropLavaEvent(f, xyb.X, xyb.Y)
		case block.TypeAcid:
			e = _dropAcidEvent(f, xyb.X, xyb.Y)
		case block.TypeWave:
			e = _dropWaveEvent(f, xyb.X, xyb.Y)
		default:
			panic("unsupported block type found in piece")
		}

		if e == nil {
			continue
		}

		// push the event
		p.Push(e)

		// apply the event to the field
		switch v := e.(type) {
		case *op.FieldBlockSet:
			if v.Op == op.OpSet {
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

	// undo the changes to the field
	for i := len(events) - 1; i >= 0; i-- {
		switch v := events[i].(type) {
		case *op.FieldBlockSet:
			if v.Op == op.OpSet {
				f.ClearXY(int(v.Col), int(v.Row), field.AnimNo, 0)
			} else {
				f.SetXY(int(v.Col), int(v.Row), field.AnimNo, 0, v.Block)
			}
		case *op.FieldBlockHardness:
			f.HardnessXY(int(v.Col), int(v.Row), -int(v.Hardness), field.AnimNo, 0)
		}
	}
}

func _setBlockEvent(f *field.Field, x, y int, b block.Block) event.Event {
	return op.NewFieldBlockSet(x, y, op.OpSet, field.AnimMeld, 0, b)
}

func _dropLavaEvent(f *field.Field, x, y int) event.Event {
	height, y0, ok := _dropLavaHeight(f, x, y)
	if !ok {
		return nil
	}

	b := block.Block{
		Type:  block.TypeRock,
		Color: block.Lava.Color,
	}

	return op.NewFieldBlockSet(x, y0, op.OpSet, field.AnimFall, height, b)
}

func _dropAcidEvent(f *field.Field, x, y int) event.Event {
	height, y0, b, ok := _dropAcidHeight(f, x, y)
	if !ok || b.Hardness == block.HardnessMax {
		return op.NewFieldExBlock(x, y, field.AnimPop, 0, block.Acid)
	}

	if b.Hardness > 0 {
		return op.NewFieldBlockHardness(x, y0, -1, field.AnimFall, height)
	}

	return op.NewFieldBlockSet(x, y0, op.OpClear, field.AnimFall, height, b)
}

func _dropWaveEvent(f *field.Field, x, y int) event.Event {
	height, y0, ok := _dropWaveHeight(f, x, y)
	if !ok {
		return op.NewFieldExBlock(x, y, field.AnimDestroy, 0, block.Wave)
	}

	b := block.Block{
		Type:  block.TypeRock,
		Color: block.Wave.Color,
	}

	return op.NewFieldBlockSet(x, y0, op.OpSet, field.AnimFall, height, b)
}

func _shoot(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	ammo := ctrl.Piece.GetActivationCount()
	if ammo == 0 {
		return
	}

	b := ctrl.Piece.Get(0, 0)
	switch b.Type {
	case block.TypeLava:
		_shootLava(f, ctrl, p)
	case block.TypeAcid:
		_shootAcid(f, ctrl, p)
	case block.TypeWave:
		_shootWave(f, ctrl, p)
	}

	if ammo == 1 {
		_clearPiece(ctrl, p)
	}
}

func _shootLava(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	height, y0, ok := _dropLavaHeight(f, ctrl.X, ctrl.Y)
	if !ok {
		return
	}

	b := block.Block{
		Type:  block.TypeRock,
		Color: block.Lava.Color,
	}

	p.Push(op.NewPieceActivate(ctrl.Idx, 1))
	p.Push(op.NewFieldBlockSet(ctrl.X, y0, op.OpSet, field.AnimShot, height, b))
}

func _shootAcid(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	height, y0, b, ok := _dropAcidHeight(f, ctrl.X, ctrl.Y)
	if !ok || b.Hardness == block.HardnessMax {
		return
	}

	p.Push(op.NewPieceActivate(ctrl.Idx, 1))

	if b.Hardness > 0 {
		p.Push(op.NewFieldBlockHardness(ctrl.X, y0, -1, field.AnimShot, height))
	} else {
		p.Push(op.NewFieldBlockSet(ctrl.X, y0, op.OpClear, field.AnimShot, height, b))
	}
}

func _shootWave(f *field.Field, ctrl *piece.Ctrl, p event.Pusher) {
	height, y0, ok := _dropWaveHeight(f, ctrl.X, ctrl.Y)
	if !ok {
		return
	}

	b := block.Block{
		Type:  block.TypeRock,
		Color: block.Wave.Color,
	}

	p.Push(op.NewPieceActivate(ctrl.Idx, 1))
	p.Push(op.NewFieldBlockSet(ctrl.X, y0, op.OpSet, field.AnimShot, height, b))
}

func _dropLavaHeight(f *field.Field, x, y int) (height, y0 int, ok bool) {
	height = f.GetHeightToTopmostEmpty(x, y)

	y0 = y - height
	b := f.GetXY(x, y0)

	if b.Type != block.TypeEmpty {
		return
	}

	ok = true

	return
}

func _dropAcidHeight(f *field.Field, x, y int) (height, y0 int, b block.Block, ok bool) {
	height = f.GetHeightToTopmostFull(x, y)
	if height == 0 {
		return
	}

	y0 = y - height
	b = f.GetXY(x, y0)
	ok = true

	return
}

func _dropWaveHeight(f *field.Field, x, y int) (height, y0 int, ok bool) {
	height = f.GetHeightToTopmostHole(x, y)
	if height == 0 {
		return
	}

	y0 = y - height
	ok = true

	return
}
