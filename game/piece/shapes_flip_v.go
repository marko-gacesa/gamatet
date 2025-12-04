// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package piece

type shapeRectV shapeRect

var (
	shapesFlipVTinyminoes = []shapeRectV{
		{width: 1, height: 1, size: 1, data: 1},
		{width: 2, height: 1, size: 2, data: 3},
		{width: 1, height: 2, size: 2, data: 3},
		{width: 3, height: 1, size: 3, data: 7},
		{width: 1, height: 3, size: 3, data: 7},
		{width: 2, height: 2, size: 3, data: 13},
		{width: 2, height: 2, size: 3, data: 14},
	}

	shapesFlipVTetrominoes = []shapeRectV{
		{width: 2, height: 2, size: 4, data: 15},
		{width: 4, height: 1, size: 4, data: 15},
		{width: 1, height: 4, size: 4, data: 15},
		{width: 3, height: 2, size: 4, data: 58},
		{width: 2, height: 3, size: 4, data: 29},
		{width: 2, height: 3, size: 4, data: 46},
		{width: 2, height: 3, size: 4, data: 53},
		{width: 2, height: 3, size: 4, data: 58},
		{width: 3, height: 2, size: 4, data: 57},
		{width: 3, height: 2, size: 4, data: 60},
		{width: 2, height: 3, size: 4, data: 45},
		{width: 3, height: 2, size: 4, data: 30},
	}

	shapesFlipVPentominoes = []shapeRectV{
		{width: 5, height: 1, size: 5, data: 31},
		{width: 1, height: 5, size: 5, data: 31},
		{width: 3, height: 3, size: 5, data: 185},
		{width: 3, height: 3, size: 5, data: 188},
		{width: 3, height: 3, size: 5, data: 410},
		{width: 3, height: 3, size: 5, data: 242},
		{width: 4, height: 2, size: 5, data: 241},
		{width: 4, height: 2, size: 5, data: 248},
		{width: 2, height: 4, size: 5, data: 213},
		{width: 2, height: 4, size: 5, data: 234},
		{width: 3, height: 2, size: 5, data: 59},
		{width: 3, height: 2, size: 5, data: 62},
		{width: 2, height: 3, size: 5, data: 61},
		{width: 2, height: 3, size: 5, data: 62},
		{width: 4, height: 2, size: 5, data: 227},
		{width: 4, height: 2, size: 5, data: 124},
		{width: 2, height: 4, size: 5, data: 181},
		{width: 2, height: 4, size: 5, data: 122},
		{width: 3, height: 3, size: 5, data: 466},
		{width: 3, height: 3, size: 5, data: 121},
		{width: 3, height: 3, size: 5, data: 316},
		{width: 3, height: 2, size: 5, data: 61},
		{width: 2, height: 3, size: 5, data: 55},
		{width: 2, height: 3, size: 5, data: 59},
		{width: 3, height: 3, size: 5, data: 457},
		{width: 3, height: 3, size: 5, data: 484},
		{width: 3, height: 3, size: 5, data: 409},
		{width: 3, height: 3, size: 5, data: 244},
		{width: 3, height: 3, size: 5, data: 186},
		{width: 4, height: 2, size: 5, data: 242},
		{width: 4, height: 2, size: 5, data: 244},
		{width: 2, height: 4, size: 5, data: 117},
		{width: 2, height: 4, size: 5, data: 186},
		{width: 3, height: 3, size: 5, data: 313},
		{width: 3, height: 3, size: 5, data: 214},
	}
)
