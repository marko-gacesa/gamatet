// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package base

import (
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/marko-gacesa/gamatet/game/action"
	"github.com/marko-gacesa/gamatet/graphics/render"
	"github.com/marko-gacesa/gamatet/graphics/texture"
	"github.com/marko-gacesa/gamatet/logic/screen"
)

func NewBase(
	renderer *render.Renderer,
	tex *texture.Manager,
) Base {
	return Base{
		renderer: renderer,
		tex:      tex,
	}
}

type Base struct {
	renderer *render.Renderer
	tex      *texture.Manager
}

func (s Base) Renderer() *render.Renderer { return s.renderer }
func (s Base) TexMgr() *texture.Manager   { return s.tex }

func (s Base) UpdateViewSize(w, h int) {}

func (s Base) Release() {}

func (s Base) InputKeyPress(key, scancode int) {}
func (s Base) InputChar(char rune)             {}

func (s Base) Prepare(time.Time) {}
func (s Base) Render()           {}

func NewBlockBase(
	renderer *render.Renderer,
	tex *texture.Manager,
	contentW, contentH int,
	usePerspective bool,
) BlockBase {
	return NewBlockBaseWithZ(renderer, tex, contentW, contentH, 3, usePerspective)
}

func NewBlockBaseWithZ(
	renderer *render.Renderer,
	tex *texture.Manager,
	contentW, contentH, contentZ int,
	usePerspective bool,
) BlockBase {
	return BlockBase{
		Base:           NewBase(renderer, tex),
		usePerspective: usePerspective,
		contentW:       contentW,
		contentH:       contentH,
		contentZ:       contentZ,
		viewW:          0,
		viewH:          0,
	}
}

type BlockBase struct {
	Base

	usePerspective bool

	contentW int
	contentH int
	contentZ int
	viewW    int
	viewH    int
}

func (b *BlockBase) InputKeyPress(key, scancode int, act screen.KeyAction) {
	if glfw.Key(key) == glfw.KeyF12 && act == screen.KeyActionPress {
		b.usePerspective = !b.usePerspective
		b.SetCamera()
	}
}

func (b *BlockBase) SetCamera() {
	if b.usePerspective {
		b.Renderer().PerspectiveFull(b.viewW, b.viewH, b.contentW, b.contentH, b.contentZ)
	} else {
		b.Renderer().OrthogonalFull(b.viewW, b.viewH, b.contentW, b.contentH, b.contentZ)
	}
}

func (b *BlockBase) UpdateViewSize(w, h int) {
	b.viewW, b.viewH = w, h
	b.SetCamera()
}

func (b *BlockBase) ViewSize() (w, h int) {
	return b.viewW, b.viewH
}

func (b *BlockBase) TopLeft() mgl32.Mat4 {
	return mgl32.Translate3D(float32(-b.contentW)/2, float32(b.contentH)/2, 0)
}

func (b *BlockBase) TopRight() mgl32.Mat4 {
	return mgl32.Translate3D(float32(b.contentW)/2, float32(b.contentH)/2, 0)
}

func (b *BlockBase) BottomLeft() mgl32.Mat4 {
	return mgl32.Translate3D(float32(-b.contentW)/2, float32(-b.contentH)/2, 0)
}

func (b *BlockBase) BottomRight() mgl32.Mat4 {
	return mgl32.Translate3D(float32(b.contentW)/2, float32(-b.contentH)/2, 0)
}

func SendAction(a action.Action, doneCh <-chan struct{}, cmdCh chan<- []byte) bool {
	if a == action.NoOp || cmdCh == nil {
		return false
	}

	select {
	case <-doneCh:
	case cmdCh <- []byte{byte(a)}:
	}

	return true
}
