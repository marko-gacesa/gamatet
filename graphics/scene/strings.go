// Copyright (c) 2025, 2026 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package scene

import (
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/graphics/render"
	. "github.com/marko-gacesa/gamatet/internal/i18n"
)

func fieldStrings() render.FieldStrings {
	return render.FieldStrings{
		TitlePanel: struct {
			Blocks string
		}{
			Blocks: T(KeyFieldTitlePanelBlocks),
		},
		SidePanel: struct {
			Player string
			Score  string
			Piece  string
			Level  string
			Next   string
		}{
			Player: T(KeyFieldSidePanelPlayer),
			Score:  T(KeyFieldSidePanelScore),
			Piece:  T(KeyFieldSidePanelPiece),
			Level:  T(KeyFieldSidePanelLevel),
			Next:   T(KeyFieldSidePanelNext),
		},
		Message: struct {
			GameOver   string
			Victory    string
			Defeat     string
			Pause      string
			Suspended  string
			ServerLost string
		}{
			GameOver:   T(KeyFieldMessageGameOver),
			Victory:    T(KeyFieldMessageVictory),
			Defeat:     T(KeyFieldMessageDefeat),
			Pause:      T(KeyFieldMessagePause),
			Suspended:  T(KeyFieldMessageSuspended),
			ServerLost: T(KeyFieldMessageServerLost),
		},
		EffectMap: map[field.Effect]string{
			field.EffectLid:   T(KeyEffectLid),
			field.EffectBigO:  T(KeyEffectBigO),
			field.EffectRaise: T(KeyEffectRaise),
			field.EffectGnaw:  T(KeyEffectGnaw),
			field.EffectPatch: T(KeyEffectPatch),
		},
	}
}
