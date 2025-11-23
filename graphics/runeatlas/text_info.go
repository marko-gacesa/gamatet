// Copyright (c) 2024 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package runeatlas

type RectUV [4]float32

const (
	U0 = 0
	V0 = 1
	U1 = 2
	V1 = 3
)

func (uv RectUV) OffsetUV() [2]float32 {
	return [2]float32{uv[0], uv[1]}
}

func (uv RectUV) ScaleUV() [2]float32 {
	dimU := uv[2] - uv[0]
	dimV := uv[3] - uv[1]
	return [2]float32{dimU, dimV}
}

func (uv RectUV) WidthToHeight() float32 {
	dimU := uv[2] - uv[0]
	dimV := uv[3] - uv[1]
	return dimU / dimV
}

func (uv RectUV) Width() float32  { return uv[2] - uv[0] }
func (uv RectUV) Height() float32 { return uv[3] - uv[1] }

func (uv RectUV) UV0() (float32, float32) { return uv[0], uv[1] }
func (uv RectUV) UV1() (float32, float32) { return uv[2], uv[3] }

func (uv RectUV) U0() float32 { return uv[U0] }
func (uv RectUV) V0() float32 { return uv[V0] }
func (uv RectUV) U1() float32 { return uv[U1] }
func (uv RectUV) V1() float32 { return uv[V1] }
