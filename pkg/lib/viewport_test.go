package lib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViewport_Draw(t *testing.T) {
	data := []string{
		"Hello there",
		"I'm a test to see if",
		"we correctly draw.",
		"Bye!",
		"should not print",
	}

	v := Viewport{
		ModelWidth:  20,
		ModelHeight: len(data) - 1,
		YPosition:   0,
		Component: NoopUpdater{
			Content(strings.Join(data, "\n")),
		},
	}

	d := v.Drawer()
	var out strings.Builder
	for i := 0; i < d.ModelHeight; i++ {
		out.WriteString(d.Draw(v.ModelWidth))
		d.Advance()
		out.WriteString("\n")
	}

	var expected strings.Builder
	for i := 0; i < v.ModelHeight; i++ {
		x := data[i]
		expected.WriteString(ExactWidth(x, v.ModelWidth) + "\n")
	}

	require.Equal(t, expected.String(), out.String())
}
