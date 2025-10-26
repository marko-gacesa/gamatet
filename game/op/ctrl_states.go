// Copyright (c) 2025 by Marko Gaćeša

package op

import (
	"gamatet/game/field"
	"gamatet/game/piece"
	"io"
)

func saveCtrlStates(f *field.Field) []byte {
	n := byte(f.Ctrls())
	ctrlStates := make([]byte, 2*n)
	for i := byte(0); i < n; i++ {
		state := f.Ctrl(i).State
		ctrlStates[i*2] = i
		ctrlStates[i*2+1] = byte(state)
	}
	return ctrlStates
}

func setCtrlStates(f *field.Field, ctrlStates []byte, newState piece.State) {
	for i := 0; i < len(ctrlStates); i += 2 {
		ctrlIdx := ctrlStates[i]
		ctrl := f.Ctrl(ctrlIdx)
		ctrl.State = newState
		ctrl.RestartTimer(0)
	}
}

func restoreCtrlStates(f *field.Field, ctrlStates []byte) {
	for i := 0; i < len(ctrlStates); i += 2 {
		ctrlIdx := ctrlStates[i]
		ctrlState := ctrlStates[i+1]
		ctrl := f.Ctrl(ctrlIdx)
		ctrl.State = piece.State(ctrlState)
		ctrl.RestartTimer(0)
	}
}

func writeCtrlStates(w io.Writer, ctrlStates []byte) error {
	if _, err := w.Write([]byte{byte(len(ctrlStates))}); err != nil {
		return err
	}
	if _, err := w.Write(ctrlStates); err != nil {
		return err
	}
	return nil
}

func readCtrlStates(r io.Reader, ctrlStates *[]byte) error {
	var buffer [1]byte
	if _, err := io.ReadFull(r, buffer[:]); err != nil {
		return err
	}

	*ctrlStates = make([]byte, buffer[0])
	if _, err := io.ReadFull(r, *ctrlStates); err != nil {
		return err
	}

	return nil
}
