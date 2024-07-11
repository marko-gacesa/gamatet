// Copyright (c) 2024 by Marko Gaćeša

package material

import "github.com/go-gl/gl/v4.1-core/gl"

var _ Material = (*Text)(nil)

func NewText(tex uint32) *Text {
	p, err := newProgramBlock(defaultVertexShader, defaultFragmentShader, tex)
	if err != nil {
		panic("failed to make text material: " + err.Error())
	}

	return &Text{
		programBlock: *p,
	}
}

// Text is a material that is used for drawing text.
type Text struct {
	programBlock
}

func (t *Text) Use() {
	t.programBlock.Use()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func (t *Text) Reset() {
	gl.Disable(gl.BLEND)
}
