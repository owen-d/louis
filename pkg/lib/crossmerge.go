package lib

import (
	"github.com/muesli/reflow/ansi"
)

/*
Idea here is to merge views which combine horizontally (on the same line)

------   ------   ------------
|    | + |    | = |          |
|    | + |    | = |          |
------   ------   ------------

*/

type CrossMergable interface {
	Drawer
	Width() int
}

type MergableSep struct {
	Sep string
}

func (s MergableSep) Draw(n int) Renderables {
	return []Renderable{
		&Index{
			xs: []rune(Truncate(s.Sep, n)),
		},
	}
}

func (s MergableSep) Advance()   {}
func (s MergableSep) Done() bool { return true }

func (s MergableSep) Width() int {
	return ansi.PrintableRuneWidth(s.Sep)
}

func (s MergableSep) Height() int {
	return 1
}

type CrossMerge []CrossMergable

func (c CrossMerge) Width() (res int) {
	for _, x := range c {
		res += x.Width()
	}
	return
}

func (c CrossMerge) Done() bool {
	for _, x := range c {
		if !x.Done() {
			return false
		}
	}
	return true
}

func (c CrossMerge) Draw(n int) (res Renderables) {
	rem := n
	var i int
	for ; rem > 0 && i < len(c); i++ {
		d := ExactWidthDrawer{c[i]}
		// draw up to the drawer's width or rem, whichever is smaller
		additions := d.Draw(min(rem, c[i].Width()))
		rem -= additions.Width()
		res = append(res, additions...)
	}
	return res
}

func (c CrossMerge) Advance() {
	for _, x := range c {
		x.Advance()
	}
}

func (c CrossMerge) Intersperse(x CrossMergable) CrossMerge {
	ln := 2*len(c) - 1
	res := make(CrossMerge, 0, ln)

	for i, y := range c {
		res = append(res, y)
		if i < (len(c) - 1) {
			res = append(res, x)
		}
	}

	return res

}

type widthDrawer struct {
	width int
	Drawer
}

func NewWidthDrawer(n int, d Drawer) CrossMergable {
	return &widthDrawer{
		width:  n,
		Drawer: d,
	}
}

func (d *widthDrawer) Width() int { return d.width }
