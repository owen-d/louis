package lib

import (
	"github.com/muesli/termenv"
)

var profile = termenv.ColorProfile()

type Drawable interface {
	Drawer() Drawer
}

type Drawer interface {
	Draw(n int) string // draw up to n width, but no longer.
	Advance()          // newline
}

type Content string

func (s Content) Drawer() Drawer {
	var o Overlay
	o.Add(string(s), nil)
	return &o
}

type ExactWidthDrawer struct {
	Drawer
}

func (d ExactWidthDrawer) Draw(n int) string {
	out := d.Drawer.Draw(n)
	return ExactWidth(out, n)
}
