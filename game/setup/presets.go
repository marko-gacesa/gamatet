// Copyright (c) 2025 by Marko Gaćeša

package setup

const (
	SinglePlayerPresetCount = 5
	MultiPlayerPresetCount  = 2
)

func SinglePlayerPreset(idx int) Setup {
	switch byte(idx) % SinglePlayerPresetCount {
	case 0:
		return SinglePlayerPresetRotTetrominoes()
	case 1:
		return SinglePlayerPresetRotPentominoes()
	case 2:
		return SinglePlayerPresetRotMiniminoes()
	case 3:
		return SinglePlayerPresetVMTetrominoes()
	case 4:
		return SinglePlayerPresetHMTetrominoes()
	}
	panic("unreachable")
}

func SinglePlayerPresetRotTetrominoes() Setup {
	return Setup{
		Name: "Rotating Tetrominoes (Classic)",
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
		Name: "Rotating Pentominoes",
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
		Name: "Rotating Miniminoes",
		GameOptions: GameOptions{
			GameType:         GameTypeFallingPolyominoes,
			FieldCount:       1,
			TeamSize:         1,
			PieceCollision:   false,
			PlayerZones:      false,
			SamePiecesForAll: false,
		},
		FieldOptions: FieldOptions{
			WidthSingle:    DefaultFieldWidth - 4,
			WidthPerPlayer: DefaultFieldWidth,
			Height:         DefaultFieldHeight,
			Speed:          9,
		},
		PieceOptions: PieceOptions{
			PieceType: PieceTypeRotatingPolyominoes,
			PieceSize: PieceSize3,
			BagSize:   BagSizeDefault,
		},
		MiscOptions: MiscOptions{},
	}
}

func SinglePlayerPresetVMTetrominoes() Setup {
	return Setup{
		Name: "V-Mirroring Tetrominoes",
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
			PieceType: PieceTypeVMirroringPolyominoes,
			PieceSize: PieceSize4,
			BagSize:   BagSizeDefault,
		},
		MiscOptions: MiscOptions{},
	}
}

func SinglePlayerPresetHMTetrominoes() Setup {
	return Setup{
		Name: "H-Mirroring Tetrominoes",
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
			PieceType: PieceTypeHMirroringPolyominoes,
			PieceSize: PieceSize4,
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
	}
	panic("unreachable")
}

func MultiPlayerPresetCoop2() Setup {
	return Setup{
		Name: "Two Players: Cooperative",
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
		Name: "Two Players: Battle",
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
