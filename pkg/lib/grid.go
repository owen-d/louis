package lib

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
)

type Grid struct {
	hSpacing,
	vSpacing,
	width,
	unitWidth,
	height,
	unitHeight int
	rows [][]Viewport
}

func NewGrid(vSpacing, hSpacing, height, width, minColumnWidth, maxColumns int, views ...Drawable) Grid {

	ln := len(views)
	cols := max(1, min(min(maxColumns, ln), width/minColumnWidth))
	rows := ln / cols
	if ln%cols > 0 {
		rows++
	}

	var unitHeight int
	for height > 0 && unitHeight < 1 {
		unitHeight = (height - (rows-1)*vSpacing) / rows
		if unitHeight > 0 {
			break
		}

		// elimnate rows that would overflow
		rows--
		views = views[:rows*cols]

	}
	unitWidth := (width - (cols-1)*hSpacing) / cols

	grid := Grid{
		vSpacing:   vSpacing,
		hSpacing:   hSpacing,
		height:     height,
		unitHeight: unitHeight,
		width:      width,
		unitWidth:  unitWidth,
		rows:       make([][]Viewport, rows),
	}

	for i, v := range views {
		row := i / cols
		vp := Viewport{
			ModelHeight: unitHeight,
			ModelWidth:  unitWidth,
			Component:   NoopUpdater{v},
		}
		grid.rows[row] = append(grid.rows[row], vp)
	}

	return grid
}

func (g Grid) View() string {
	mergers := make([]CrossMerge, 0, len(g.rows))

	for _, row := range g.rows {
		merger := make(CrossMerge, 0, len(row))
		for _, v := range row {
			merger = append(merger, v.Drawer())
		}
		mergers = append(mergers, merger.Intersperse(MergableSep{
			Sep: " ",
		}))
	}

	var output string

	if len(g.rows) < 2 {
		output = mergers[0].View()
	} else {
		// build vertical separator
		var sb strings.Builder
		unit := strings.Repeat(" ", g.width) + "\n"
		for i := 0; i < g.vSpacing; i++ {
			sb.WriteString(unit)
		}

		if g.vSpacing == 0 {
			// If there is zero vSpacing specified, just write a newline.
			sb.WriteString("\n")
		}

		vSep := sb.String()

		var result strings.Builder
		for i, m := range mergers {
			result.WriteString(m.View())
			if i < len(mergers)-1 {
				result.WriteString(vSep)
			}

		}
		output = result.String()
	}

	// Finally, bound it to a viewport to ensure desired size.
	var v viewport.Model
	v.Height = g.height
	v.Width = g.width
	v.SetContent(output)
	return viewport.View(v)

}

func min(x, y int) int {
	res := x
	if y < x {
		res = y
	}
	return res
}

func max(x, y int) int {
	res := x
	if y > x {
		res = y
	}
	return res
}
