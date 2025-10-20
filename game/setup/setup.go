// Copyright (c) 2025 by Marko Gaćeša

package setup

import (
	"fmt"
	"gamatet/game/field"
	"gamatet/game/piece"
	"github.com/marko-gacesa/bitdata"
	"github.com/marko-gacesa/udpstar/udpstar/message"
	"github.com/marko-gacesa/udpstar/udpstar/message/lobby"
	"slices"
	"strconv"
	"strings"
)

const (
	MaxPlayers      = 8
	MaxLocalPlayers = 4

	MaxLenName = lobby.MaxLenName

	MaxFieldCount = 8
	MaxTeamSize   = field.MaxPieces

	MaxFieldWidthSingle    = field.MaxWidth
	MaxFieldWidthPerPlayer = field.MaxWidth / field.MaxPieces
	MinFieldWidthPerPlayer = field.MinWidth
	DefaultFieldWidth      = 10

	MaxFieldHeight     = field.MaxHeight
	MinFieldHeight     = field.MinHeight
	DefaultFieldHeight = 20

	MaxSpeed     = piece.MaxLevel
	MinSpeed     = 0
	DefaultSpeed = 4
)

type Setup struct {
	Name string
	GameOptions
	FieldOptions
	PieceOptions
	MiscOptions
}

func (o *Setup) Def() []byte {
	w := bitdata.NewWriter()
	o.Write(w)
	gameDef := w.BitData()
	return gameDef
}

func (o *Setup) Write(w *bitdata.Writer) {
	o.GameOptions.Write(w)
	o.FieldOptions.Write(w)
	o.PieceOptions.Write(w)
	o.MiscOptions.Write(w)
}

func (o *Setup) Read(r *bitdata.ReaderError) {
	o.GameOptions.Read(r)
	o.FieldOptions.Read(r)
	o.PieceOptions.Read(r)
	o.MiscOptions.Read(r)
}

func (o *Setup) Sanitize() bool {
	s1 := o.GameOptions.Sanitize()
	s2 := o.FieldOptions.Sanitize()
	s3 := o.PieceOptions.Sanitize()
	s4 := o.MiscOptions.Sanitize()
	return s1 || s2 || s3 || s4
}

func (o *Setup) String() string {
	sb := strings.Builder{}

	sb.WriteString(o.GameOptions.String())
	sb.WriteByte(' ')
	sb.WriteString(o.FieldOptions.String(o.GameOptions.TeamSize))

	if s := o.PieceOptions.String(); s != "" {
		sb.WriteByte(' ')
		sb.WriteString(s)
	}

	if s := o.MiscOptions.String(); s != "" {
		sb.WriteByte(' ')
		sb.WriteString(s)
	}

	return sb.String()
}

type GameOptions struct {
	GameType         GameType `json:"game_type"`
	FieldCount       byte     `json:"field_count"` // 1..MaxFieldCount
	TeamSize         byte     `json:"team_size"`   // 1..MaxTeamSize
	PlayerZones      bool     `json:"player_zones"`
	PieceCollision   bool     `json:"piece_collision"`
	SamePiecesForAll bool     `json:"same_pieces_for_all"`
}

func (o *GameOptions) PlayerCount() byte {
	return o.FieldCount * o.TeamSize
}

func (o *GameOptions) CreateSlotsStories() []message.Token {
	fieldCount := o.FieldCount
	teamSize := o.TeamSize
	slotStories := make([]message.Token, fieldCount*teamSize)
	for i := range fieldCount {
		storyToken := message.RandomToken()
		for j := range teamSize {
			slotStories[j+i*teamSize] = storyToken
		}
	}

	return slotStories
}

func (o *GameOptions) Write(w *bitdata.Writer) {
	w.Write8(byte(o.GameType), 4)
	w.Write8(o.FieldCount-1, 3)
	w.Write8(o.TeamSize-1, 2)
	w.WriteBool(o.PlayerZones)
	w.WriteBool(o.PieceCollision)
	w.WriteBool(o.SamePiecesForAll)
}

func (o *GameOptions) Read(r *bitdata.ReaderError) {
	o.GameType = GameType(r.Read8(4))
	o.FieldCount = r.Read8(3) + 1
	o.TeamSize = r.Read8(2) + 1
	o.PlayerZones = r.ReadBool()
	o.PieceCollision = r.ReadBool()
	o.SamePiecesForAll = r.ReadBool()
}

func (o *GameOptions) Sanitize() bool {
	var sanitized bool
	if o.GameType > GameTypeVMirroringPolyominoes {
		o.GameType = GameTypeRotatingPolyominoes
		sanitized = true
	}
	if o.FieldCount < 1 || o.FieldCount > MaxFieldCount {
		o.FieldCount = 1
		sanitized = true
	}
	if o.TeamSize < 1 || o.TeamSize > MaxTeamSize {
		o.TeamSize = MaxTeamSize
		sanitized = true
	}
	if o.PieceCollision && o.TeamSize == 1 {
		o.PieceCollision = false
		sanitized = true
	}
	if o.PlayerZones && o.TeamSize == 1 {
		o.PlayerZones = false
		sanitized = true
	}
	return sanitized
}

func (o *GameOptions) String() string {
	const smileyEmpty = "☺"
	const smileyFull = "☻"

	sb := strings.Builder{}

	sb.WriteString(func() string {
		players := o.PlayerCount()
		if o.FieldCount == 1 {
			if o.TeamSize == 1 {
				return "1" + smileyFull
			}
			return fmt.Sprintf("%d%s co-op", players, smileyFull)
		}

		if o.TeamSize == 1 {
			return fmt.Sprintf("%d%s battle", players, smileyFull)
		}

		return fmt.Sprintf("%d☻ ", players) +
			strings.Join(slices.Repeat([]string{string('0' + o.TeamSize)}, int(o.FieldCount)), "vs")
	}())

	if o.TeamSize > 1 && o.PlayerZones {
		sb.WriteString(" ZONES")
	}

	if o.TeamSize > 1 && !o.PlayerZones && o.PieceCollision {
		sb.WriteString(" COLLIDE")
	}

	if o.FieldCount*o.TeamSize > 1 && o.SamePiecesForAll {
		sb.WriteString(" SAME")
	}

	return sb.String()
}

type FieldOptions struct {
	WidthSingle    byte
	WidthPerPlayer byte
	Height         byte
	Speed          byte
}

func (o *FieldOptions) Write(w *bitdata.Writer) {
	w.Write8(o.WidthSingle-MinFieldWidthPerPlayer, 6)    // 4..40 -> 0..36 -> 6 bits
	w.Write8(o.WidthPerPlayer-MinFieldWidthPerPlayer, 3) // 4..10 -> 0..6 -> 3 bits
	w.Write8(o.Height-MinFieldHeight, 6)                 // 4..40 -> 0..36 -> 6 bits
	w.Write8(o.Speed, 4)                                 // 0..15 -> 4 bits
}

func (o *FieldOptions) Read(r *bitdata.ReaderError) {
	o.WidthSingle = r.Read8(6) + MinFieldWidthPerPlayer
	o.WidthPerPlayer = r.Read8(3) + MinFieldWidthPerPlayer
	o.Height = r.Read8(6) + MinFieldHeight
	o.Speed = r.Read8(4)
}

func (o *FieldOptions) Sanitize() bool {
	var sanitized bool
	if o.WidthSingle < MinFieldWidthPerPlayer || o.WidthSingle > MaxFieldWidthSingle {
		o.WidthSingle = DefaultFieldWidth
		sanitized = true
	}
	if o.WidthPerPlayer < MinFieldWidthPerPlayer || o.WidthPerPlayer > MaxFieldWidthPerPlayer {
		o.WidthPerPlayer = DefaultFieldWidth
		sanitized = true
	}
	if o.Height < MinFieldHeight || o.Height > MaxFieldHeight {
		o.Height = DefaultFieldHeight
		sanitized = true
	}
	if o.Speed > MaxSpeed {
		o.Speed = DefaultSpeed
		sanitized = true
	}
	return sanitized
}

func (o *FieldOptions) String(teamSize byte) string {
	sb := strings.Builder{}

	if teamSize == 1 {
		sb.WriteString(strconv.Itoa(int(o.WidthSingle)))
	} else {
		sb.WriteString(strconv.Itoa(int(o.WidthPerPlayer)))
	}

	sb.WriteByte('x')

	sb.WriteString(strconv.Itoa(int(o.Height)))

	sb.WriteByte('@')

	sb.WriteString(strconv.Itoa(int(o.Speed)))

	return sb.String()
}

type PieceOptions struct {
	/*
		Monomino    bool
		Dominoes    bool
		Trominoes   bool
		Tetrominoes bool
		Pentominoes bool

		Shooters bool
		Special  bool
		Lava     bool
		Acid     bool
		Curl     bool
		Wave     bool
	*/
}

func (o *PieceOptions) Write(w *bitdata.Writer) {
}

func (o *PieceOptions) Read(r *bitdata.ReaderError) {
}

func (o *PieceOptions) Sanitize() bool {
	return false
}

func (o *PieceOptions) String() string {
	return ""
}

type MiscOptions struct {
	CustomSeed bool
	Seed       int64
}

func (o *MiscOptions) Write(w *bitdata.Writer) {
	w.WriteBool(o.CustomSeed)
	w.Write64(uint64(o.Seed), 64)
}

func (o *MiscOptions) Read(r *bitdata.ReaderError) {
	o.CustomSeed = r.ReadBool()
	o.Seed = int64(r.Read64(64))
}

func (o *MiscOptions) Sanitize() bool {
	return false
}

func (o *MiscOptions) String() string {
	sb := strings.Builder{}

	if o.CustomSeed {
		sb.WriteString("seed=")
		sb.WriteString(strconv.FormatInt(o.Seed, 10))
	}

	return sb.String()
}

func SinglePlayerSetupDefault() Setup {
	return Setup{
		Name: "",
		GameOptions: GameOptions{
			GameType:         GameTypeRotatingPolyominoes,
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
			//TODO
		},
		MiscOptions: MiscOptions{},
	}
}

func MultiplayerPlayerSetupDefault() Setup {
	return Setup{
		Name: "",
		GameOptions: GameOptions{
			GameType:         GameTypeRotatingPolyominoes,
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
			//TODO
		},
		MiscOptions: MiscOptions{},
	}
}

type PlayerConfig piece.Config

func (c *PlayerConfig) Write(w *bitdata.Writer) {
	w.WriteBool(c.RotationDirectionCW)
	w.WriteBool(c.SlideDisabled)
	w.Write8(c.WallKick, 2)
}

func (c *PlayerConfig) Read(r *bitdata.ReaderError) {
	c.RotationDirectionCW = r.ReadBool()
	c.SlideDisabled = r.ReadBool()
	c.WallKick = r.Read8(2)
}
