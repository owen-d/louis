package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntent(t *testing.T) {
	h := DefaultHelp()
	h.Width = 100
	h.Height = 10

	require.Equal(t, "x", h.View())

}
