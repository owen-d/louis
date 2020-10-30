package lib

import (
	"testing"

	"github.com/muesli/reflow/ansi"
	"github.com/stretchr/testify/require"
)

func Test_ExactWidthDrawer(t *testing.T) {
	var o Overlay
	o.Add("foo\nbar", nil)

	for _, tc := range []string{
		"fo", "foo", "foo ",
	} {
		var r Renderer
		ln := ansi.PrintableRuneWidth(tc)
		out := r.Render(ExactWidthDrawer{o.Drawer()}.Draw(ln))
		require.Equal(t, tc, out)
	}

	// test advance
	var r Renderer
	d := ExactWidthDrawer{o.Drawer()}
	require.Equal(t, "foo ", r.Render(d.Draw(4)))
	d.Advance()
	require.Equal(t, "bar ", r.Render(d.Draw(4)))

}
