package lib

import (
	"github.com/muesli/termenv"
)

var profile = termenv.ColorProfile()

type Drawable interface {
	Drawer() Drawer
}

type Drawer interface {
	Draw(n int) string // draw n width
	Advance()          // newline
}

type Content string

func (s Content) Drawer() Drawer {
	var o Overlay
	o.Add(string(s), nil)
	return &o
}
