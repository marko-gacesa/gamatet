// Copyright (c) 2020-2024 by Marko Gaćeša

package texture

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"image"
	"image/draw"
	"os"
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
		m.textures[idx] = new4Channel(v.Rect.Size(), v.Pix, texID)
	case *image.Gray:
		m.textures[idx] = new1Channel(v.Rect.Size(), v.Pix, texID)
	case *image.Alpha:
		m.textures[idx] = new1Channel(v.Rect.Size(), v.Pix, texID)
	default:
		panic("unsupported color model")
	}

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

	defer imgFile.Close()

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

func new4Channel(size image.Point, data []byte, textureNumber uint32) uint32 {
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
		int32(size.X),
		int32(size.Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(data))

	return texture
}

func new1Channel(size image.Point, data []byte, textureNumber uint32) uint32 {
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
		int32(size.X),
		int32(size.Y),
		0,
		gl.RED,
		gl.UNSIGNED_BYTE,
		gl.Ptr(data))

	return texture
}
