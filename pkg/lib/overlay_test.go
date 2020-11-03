package lib

import (
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/require"
)

func TestPad(t *testing.T) {
	s := termenv.String("gh").Foreground(profile.Convert(termenv.ANSIYellow)).String()
	require.Equal(t, s+"  ", RPad(s, 4))
}

func TestOverlayDraw(t *testing.T) {
	var o Overlay

	require.Equal(t, "", quickRender(2, o.Drawer()))

	o.Add("abc\ndef", nil)
	o.Add("ghi", termenv.ANSIYellow)
	o.Add("jkl\nmno\np", nil)

	d := o.Drawer()
	require.Equal(t, "abc", quickRender(4, d))
	require.Equal(t, "def"+termenv.String("gh").Foreground(profile.Convert(termenv.ANSIYellow)).String(), quickRender(5, d))
	require.Equal(
		t,
		termenv.String("i").Foreground(profile.Convert(termenv.ANSIYellow)).String()+"jkl",
		quickRender(5, d),
	)
	d.Advance()
	require.Equal(t, "mn", quickRender(2, d))
	d.Advance() // don't skip because we word wrap
	require.Equal(t, "o", quickRender(1, d))
	d.Advance() // advance should skip the newline
	require.Equal(t, "p", quickRender(2, d))
}

func TestNewlineEnding(t *testing.T) {
	var o Overlay
	o.Add("abcde\nf", nil)
	d := o.Drawer()

	require.Equal(t, "abc", quickRender(3, d))
	require.Equal(t, "de", quickRender(3, d))
	require.Equal(t, "f", quickRender(3, d))

}

func TestOverlayMultiDraw(t *testing.T) {
	var o Overlay
	o.Add(`ok`, nil)
	require.Equal(t, "ok", quickRender(2, o.Drawer()))
	require.Equal(t, "ok", quickRender(2, o.Drawer()))

	d := o.Drawer()
	require.Equal(t, "ok", quickRender(2, d))
	require.Equal(t, "", quickRender(2, d))
}

func TestNoWrap(t *testing.T) {
	var o Overlay
	o.Strategy(NoWrap)

	o.Add("abc\ndef", nil)
	d := o.Drawer()
	d.Advance()
	require.Equal(t, "def", quickRender(3, d))
}
