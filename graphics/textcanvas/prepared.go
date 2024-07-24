// Copyright (c) 2024 by Marko Gaćeša

package textcanvas

import "image/color"

type Prepared struct {
	canvas  TextCanvas
	textMap map[int][4]float32
	nextID  int
}

func NewPrepared(size int) *Prepared {
	return &Prepared{
		canvas:  *NewTextCanvas(size),
		textMap: make(map[int][4]float32),
		nextID:  1,
	}
}

func (p *Prepared) Prepare(text string, face Face, color color.Color, lrPad bool) int {
	pos := p.canvas.TextUV(text, face, color, lrPad)
	id := p.nextID
	p.textMap[id] = pos
	p.nextID++
	return id
}

func (p *Prepared) Text(id int) [4]float32 {
	return p.textMap[id]
}

func (p *Prepared) Save(path string) error {
	return p.canvas.Save(path)
}
