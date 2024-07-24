// Copyright (c) 2020-2023 by Marko Gaćeša

package texture

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"image"
	"image/draw"
	"os"
)

var Instance *Manager

type Manager struct {
	textures []uint32
}

func Init() *Manager {
	var maxTextures int32
	gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &maxTextures)

	textures := make([]uint32, maxTextures)

	return &Manager{textures: textures}
}

// Bind binds the provided image as OpenGL texture.
// The idea is to first bind the texture: `texName := manager.Bind(img)`.
// And then to activate it: `gl.ActiveTexture(texName)`.
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

	switch v := img.(type) {
	case *image.RGBA:
		m.textures[idx] = NewTexture2DFromRGB(v, texID)
	case *image.Gray:
		m.textures[idx] = NewTexture2DFromGray(v, texID)
	default:
		panic("unsupported color model")
	}

	return texID
}

func (m *Manager) Bind3D(values []float32, dim int) uint32 {
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

	texture := NewTexture3DFromSlice(values, dim, texID)
	m.textures[idx] = texture

	return texID
}

func (m *Manager) Delete(texID uint32) {
	idx := texID - gl.TEXTURE0
	texture := m.textures[idx]
	if texture == 0 {
		return
	}
	gl.DeleteTextures(1, &texture)
	m.textures[idx] = 0
}

func (m *Manager) DeleteAll() {
	for idx, texture := range m.textures {
		if texture == 0 {
			continue
		}
		gl.DeleteTextures(1, &texture)
		m.textures[idx] = 0
	}
}

func LoadFile(file string) (image.Image, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported texture %q stride=%d", file, rgba.Stride)
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Src)

	return rgba, nil
}

func NewTexture2DFromRGB(rgba *image.RGBA, textureNumber uint32) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(textureNumber) // gl.TEXTURE0
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

func NewTexture2DFromGray(rgba *image.Gray, textureNumber uint32) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(textureNumber) // gl.TEXTURE0
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RED,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RED,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

func NewTexture3DFromSlice(values []float32, dim int, textureNumber uint32) uint32 {
	if len(values) != dim*dim*dim {
		panic("length of values and dimension don't match")
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(textureNumber) // gl.TEXTURE0
	gl.BindTexture(gl.TEXTURE_3D, texture)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_3D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage3D(
		gl.TEXTURE_3D,
		0,
		gl.RED,
		int32(dim), int32(dim), int32(dim),
		0,
		gl.RED,
		gl.FLOAT,
		gl.Ptr(values))

	return texture
}
