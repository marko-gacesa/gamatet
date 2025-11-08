// Copyright (c) 2025 by Marko Gaćeša

package setup

var (
	ColorRGB = [MaxPlayers][3]float32{
		{0.3, 0.3, 1.0}, // blue
		{1.0, 0.3, 0.3}, // red
		{0.0, 1.0, 0.2}, // green
		{1.0, 1.0, 0.2}, // yellow
		{0.9, 0.1, 0.9}, // magenta
		{0.1, 0.8, 0.9}, // cyan
		{1.0, 0.6, 0.2}, // orange
		{0.6, 1.0, 0.2}, // lime
	}
	ColorBackRGB = [MaxPlayers][3]float32{}
)

func init() {
	for i := range MaxPlayers {
		for j := range 3 {
			ColorBackRGB[i][j] = ColorRGB[i][j] * 0.15
		}
	}
}
