// Copyright (c) 2024 by Marko Gaćeša

package menu

type Command struct {
	base
	text string
	fn   func()
}

func NewCommand(text, description string, fn func()) *Command {
	return &Command{
		base: base{description: description},
		text: text,
		fn:   fn,
	}
}

func (c *Command) Text() string {
	return c.text
}

func (c *Command) Increase() {}
func (c *Command) Decrease() {}

func (c *Command) Input(r rune) {
	if fn := c.fn; fn != nil {
		c.fn = nil
		fn()
	}
}
