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

func (c CrossMerge) Lines() (lines []string) {
	var maxLines int
	for _, x := range c {
		if height := x.Height(); height > maxLines {
			maxLines = height
		}
	}

	for i := 0; i < maxLines; i++ {
		var line string
		for _, x := range c {
			line += x.Draw(x.Width())
			x.Advance()
		}

		lines = append(lines, line)
	}

	return lines

}

func (c CrossMerge) View() string {
	return strings.Join(c.Lines(), "\n")
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
