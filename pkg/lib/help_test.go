package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntent(t *testing.T) {
	h := DefaultHelp()
	h.Width = 100
	h.Height = 10

	expected := `────────────────────────────────────────────────────────────────────────────────────────────────────
q (or C-c) -> quit     n -> next pane         p -> previous pane     h (or ←) -> move left 
                                                                                           
                                                                                           
                                                                                           
j (or ↓) -> move down  k (or ↑) -> move up    l (or →) -> move right
                                                                    
                                                                    
                                                                    
`

	require.Equal(t, expected, h.View())

}
