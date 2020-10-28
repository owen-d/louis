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

	require.Equal(t, "", o.Draw(2))

	o.Add("abc\ndef", nil)
	o.Add("ghi", termenv.ANSIYellow)
	o.Add("jkl\nmno\np", nil)

	require.Equal(t, "abc", o.Draw(4))
	require.Equal(t, "def"+termenv.String("gh").Foreground(profile.Convert(termenv.ANSIYellow)).String(), o.Draw(5))
	require.Equal(
		t,
		termenv.String("i").Foreground(profile.Convert(termenv.ANSIYellow)).String()+"jkl",
		o.Draw(5),
	)
	o.Advance()
	require.Equal(t, "mn", o.Draw(2))
	o.Advance() // skip remaining "o" on this line
	require.Equal(t, "p", o.Draw(2))
}

type s struct{ xs []int }

func (in *s) copy() *s {
	tmp := *in
	return &tmp
}

func TestOverlayMultiDraw(t *testing.T) {
	var o Overlay
	o.Add(`ok`, nil)
	require.Equal(t, "ok", o.Drawer().Draw(2))
	require.Equal(t, "ok", o.Drawer().Draw(2))
	require.Equal(t, "ok", o.Draw(2))
	require.Equal(t, "", o.Draw(2))
}
