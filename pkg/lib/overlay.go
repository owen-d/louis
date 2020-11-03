package lib

import (
	"strings"

	"github.com/muesli/termenv"
)

var profile = termenv.ColorProfile() // keep a process wide reference to the color profile.

type Overlay struct {
	xs []Index
}

func (o Overlay) Drawer() Drawer {
	x := overlayDraw(o)
	x.xs = make([]Index, len(o.xs))
	_ = copy(x.xs, o.xs)
	return &x

}

func (o *Overlay) Add(s string, c termenv.Color) {
	if len(s) < 1 {
		return
	}

	delimited := strings.Split(s, "\n")
	for i, x := range delimited {
		o.xs = append(o.xs, Index{
			color: c,
			xs:    []rune(x),
		})
		if ln := len(delimited); ln > 1 && i < ln-1 {
			o.xs = append(o.xs, Index{newline: true})
		}
	}
}

// Ensure the drawable type is hidden to prevent accidentally writing to an overlay while it's being drawn.
type overlayDraw Overlay

func (o *overlayDraw) Done() bool {
	return len(o.xs) == 0
}

// Draw uses a line wrapping strategy and helps implement Drawer.
func (o *overlayDraw) Draw(n int) (results Renderables) {
	if o.Done() {
		return nil
	}

	var ln int
	var newStart int

	for i := range o.xs {
		// We want to mutate the indices, so grab a ref.
		x := &o.xs[i]

		diff := n - ln
		if diff < 1 {
			break
		}

		if x.newline {
			newStart = i + 1
			break
		}

		usable := min(diff, len(x.xs))

		results = append(results, &Index{
			color: x.color,
			xs:    x.xs[0:usable],
		})
		ln = ln + usable

		x.xs = x.xs[usable:]
		// if the current index is exhausted, remove it.
		if len(x.xs) == 0 {
			newStart = i + 1
		}
	}

	o.xs = o.xs[newStart:]
	return results
}

func (o *overlayDraw) Advance() {
	if o.Done() {
		return
	}

	if o.xs[0].newline {
		o.xs = o.xs[1:]
	}

}

type Index struct {
	color   termenv.Color
	xs      []rune
	newline bool // if it holds a newline delimiter instead of a string
}

func (i *Index) Style() Style {
	return Style{
		Foreground: i.color,
	}
}

func (i *Index) String() string { return string(i.xs) }

func (i *Index) Width() int { return len(i.xs) }
