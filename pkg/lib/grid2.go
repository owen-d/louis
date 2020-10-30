package lib

type Grid struct {
	hSpacing,
	vSpacing,
	width,
	unitWidth,
	height,
	unitHeight int
	rows []CrossMerge

	// drawing internals
	curLine int
}

func NewGrid(vSpacing, hSpacing, height, width, minColumnWidth, maxColumns int, views ...Drawable) *Grid {
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

	grid := &Grid{
		vSpacing:   vSpacing,
		hSpacing:   hSpacing,
		height:     height,
		unitHeight: unitHeight,
		width:      width,
		unitWidth:  unitWidth,
		rows:       make([]CrossMerge, rows),
	}

	for i, v := range views {
		row := i / cols
		vp := Viewport{
			ModelHeight: unitHeight,
			ModelWidth:  unitWidth,
			Component:   NoopUpdater{v},
		}.Drawer()
		grid.rows[row] = append(grid.rows[row], vp)
	}

	for i, x := range grid.rows {
		grid.rows[i] = x.Intersperse(MergableSep{
			Sep: " ",
		})
	}

	return grid
}

func (g *Grid) Height() int { return g.height }

func (g *Grid) Draw(n int) (res Renderables) {
	block := g.unitHeight + g.vSpacing
	// determine if we're drawing a vertical spacer or not
	row := g.curLine / block
	isSpacing := (g.curLine % block) > g.unitHeight

	if isSpacing || row >= len(g.rows) {
		d := ExactWidthDrawer{Content("").Drawer()}
		return d.Draw(min(g.width, n))
	}

	return g.rows[row].Draw(n)
}

func (g *Grid) Advance() {
	block := g.unitHeight + g.vSpacing
	// determine if we're drawing a vertical spacer or not
	row := g.curLine / block

	// advance the current row being drawn
	if row < len(g.rows) {
		g.rows[row].Advance()
	}
	g.curLine++
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
