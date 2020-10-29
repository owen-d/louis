package lib

// A Drawable can "spawn" a Drawer. The idea here is to allow component types,
// typically used as part of a Model, to create these on demand. This allows
// a Component to spawn these on each draw cycle for use in more complex
// drawing strategies (see the CrossMerge, a higher order Drawer)
type Drawable interface {
	Drawer() Drawer
}

// Drawer is a view specific interface that can be used to draw something to screen.
type Drawer interface {
	Draw(n int) []Renderable // Request up to n length drawing, but no longer.
	// Advance a newline. Different strategies could be implemented here,
	// for instance line-wrapping vs clipping.
	Advance()
}

type Content string

func (s Content) Drawer() Drawer {
	var o Overlay
	o.Add(string(s), nil)
	return o.Drawer()
}

type ExactWidthDrawer struct {
	Drawer
}

func (d ExactWidthDrawer) Draw(n int) []Renderable {
	toDraw := d.Drawer.Draw(n)

	remaining := n - Renderables(toDraw).Width()
	if remaining > 0 {
		toDraw = append(toDraw, Content(Padding(remaining)).Drawer().Draw(remaining)...)
	}

	return toDraw

}
