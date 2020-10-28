package lib

import (
	"strings"

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
	Height() int
	Width() int
}

type MergableSep struct {
	Sep string
}

func (s MergableSep) Draw(n int) string {
	return Truncate(s.Sep, n)
}

func (s MergableSep) Advance() {}

func (s MergableSep) Width() int {
	return ansi.PrintableRuneWidth(s.Sep)
}

func (s MergableSep) Height() int {
	return 0
}

type CrossMerge []CrossMergable

func (c CrossMerge) Width() (res int) {
	for _, x := range c {
		res += x.Width()
	}
	return

}

func (c CrossMerge) View() string {
	var maxLines int
	for _, x := range c {
		if height := x.Height(); height > maxLines {
			maxLines = height
		}
	}

	var sb strings.Builder

	for i := 0; i < maxLines; i++ {
		for _, x := range c {
			drawer := ExactWidthDrawer{x}
			addition := drawer.Draw(x.Width())
			sb.WriteString(addition)
			drawer.Advance()
		}

		if i < maxLines-1 {
			sb.WriteString("\n")
		}

	}

	return sb.String()

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
