// Copyright (c) 2025 by Marko Gaćeša

package key

import (
	"slices"
)

type Input struct {
	Left   Key `json:"left"`
	Right  Key `json:"right"`
	Action Key `json:"action"`
	Drop   Key `json:"drop"`
}

func (in *Input) Sanitize(idx int) {
	keys := [4]Key{in.Left, in.Right, in.Action, in.Drop}
	slices.Sort(keys[:])
	uq := slices.Compact(keys[:])
	if len(uq) != len(keys) {
		*in = DefaultInput[idx%len(DefaultInput)]
	}
}

var (
	InputWASD        = Input{Left: A, Right: D, Action: W, Drop: S}
	InputWASDSpace   = Input{Left: A, Right: D, Action: W, Drop: Space}
	InputHJKL        = Input{Left: H, Right: L, Action: K, Drop: J}
	InputArrows      = Input{Left: Left, Right: Right, Action: Up, Drop: Down}
	InputArrowsSpace = Input{Left: Left, Right: Right, Action: Up, Drop: Space}
	InputNumPad      = Input{Left: KP4, Right: KP6, Action: KP8, Drop: KP5}
	InputNumPad0     = Input{Left: KP4, Right: KP6, Action: KP8, Drop: KP0}

	DefaultInput = []Input{InputArrows, InputWASD, InputHJKL, InputNumPad}
)

type Key byte

func (k *Key) UnmarshalJSON(b []byte) error {
	s := Str(b)
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
		return []byte(StrUnknown), nil
	}

	return []byte(s), nil
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

type Str string

const (
	StrUnknown      Str = "???"
	StrSpace        Str = "SPACE"
	StrApostrophe   Str = "'"
	StrComma        Str = ","
	StrMinus        Str = "-"
	StrPeriod       Str = "."
	StrSlash        Str = "/"
	StrN0           Str = "0"
	StrN1           Str = "1"
	StrN2           Str = "2"
	StrN3           Str = "3"
	StrN4           Str = "4"
	StrN5           Str = "5"
	StrN6           Str = "6"
	StrN7           Str = "7"
	StrN8           Str = "8"
	StrN9           Str = "9"
	StrSemicolon    Str = ";"
	StrEqual        Str = "="
	StrA            Str = "A"
	StrB            Str = "B"
	StrC            Str = "C"
	StrD            Str = "D"
	StrE            Str = "E"
	StrF            Str = "F"
	StrG            Str = "G"
	StrH            Str = "H"
	StrI            Str = "I"
	StrJ            Str = "J"
	StrK            Str = "K"
	StrL            Str = "L"
	StrM            Str = "M"
	StrN            Str = "N"
	StrO            Str = "O"
	StrP            Str = "P"
	StrQ            Str = "Q"
	StrR            Str = "R"
	StrS            Str = "S"
	StrT            Str = "T"
	StrU            Str = "U"
	StrV            Str = "V"
	StrW            Str = "W"
	StrX            Str = "X"
	StrY            Str = "Y"
	StrZ            Str = "Z"
	StrLeftBracket  Str = "["
	StrBackslash    Str = "\\"
	StrRightBracket Str = "]"
	StrGraveAccent  Str = "`"
	StrWorld1       Str = "W1"
	StrWorld2       Str = "W2"
	StrEscape       Str = "ESC"
	StrEnter        Str = "ENTER"
	StrTab          Str = "TAB"
	StrBackspace    Str = "BSP"
	StrInsert       Str = "INS"
	StrDelete       Str = "DEL"
	StrRight        Str = "RIGHT"
	StrLeft         Str = "LEFT"
	StrDown         Str = "DOWN"
	StrUp           Str = "UP"
	StrPageUp       Str = "PAGE_UP"
	StrPageDown     Str = "PAGE_DOWN"
	StrHome         Str = "HOME"
	StrEnd          Str = "END"
	StrCapsLock     Str = "CAPS"
	StrScrollLock   Str = "SCROLL"
	StrNumLock      Str = "NUM"
	StrPrintScreen  Str = "PRINT"
	StrPause        Str = "PAUSE"
	StrF1           Str = "F1"
	StrF2           Str = "F2"
	StrF3           Str = "F3"
	StrF4           Str = "F4"
	StrF5           Str = "F5"
	StrF6           Str = "F6"
	StrF7           Str = "F7"
	StrF8           Str = "F8"
	StrF9           Str = "F9"
	StrF10          Str = "F10"
	StrF11          Str = "F11"
	StrF12          Str = "F12"
	StrF13          Str = "F13"
	StrF14          Str = "F14"
	StrF15          Str = "F15"
	StrF16          Str = "F16"
	StrF17          Str = "F17"
	StrF18          Str = "F18"
	StrF19          Str = "F19"
	StrF20          Str = "F20"
	StrF21          Str = "F21"
	StrF22          Str = "F22"
	StrF23          Str = "F23"
	StrF24          Str = "F24"
	StrF25          Str = "F25"
	StrKP0          Str = "KP0"
	StrKP1          Str = "KP1"
	StrKP2          Str = "KP2"
	StrKP3          Str = "KP3"
	StrKP4          Str = "KP4"
	StrKP5          Str = "KP5"
	StrKP6          Str = "KP6"
	StrKP7          Str = "KP7"
	StrKP8          Str = "KP8"
	StrKP9          Str = "KP9"
	StrKPDecimal    Str = "KP."
	StrKPDivide     Str = "KP/"
	StrKPMultiply   Str = "KP*"
	StrKPSubtract   Str = "KP-"
	StrKPAdd        Str = "KP+"
	StrKPEnter      Str = "KP_ENTER"
	StrKPEqual      Str = "KP="
	StrLeftShift    Str = "L-SHIFT"
	StrLeftControl  Str = "L-CTRL"
	StrLeftAlt      Str = "L-ALT"
	StrLeftSuper    Str = "L-SUPER"
	StrRightShift   Str = "R-SHIFT"
	StrRightControl Str = "R-CTRL"
	StrRightAlt     Str = "R-ALT"
	StrRightSuper   Str = "R-SUPER"
	StrMenu         Str = "MENU"
)

var Map = map[Key]Str{
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
