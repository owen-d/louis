package lib

import (
	"strings"

	"github.com/muesli/termenv"
)

type index struct {
	color   termenv.Color
	xs      []rune
	newline bool // if it holds a newline delimiter instead of a string
}

type Overlay struct {
	xs []*index
}

func (o *Overlay) IsEmpty() bool {
	return len(o.xs) == 0
}

// impl Drawable as well for convenience
func (o *Overlay) Drawer() Drawer {
	x := *o
	return &x
}

func (o *Overlay) Add(s string, c termenv.Color) {
	delimited := strings.Split(s, "\n")
	for i, x := range delimited {
		o.xs = append(o.xs, &index{
			color: c,
			xs:    []rune(x),
		})
		if ln := len(delimited); ln > 1 && i < ln-1 {
			o.xs = append(o.xs, &index{newline: true})
		}
	}
}

func (o *Overlay) Draw(n int) string {
	if o.IsEmpty() {
		return RPad(" ", n)
	}

	var ln int
	var b strings.Builder
	var newStart int

	for i, x := range o.xs {
		diff := n - ln
		if x.newline {
			newStart = i + 1
			break
		}

		if diff < 1 {
			break
		}

		usable := min(diff, len(x.xs))
		sub := termenv.String(string(x.xs[0:usable])).Foreground(profile.Convert(x.color))
		b.WriteString(sub.String())
		ln = ln + usable

		x.xs = x.xs[usable:]
		// if the current index is exhausted, remove it.
		if len(x.xs) == 0 {
			newStart = i + 1
		}
	}

	o.xs = o.xs[newStart:]

	// add any additional padding if needed
	out := RPad(b.String(), n)
	return out
}

func (o *Overlay) Advance() {
	if o.IsEmpty() {
		return
	}

	if o.xs[0].newline {
		o.xs = o.xs[1:]
	}

}
