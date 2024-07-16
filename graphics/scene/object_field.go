// Copyright (c) 2024 by Marko Gaćeša

package scene

import (
	"context"
	"gamatet/game/core"
	"gamatet/game/field"
	"gamatet/graphics/render"
	"gamatet/graphics/textcanvas"
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type Field struct {
	fIdx  int
	field *field.Field

	model *mgl32.Mat4

	renderer *render.FieldRender

	renderRequester core.RenderRequester
	renderInfo      *field.RenderInfo
	renderInfoCh    chan *field.RenderInfo

	text *textcanvas.TextCanvas
}

var _ Object = (*Field)(nil)

func NewField(fIdx int, f *field.Field, model *mgl32.Mat4, resources *render.Resources) *Field {
	return &Field{
		fIdx:         fIdx,
		field:        f,
		model:        model,
		renderer:     render.NewFieldRenderer(resources),
		renderInfoCh: make(chan *field.RenderInfo, 1),
	}
}

func (f *Field) StartPrepare(ctx context.Context, now time.Time) {
	f.renderRequester.RenderRequest(ctx, f.fIdx, now, f.renderInfoCh)
}

func (f *Field) EndPrepare(ctx context.Context) {
	select {
	case <-ctx.Done():
	case f.renderInfo = <-f.renderInfoCh:
	}
}

func (f *Field) Render(r *render.Renderer) {
	if f.renderInfo == nil {
		return
	}
	f.renderer.Render(r, f.model, f.renderInfo)
	field.ReturnRenderInfo(f.renderInfo)
	f.renderInfo = nil
}
