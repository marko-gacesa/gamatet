// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"

	"github.com/marko-gacesa/gamatet/internal/config/key"

	"github.com/marko-gacesa/gamatet/game/setup"
)

type setupSections struct {
	showField    bool
	showPiece    bool
	showMisc     bool
	showFieldMap map[bool]string
	showPieceMap map[bool]string
	showMiscMap  map[bool]string
}

func newSetupSections() *setupSections {
	return &setupSections{
		showField:    false,
		showPiece:    false,
		showMisc:     false,
		showFieldMap: make(map[bool]string),
		showPieceMap: make(map[bool]string),
		showMiscMap:  make(map[bool]string),
	}
}

func (sections *setupSections) refresh(s *setup.Setup) {
	sections.showFieldMap[false] = fmt.Sprintf("%s (%s)", "Hide", s.FieldOptions.String(s.GameOptions.TeamSize))
	sections.showFieldMap[true] = "Show"
	sections.showPieceMap[false] = fmt.Sprintf("%s (%s)", "Hide", s.PieceOptions.String())
	sections.showPieceMap[true] = "Show"
	sections.showMiscMap[false] = fmt.Sprintf("%s (%s)", "Hide", s.MiscOptions.String())
	sections.showMiscMap[true] = "Show"
}

type setupKeySection struct {
	showKeys    bool
	showKeysMap map[bool]string
}

func newSetupKeySection() *setupKeySection {
	return &setupKeySection{
		showKeys:    false,
		showKeysMap: make(map[bool]string),
	}
}

func (k *setupKeySection) refresh(i *key.Input) {
	k.showKeysMap[false] = i.String()
	k.showKeysMap[true] = "Redefine"
}
