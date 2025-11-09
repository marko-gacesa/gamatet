// Copyright (c) 2025 by Marko Gaćeša

package setup

import (
	"fmt"
	"github.com/marko-gacesa/bitdata"
	"github.com/marko-gacesa/gamatet/game/field"
	"github.com/marko-gacesa/gamatet/game/piece"
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

	MaxFieldWidthSingle    = 20
	MaxFieldWidthPerPlayer = 10
	MinFieldWidthPerPlayer = 6
	DefaultFieldWidth      = 10

	MaxFieldHeight     = 24
	MinFieldHeight     = 12
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

func (o *Setup) Empty() bool {
	return o.GameOptions.FieldCount == 0 || o.GameOptions.TeamSize == 0 ||
		o.FieldOptions.WidthSingle == 0 || o.FieldOptions.WidthPerPlayer == 0 || o.FieldOptions.Height == 0
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

func (o *Setup) SanitizeName() bool {
	if o.Name == "" {
		o.Name = "Game"
		return true
	}

	runeCount := 0
	for i := range o.Name {
		if runeCount == MaxLenName {
			o.Name = o.Name[:i]
			return true
		}
		runeCount++
	}

	return false

}

func (o *Setup) Sanitize() bool {
	s1 := o.GameOptions.Sanitize()
	s2 := o.FieldOptions.Sanitize()
	s3 := o.PieceOptions.Sanitize()
	s4 := o.MiscOptions.Sanitize()
	return s1 || s2 || s3 || s4
}

func (o *Setup) SanitizeSingle() bool {
	s1 := o.GameOptions.SanitizeSingle()
	s2 := o.FieldOptions.Sanitize()
	s3 := o.PieceOptions.Sanitize()
	s4 := o.MiscOptions.Sanitize()
	return s1 || s2 || s3 || s4
}

func (o *Setup) SanitizeMulti() bool {
	s1 := o.GameOptions.SanitizeMulti()
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
	FieldCount       byte     `json:"field_count"`
	TeamSize         byte     `json:"team_size"`
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
	w.Write8(o.FieldCount-1, 3) // 1..MaxFieldCount
	w.Write8(o.TeamSize-1, 2)   // 1..MaxTeamSize
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
	if o.GameType != GameTypeFallingPolyominoes {
		o.GameType = GameTypeFallingPolyominoes
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

func (o *GameOptions) SanitizeSingle() bool {
	sanitized := o.Sanitize()
	if o.FieldCount != 1 {
		o.FieldCount = 1
		sanitized = true
	}
	if o.TeamSize != 1 {
		o.TeamSize = 1
		sanitized = true
	}
	if o.PieceCollision {
		o.PieceCollision = false
		sanitized = true
	}
	if o.PlayerZones {
		o.PlayerZones = false
		sanitized = true
	}
	if o.SamePiecesForAll {
		o.SamePiecesForAll = false
		sanitized = true
	}
	return sanitized
}

func (o *GameOptions) SanitizeMulti() bool {
	sanitized := o.Sanitize()
	if o.FieldCount == 1 && o.TeamSize == 1 {
		o.FieldCount = 2
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

		return fmt.Sprintf("%d%s ", players, smileyFull) +
			strings.Join(slices.Repeat([]string{string('0' + o.TeamSize)}, int(o.FieldCount)), "vs")
	}())

	opts := strings.Builder{}

	if o.TeamSize > 1 && o.PlayerZones {
		opts.WriteByte('Z')
	}

	if o.TeamSize > 1 && !o.PlayerZones && o.PieceCollision {
		opts.WriteByte('!')
	}

	if o.FieldCount*o.TeamSize > 1 && !o.SamePiecesForAll {
		opts.WriteRune('≠')
	}

	if opts.Len() > 0 {
		sb.WriteByte(' ')
		sb.WriteString(opts.String())
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
	PieceType PieceType
	PieceSize byte
	BagSize   byte
}

func (o *PieceOptions) Write(w *bitdata.Writer) {
	w.Write8(byte(o.PieceType), 4)
	w.Write8(o.PieceSize-3, 2) // 3, 4, 5
	w.Write8(o.BagSize-1, 3)   // 1..8
}

func (o *PieceOptions) Read(r *bitdata.ReaderError) {
	o.PieceType = PieceType(r.Read8(4))
	o.PieceSize = r.Read8(2) + 3
	o.BagSize = r.Read8(3) + 1
}

func (o *PieceOptions) Sanitize() bool {
	var sanitized bool
	if o.PieceType > PieceTypeVMirroringPolyominoes {
		o.PieceType = PieceTypeRotatingPolyominoes
		sanitized = true
	}

	if o.PieceSize == 0 {
		o.PieceSize = PieceSizeDefault
		sanitized = true
	} else if o.PieceSize < PieceSizeMin {
		o.PieceSize = PieceSizeMin
		sanitized = true
	} else if o.PieceSize > PieceSizeMax {
		o.PieceSize = PieceSizeMax
		sanitized = true
	}

	if o.BagSize == 0 || o.BagSize > BagSizeMax {
		o.BagSize = BagSizeDefault
	}

	return sanitized
}

func (o *PieceOptions) String() string {
	sb := strings.Builder{}
	switch o.PieceType {
	case PieceTypeRotatingPolyominoes:
		sb.WriteByte('R')
	case PieceTypeHMirroringPolyominoes:
		sb.WriteString("HM")
	case PieceTypeVMirroringPolyominoes:
		sb.WriteString("VM")
	}

	sb.WriteString(strconv.Itoa(int(o.PieceSize)))

	return sb.String()
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

type PlayerConfig piece.Config

func DefaultPlayerConfig() PlayerConfig {
	return PlayerConfig{
		RotationDirectionCW: false,
		SlideDisabled:       false,
		WallKick:            piece.WallKickDefault,
	}
}

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
