package lib

// Viewport just embeds a component alongside height/width methods
type Viewport struct {
	ModelWidth, ModelHeight int
	Component
}

func (v *Viewport) Width() int { return v.ModelWidth }

func (v *Viewport) Height() int { return v.ModelHeight }

func (v Viewport) Drawer() *ViewportDrawer {
	return &ViewportDrawer{
		Viewport: v,
		Drawer:   v.Component.Drawer(),
	}
}

type ViewportDrawer struct {
	Viewport // embed for height/width methods
	Drawer
}
