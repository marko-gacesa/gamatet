// Copyright (c) 2025 by Marko Gaćeša

package app

import (
	"fmt"

	"github.com/marko-gacesa/gamatet/game/setup"
	"github.com/marko-gacesa/gamatet/internal/config/key"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
)

func pieceTypeStr(t setup.PieceType) string {
	switch t {
	case setup.PieceTypeRotatingPolyominoes:
		return T(KeyPieceTypeRotatingPolyominoes)
	case setup.PieceTypeHMirroringPolyominoes:
		return T(KeyPieceTypeHMirroringPolyominoes)
	case setup.PieceTypeVMirroringPolyominoes:
		return T(KeyPieceTypeVMirroringPolyominoes)
	default:
		return "?"
	}
}

func pieceSizeStr(s byte) string {
	switch s {
	case setup.PieceSize3:
		return T(KeyPieceSize3)
	case setup.PieceSize4:
		return T(KeyPieceSize4)
	case setup.PieceSize5:
		return T(KeyPieceSize5)
	default:
		return "?"
	}
}

func rotationDirCWStr(dir bool) string {
	if dir {
		return T(KeyRotationDirCW)
	} else {
		return T(KeyRotationDirCCW)
	}
}

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
	sections.showFieldMap[false] = fmt.Sprintf("%s (%s)", T(KeySetupHide), s.FieldOptions.String(s.GameOptions.TeamSize))
	sections.showFieldMap[true] = T(KeySetupShow)
	sections.showPieceMap[false] = fmt.Sprintf("%s (%s)", T(KeySetupHide), s.PieceOptions.String())
	sections.showPieceMap[true] = T(KeySetupShow)
	sections.showMiscMap[false] = fmt.Sprintf("%s (%s)", T(KeySetupHide), s.MiscOptions.String())
	sections.showMiscMap[true] = T(KeySetupShow)
}

func (sections *setupSections) showFieldsStr(b bool) string { return sections.showFieldMap[b] }
func (sections *setupSections) showPieceStr(b bool) string  { return sections.showPieceMap[b] }
func (sections *setupSections) showMiscStr(b bool) string   { return sections.showMiscMap[b] }

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
	k.showKeysMap[true] = T(KeySetupRedefine)
}

func (k *setupKeySection) showKeysStr(b bool) string { return k.showKeysMap[b] }
