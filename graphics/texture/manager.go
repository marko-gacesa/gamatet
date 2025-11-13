// Copyright (c) 2020-2024 by Marko Gaćeša

package texture

import (
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Manager struct {
	textures []uint32
}

func Init() *Manager {
	var maxTextures int32
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &maxTextures)

	textures := make([]uint32, maxTextures)

	return &Manager{textures: textures}
}

func (m *Manager) MaxTextures() int {
	return len(m.textures)
}

// Bind binds the provided image as OpenGL texture.
// The idea is to first bind the texture: `texID := manager.Bind(img)`.
// And then to activate it: `gl.ActiveTexture(texID)`.
func (m *Manager) Bind(img image.Image) uint32 {
	var idx int
	for idx = 0; idx < len(m.textures); idx++ {
		if m.textures[idx] == 0 {
			break
		}
	}
	if idx == len(m.textures) {
		return 0
	}

	var texID = gl.TEXTURE0 + uint32(idx)

	texture := genTex()
	m.textures[idx] = texture

	loadData(texID, texture, img)

	return texID
}

func (m *Manager) ReBind(texID uint32, img image.Image) uint32 {
	idx := texID - gl.TEXTURE0
	texture := m.textures[idx]

	loadData(texID, texture, img)

	return texID
}

func (m *Manager) Delete(texID uint32) {
	idx := texID - gl.TEXTURE0
	texture := m.textures[idx]

	gl.ActiveTexture(texID)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.DeleteTextures(1, &texture)

	m.textures[idx] = 0
}

func genTex() uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	return texture
}

func loadData(texID, texture uint32, img image.Image) {
	switch v := img.(type) {
	case *image.RGBA:
		loadData4Channel(v.Rect.Size(), v.Pix, texID, texture)
	case *image.Gray:
		loadData1Channel(v.Rect.Size(), v.Pix, texID, texture)
	case *image.Alpha:
		loadData1Channel(v.Rect.Size(), v.Pix, texID, texture)
	default:
		panic("unsupported color model")
	}
}

func loadData4Channel(size image.Point, data []byte, texID, texture uint32) uint32 {
	gl.ActiveTexture(texID)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(size.X),
		int32(size.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(data))

	return texture
}

func loadData1Channel(size image.Point, data []byte, texID, texture uint32) uint32 {
	gl.ActiveTexture(texID)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.R8,
		int32(size.X),
		int32(size.Y),
		0,
		gl.RED,
		gl.UNSIGNED_BYTE,
		gl.Ptr(data))

	return texture
}
