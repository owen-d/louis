package lib

type VMergable interface {
	Drawer
	Height() int
}

type VMerge struct {
	xs  []VMergable
	cur int
}

func NewVMerge(xs ...VMergable) *VMerge {
	return &VMerge{
		xs: xs,
	}
}

func (v *VMerge) Done() bool {
	return len(v.xs) == 0
}

func (v *VMerge) Draw(n int) (res Renderables) {
	if v.Done() {
		return nil
	}

	return v.xs[0].Draw(n)
}

func (v *VMerge) Advance() {
	if len(v.xs) == 0 {
		return
	}

	x := v.xs[0]
	x.Advance()
	v.cur++
	if v.cur >= x.Height() {
		// advance internal drawer & reset current line within said drawer.
		v.xs = v.xs[1:]
		v.cur = 0
	}
}

type heightDrawer struct {
	height int
	Drawer
}

func NewHeightDrawer(n int, d Drawer) VMergable {
	return &heightDrawer{
		height: n,
		Drawer: d,
	}
}

func (d *heightDrawer) Height() int { return d.height }
