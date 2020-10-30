package lib

import (
	"strings"
	"testing"

	"github.com/muesli/reflow/ansi"
	"github.com/stretchr/testify/require"
)

func TestIntent(t *testing.T) {
	h := DefaultHelp()
	h.width = 100
	h.height = 10

	expected := `q (or C-c) -> quit     n -> next pane         p -> previous pane     h (or ←) -> move left 
                                                                                           
                                                                                           
                                                                                           
j (or ↓) -> move down  k (or ↑) -> move up    l (or →) -> move right
                                                                    
                                                                    
                                                                    
`

	d := h.Drawer()
	for _, line := range strings.Split(expected, "\n") {
		out := quickRender(ansi.PrintableRuneWidth(line), d)
		require.Equal(t, line, out)
		d.Advance()
	}

}
