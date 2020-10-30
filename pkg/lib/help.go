package lib

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/ansi"
)

type Intent struct {
	Primary string
	Aliases []string
	Msg     string
}

func (i Intent) String() string {
	strs := []string{i.Primary}
	if len(i.Aliases) > 0 {
		strs = append(strs, fmt.Sprintf("(or %s)", strings.Join(i.Aliases, ",")))
	}
	strs = append(strs, []string{"->", i.Msg}...)
	return strings.Join(strs, " ")
}

func (i Intent) Drawer() Drawer {
	var o Overlay
	o.Add(i.String(), nil)
	return o.Drawer()
}

type HelpPane struct {
	height, width int
	intents       []Intent
}

func (h HelpPane) Height() int {
	return h.height
}

func (h HelpPane) Drawer() *Grid {
	var minColumWidth int
	vs := make([]Drawable, 0, len(h.intents))
	for _, intent := range h.intents {
		vs = append(vs, intent)
		if w := ansi.PrintableRuneWidth(intent.String()); minColumWidth < w {
			minColumWidth = w
		}
	}

	return NewGrid(
		0,
		4,
		h.height-1,
		h.width,
		minColumWidth,
		7,
		vs...,
	)
}

func DefaultHelp() HelpPane {
	return HelpPane{
		height: 0,
		width:  0,
		intents: []Intent{
			{
				Primary: "q",
				Aliases: []string{"C-c"},
				Msg:     "quit",
			},
			{
				Primary: "n",
				Msg:     "next pane",
			},
			{
				Primary: "p",
				Msg:     "previous pane",
			},
			{
				Primary: "h",
				Aliases: []string{"←"},
				Msg:     "move left",
			},
			{
				Primary: "j",
				Aliases: []string{"↓"},
				Msg:     "move down",
			},
			{
				Primary: "k",
				Aliases: []string{"↑"},
				Msg:     "move up",
			},
			{
				Primary: "l",
				Aliases: []string{"→"},
				Msg:     "move right",
			},
		},
	}
}
