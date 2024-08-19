// Copyright (c) 2024 by Marko Gaćeša

package menu

type Static struct {
	base
	text string
}

func NewStatic(text, description string) *Static {
	return &Static{
		base: base{description: description},
		text: text,
	}
}

func (c *Static) Text() string { return c.text }

func (c *Static) Increase() {}
func (c *Static) Decrease() {}
