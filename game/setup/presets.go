// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package setup

const (
	SinglePlayerPresetCount = 3
	MultiPlayerPresetCount  = 3
)

func SinglePlayerPreset(idx int) Setup {
	switch byte(idx) % SinglePlayerPresetCount {
	case 0:
		return SinglePlayerPresetRotTetrominoes()
	case 1:
		return SinglePlayerPresetRotPentominoes()
	case 2:
		return SinglePlayerPresetRotMiniminoes()
	}
	panic("unreachable")
}

func SinglePlayerPresetRotTetrominoes() Setup {
	return Setup{
		Name: "Classic",
		GameOptions: GameOptions{
			GameType:         GameTypeFallingPolyominoes,
			FieldCount:       1,
			TeamSize:         1,
			PieceCollision:   false,
			PlayerZones:      false,
			SamePiecesForAll: false,
		},
		FieldOptions: FieldOptions{
			WidthSingle:    DefaultFieldWidth,
			WidthPerPlayer: DefaultFieldWidth,
			Height:         DefaultFieldHeight,
			Speed:          DefaultSpeed,
		},
		PieceOptions: PieceOptions{
			PieceType: PieceTypeRotatingPolyominoes,
			PieceSize: PieceSize4,
			BagSize:   BagSizeDefault,
		},
		MiscOptions: MiscOptions{},
	}
}

func SinglePlayerPresetRotPentominoes() Setup {
	return Setup{
		Name: "Pentix",
		GameOptions: GameOptions{
			GameType:         GameTypeFallingPolyominoes,
			FieldCount:       1,
			TeamSize:         1,
			PieceCollision:   false,
			PlayerZones:      false,
			SamePiecesForAll: false,
		},
		FieldOptions: FieldOptions{
			WidthSingle:    DefaultFieldWidth + 2,
			WidthPerPlayer: DefaultFieldWidth,
			Height:         DefaultFieldHeight,
			Speed:          DefaultSpeed,
		},
		PieceOptions: PieceOptions{
			PieceType: PieceTypeRotatingPolyominoes,
			PieceSize: PieceSize5,
			BagSize:   BagSizeDefault,
		},
		MiscOptions: MiscOptions{},
	}
}

func SinglePlayerPresetRotMiniminoes() Setup {
	return Setup{
		Name: "Speedway",
		GameOptions: GameOptions{
			GameType:         GameTypeFallingPolyominoes,
			FieldCount:       1,
			TeamSize:         1,
			PieceCollision:   false,
			PlayerZones:      false,
			SamePiecesForAll: false,
		},
		FieldOptions: FieldOptions{
			WidthSingle:    MinFieldWidthPerPlayer,
			WidthPerPlayer: MinFieldWidthPerPlayer,
			Height:         MaxFieldHeight,
			Speed:          10,
		},
		PieceOptions: PieceOptions{
			PieceType: PieceTypeRotatingPolyominoes,
			PieceSize: PieceSize3,
			BagSize:   BagSizeDefault,
		},
		MiscOptions: MiscOptions{},
	}
}

func MultiPlayerPreset(idx int) Setup {
	switch byte(idx) % MultiPlayerPresetCount {
	case 0:
		return MultiPlayerPresetCoop2()
	case 1:
		return MultiPlayerPresetBattle1vs1()
	case 2:
		return MultiPlayerPresetBattle2vs2()
	}
	panic("unreachable")
}

func MultiPlayerPresetCoop2() Setup {
	return Setup{
		Name: "Co-op",
		GameOptions: GameOptions{
			GameType:         GameTypeFallingPolyominoes,
			FieldCount:       1,
			TeamSize:         2,
			PieceCollision:   true,
			PlayerZones:      false,
			SamePiecesForAll: true,
		},
		FieldOptions: FieldOptions{
			WidthSingle:    DefaultFieldWidth,
			WidthPerPlayer: DefaultFieldWidth - 2,
			Height:         DefaultFieldHeight,
			Speed:          DefaultSpeed,
		},
		PieceOptions: PieceOptions{
			PieceType: PieceTypeRotatingPolyominoes,
			PieceSize: PieceSize4,
			BagSize:   BagSizeDefault,
		},
		MiscOptions: MiscOptions{},
	}
}

func MultiPlayerPresetBattle1vs1() Setup {
	return Setup{
		Name: "Battle",
		GameOptions: GameOptions{
			GameType:         GameTypeFallingPolyominoes,
			FieldCount:       2,
			TeamSize:         1,
			PieceCollision:   false,
			PlayerZones:      false,
			SamePiecesForAll: true,
		},
		FieldOptions: FieldOptions{
			WidthSingle:    DefaultFieldWidth,
			WidthPerPlayer: DefaultFieldWidth,
			Height:         DefaultFieldHeight,
			Speed:          DefaultSpeed,
		},
		PieceOptions: PieceOptions{
			PieceType: PieceTypeRotatingPolyominoes,
			PieceSize: PieceSize4,
			BagSize:   BagSizeDefault,
		},
		MiscOptions: MiscOptions{},
	}
}

func MultiPlayerPresetBattle2vs2() Setup {
	return Setup{
		Name: "2 versus 2",
		GameOptions: GameOptions{
			GameType:         GameTypeFallingPolyominoes,
			FieldCount:       2,
			TeamSize:         2,
			PieceCollision:   false,
			PlayerZones:      true,
			SamePiecesForAll: true,
		},
		FieldOptions: FieldOptions{
			WidthSingle:    DefaultFieldWidth,
			WidthPerPlayer: DefaultFieldWidth,
			Height:         DefaultFieldHeight,
			Speed:          DefaultSpeed,
		},
		PieceOptions: PieceOptions{
			PieceType: PieceTypeRotatingPolyominoes,
			PieceSize: PieceSize4,
			BagSize:   BagSizeDefault,
		},
		MiscOptions: MiscOptions{},
	}
}
