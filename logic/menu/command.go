// Copyright (c) 2024 by Marko Gaćeša

package menu

type Command struct {
	base
	text   string
	parent *Menu
	fn     func(*Menu, *Command)
}

func NewCommand(text, description string, fn func(*Menu, *Command)) *Command {
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
	if r != '\n' {
		return
	}
	if fn := c.fn; fn != nil {
		fn(c.parent, c)
	}
}

func (c *Command) setParent(menu *Menu) {
	c.parent = menu
}

func (c *Command) SetText(text string) {
	c.text = text
}

func (c *Command) SetDescription(desc string) {
	c.description = desc
}

func (c *Command) ClearFunction() {
	c.fn = nil
}
