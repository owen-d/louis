package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExactWidthDrawer(t *testing.T) {
	var o Overlay
	o.Add("foo\nbar", nil)

	require.Equal(t, "fo", ExactWidthDrawer{o.Drawer()}.Draw(2))
	require.Equal(t, "foo", ExactWidthDrawer{o.Drawer()}.Draw(3))
	require.Equal(t, "foo ", ExactWidthDrawer{o.Drawer()}.Draw(4))

	d := ExactWidthDrawer{o.Drawer()}
	require.Equal(t, "foo ", d.Draw(4))
	d.Advance()
	require.Equal(t, "bar ", d.Draw(4))

}
