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

func (i Intent) View() string {
	strs := []string{i.Primary}
	if len(i.Aliases) > 0 {
		strs = append(strs, fmt.Sprintf("(or %s)", strings.Join(i.Aliases, ",")))
	}
	strs = append(strs, []string{"->", i.Msg}...)
	return strings.Join(strs, " ")
}

func (i Intent) Drawer() Drawer {
	var o Overlay
	o.Add(i.View(), nil)
	return &o
}

type HelpPane struct {
	Height, Width int
	intents       []Intent
}

func (h HelpPane) View() string {
	var minColumWidth int
	vs := make([]Drawable, 0, len(h.intents))
	for _, intent := range h.intents {
		vs = append(vs, intent)
		if w := ansi.PrintableRuneWidth(intent.View()); minColumWidth < w {
			minColumWidth = w
		}
	}

	topBorder := strings.Repeat("─", h.Width)
	return topBorder + NewGrid(
		0,
		4,
		h.Height-1,
		h.Width,
		minColumWidth,
		7,
		vs...,
	).View()
}

func DefaultHelp() HelpPane {
	return HelpPane{
		Height: 0,
		Width:  0,
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
