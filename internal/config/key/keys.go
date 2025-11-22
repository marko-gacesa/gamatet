// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package key

import (
	"fmt"
	"slices"

	"github.com/marko-gacesa/gamatet/internal/i18n"
)

type Input struct {
	Left     Key `json:"left"`
	Right    Key `json:"right"`
	Activate Key `json:"activate"`
	Boost    Key `json:"boost"`
	Drop     Key `json:"drop"`
}

func (in *Input) Sanitize(idx int) {
	keys := [5]Key{in.Left, in.Right, in.Activate, in.Boost, in.Drop}
	slices.Sort(keys[:])
	uq := slices.Compact(keys[:])
	if len(uq) != len(keys) {
		*in = DefaultInput[idx%len(DefaultInput)]
	}
}

func (in *Input) String() string {
	return fmt.Sprintf("%s=%s %s=%s %s=%s %s=%s %s=%s",
		i18n.T(i18n.KeyConfigPlayerKeyLeft), Map[in.Left],
		i18n.T(i18n.KeyConfigPlayerKeyRight), Map[in.Right],
		i18n.T(i18n.KeyConfigPlayerKeyActivate), Map[in.Activate],
		i18n.T(i18n.KeyConfigPlayerKeyBoost), Map[in.Boost],
		i18n.T(i18n.KeyConfigPlayerKeyDrop), Map[in.Drop],
	)
}

var (
	InputWASD   = Input{Left: A, Right: D, Activate: W, Boost: S, Drop: Tab}
	InputHJKL   = Input{Left: H, Right: L, Activate: K, Boost: J, Drop: Space}
	InputArrows = Input{Left: Left, Right: Right, Activate: Up, Boost: Down, Drop: Space}
	InputNumPad = Input{Left: KP4, Right: KP6, Activate: KP8, Boost: KP2, Drop: KP0}

	DefaultInput = []Input{InputArrows, InputWASD, InputHJKL, InputNumPad}
)

type Key byte

func (k *Key) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return fmt.Errorf("invalid key: %s", s)
	}
	for kk, v := range Map {
		if v == s {
			*k = kk
		}
	}
	*k = Unknown
	return nil
}

func (k Key) MarshalJSON() ([]byte, error) {
	s, ok := Map[k]
	if !ok {
		s = StrUnknown
	}

	return fmt.Appendf(nil, `"%s"`, s), nil
}

const (
	Unknown Key = iota
	Space
	Apostrophe
	Comma
	Minus
	Period
	Slash
	N0
	N1
	N2
	N3
	N4
	N5
	N6
	N7
	N8
	N9
	Semicolon
	Equal
	A
	B
	C
	D
	E
	F
	G
	H
	I
	J
	K
	L
	M
	N
	O
	P
	Q
	R
	S
	T
	U
	V
	W
	X
	Y
	Z
	LeftBracket
	Backslash
	RightBracket
	GraveAccent
	World1
	World2
	Escape
	Enter
	Tab
	Backspace
	Insert
	Delete
	Right
	Left
	Down
	Up
	PageUp
	PageDown
	Home
	End
	CapsLock
	ScrollLock
	NumLock
	PrintScreen
	Pause
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	F13
	F14
	F15
	F16
	F17
	F18
	F19
	F20
	F21
	F22
	F23
	F24
	F25
	KP0
	KP1
	KP2
	KP3
	KP4
	KP5
	KP6
	KP7
	KP8
	KP9
	KPDecimal
	KPDivide
	KPMultiply
	KPSubtract
	KPAdd
	KPEnter
	KPEqual
	LeftShift
	LeftControl
	LeftAlt
	LeftSuper
	RightShift
	RightControl
	RightAlt
	RightSuper
	Menu
)

const (
	StrUnknown      = "???"
	StrSpace        = "SPACE"
	StrApostrophe   = "'"
	StrComma        = ","
	StrMinus        = "-"
	StrPeriod       = "."
	StrSlash        = "/"
	StrN0           = "0"
	StrN1           = "1"
	StrN2           = "2"
	StrN3           = "3"
	StrN4           = "4"
	StrN5           = "5"
	StrN6           = "6"
	StrN7           = "7"
	StrN8           = "8"
	StrN9           = "9"
	StrSemicolon    = ";"
	StrEqual        = "="
	StrA            = "A"
	StrB            = "B"
	StrC            = "C"
	StrD            = "D"
	StrE            = "E"
	StrF            = "F"
	StrG            = "G"
	StrH            = "H"
	StrI            = "I"
	StrJ            = "J"
	StrK            = "K"
	StrL            = "L"
	StrM            = "M"
	StrN            = "N"
	StrO            = "O"
	StrP            = "P"
	StrQ            = "Q"
	StrR            = "R"
	StrS            = "S"
	StrT            = "T"
	StrU            = "U"
	StrV            = "V"
	StrW            = "W"
	StrX            = "X"
	StrY            = "Y"
	StrZ            = "Z"
	StrLeftBracket  = "["
	StrBackslash    = "\\"
	StrRightBracket = "]"
	StrGraveAccent  = "`"
	StrWorld1       = "W1"
	StrWorld2       = "W2"
	StrEscape       = "ESC"
	StrEnter        = "ENTER"
	StrTab          = "TAB"
	StrBackspace    = "BSP"
	StrInsert       = "INS"
	StrDelete       = "DEL"
	StrRight        = "RIGHT"
	StrLeft         = "LEFT"
	StrDown         = "DOWN"
	StrUp           = "UP"
	StrPageUp       = "PAGE_UP"
	StrPageDown     = "PAGE_DOWN"
	StrHome         = "HOME"
	StrEnd          = "END"
	StrCapsLock     = "CAPS"
	StrScrollLock   = "SCROLL"
	StrNumLock      = "NUM"
	StrPrintScreen  = "PRINT"
	StrPause        = "PAUSE"
	StrF1           = "F1"
	StrF2           = "F2"
	StrF3           = "F3"
	StrF4           = "F4"
	StrF5           = "F5"
	StrF6           = "F6"
	StrF7           = "F7"
	StrF8           = "F8"
	StrF9           = "F9"
	StrF10          = "F10"
	StrF11          = "F11"
	StrF12          = "F12"
	StrF13          = "F13"
	StrF14          = "F14"
	StrF15          = "F15"
	StrF16          = "F16"
	StrF17          = "F17"
	StrF18          = "F18"
	StrF19          = "F19"
	StrF20          = "F20"
	StrF21          = "F21"
	StrF22          = "F22"
	StrF23          = "F23"
	StrF24          = "F24"
	StrF25          = "F25"
	StrKP0          = "KP0"
	StrKP1          = "KP1"
	StrKP2          = "KP2"
	StrKP3          = "KP3"
	StrKP4          = "KP4"
	StrKP5          = "KP5"
	StrKP6          = "KP6"
	StrKP7          = "KP7"
	StrKP8          = "KP8"
	StrKP9          = "KP9"
	StrKPDecimal    = "KP."
	StrKPDivide     = "KP/"
	StrKPMultiply   = "KP*"
	StrKPSubtract   = "KP-"
	StrKPAdd        = "KP+"
	StrKPEnter      = "KP_ENTER"
	StrKPEqual      = "KP="
	StrLeftShift    = "L-SHIFT"
	StrLeftControl  = "L-CTRL"
	StrLeftAlt      = "L-ALT"
	StrLeftSuper    = "L-SUPER"
	StrRightShift   = "R-SHIFT"
	StrRightControl = "R-CTRL"
	StrRightAlt     = "R-ALT"
	StrRightSuper   = "R-SUPER"
	StrMenu         = "MENU"
)

var Map = map[Key]string{
	Unknown:      StrUnknown,
	Space:        StrSpace,
	Apostrophe:   StrApostrophe,
	Comma:        StrComma,
	Minus:        StrMinus,
	Period:       StrPeriod,
	Slash:        StrSlash,
	N0:           StrN0,
	N1:           StrN1,
	N2:           StrN2,
	N3:           StrN3,
	N4:           StrN4,
	N5:           StrN5,
	N6:           StrN6,
	N7:           StrN7,
	N8:           StrN8,
	N9:           StrN9,
	Semicolon:    StrSemicolon,
	Equal:        StrEqual,
	A:            StrA,
	B:            StrB,
	C:            StrC,
	D:            StrD,
	E:            StrE,
	F:            StrF,
	G:            StrG,
	H:            StrH,
	I:            StrI,
	J:            StrJ,
	K:            StrK,
	L:            StrL,
	M:            StrM,
	N:            StrN,
	O:            StrO,
	P:            StrP,
	Q:            StrQ,
	R:            StrR,
	S:            StrS,
	T:            StrT,
	U:            StrU,
	V:            StrV,
	W:            StrW,
	X:            StrX,
	Y:            StrY,
	Z:            StrZ,
	LeftBracket:  StrLeftBracket,
	Backslash:    StrBackslash,
	RightBracket: StrRightBracket,
	GraveAccent:  StrGraveAccent,
	World1:       StrWorld1,
	World2:       StrWorld2,
	Escape:       StrEscape,
	Enter:        StrEnter,
	Tab:          StrTab,
	Backspace:    StrBackspace,
	Insert:       StrInsert,
	Delete:       StrDelete,
	Right:        StrRight,
	Left:         StrLeft,
	Down:         StrDown,
	Up:           StrUp,
	PageUp:       StrPageUp,
	PageDown:     StrPageDown,
	Home:         StrHome,
	End:          StrEnd,
	CapsLock:     StrCapsLock,
	ScrollLock:   StrScrollLock,
	NumLock:      StrNumLock,
	PrintScreen:  StrPrintScreen,
	Pause:        StrPause,
	F1:           StrF1,
	F2:           StrF2,
	F3:           StrF3,
	F4:           StrF4,
	F5:           StrF5,
	F6:           StrF6,
	F7:           StrF7,
	F8:           StrF8,
	F9:           StrF9,
	F10:          StrF10,
	F11:          StrF11,
	F12:          StrF12,
	F13:          StrF13,
	F14:          StrF14,
	F15:          StrF15,
	F16:          StrF16,
	F17:          StrF17,
	F18:          StrF18,
	F19:          StrF19,
	F20:          StrF20,
	F21:          StrF21,
	F22:          StrF22,
	F23:          StrF23,
	F24:          StrF24,
	F25:          StrF25,
	KP0:          StrKP0,
	KP1:          StrKP1,
	KP2:          StrKP2,
	KP3:          StrKP3,
	KP4:          StrKP4,
	KP5:          StrKP5,
	KP6:          StrKP6,
	KP7:          StrKP7,
	KP8:          StrKP8,
	KP9:          StrKP9,
	KPDecimal:    StrKPDecimal,
	KPDivide:     StrKPDivide,
	KPMultiply:   StrKPMultiply,
	KPSubtract:   StrKPSubtract,
	KPAdd:        StrKPAdd,
	KPEnter:      StrKPEnter,
	KPEqual:      StrKPEqual,
	LeftShift:    StrLeftShift,
	LeftControl:  StrLeftControl,
	LeftAlt:      StrLeftAlt,
	LeftSuper:    StrLeftSuper,
	RightShift:   StrRightShift,
	RightControl: StrRightControl,
	RightAlt:     StrRightAlt,
	RightSuper:   StrRightSuper,
	Menu:         StrMenu,
}
