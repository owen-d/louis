package lib

import (
	"strings"

	"github.com/muesli/termenv"
)

type Renderable interface {
	Style() Style
	Width() int
	String() string
}

type Style struct {
	Foreground termenv.Color
}

type Renderer struct{}

func (r Renderer) Render(xs []Renderable) string {
	var b strings.Builder

	for _, x := range xs {
		s := termenv.String(x.String())
		if f := x.Style().Foreground; f != nil {
			s = s.Foreground(profile.Convert(f))
		}
		b.WriteString(s.String())
	}
	return b.String()

}

func quickRender(n int, d Drawer) string {
	var r Renderer
	return r.Render(d.Draw(n))
}

type Renderables []Renderable

func (xs Renderables) Width() (res int) {
	for _, x := range xs {
		res += x.Width()
	}
	return res
}
