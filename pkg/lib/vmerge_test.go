package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_VMerge(t *testing.T) {
	a := NewHeightDrawer(2, Content("a").Drawer())
	b := NewHeightDrawer(1, Content("b").Drawer())

	v := NewVMerge(a, b)

	require.Equal(t, "a", quickRender(1, v))
	v.Advance()
	require.Equal(t, "", quickRender(1, v))
	v.Advance()
	require.Equal(t, "b", quickRender(1, v))
	v.Advance()
	require.Equal(t, "", quickRender(1, v))
}
