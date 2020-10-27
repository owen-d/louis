package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/owen-d/louis/pkg/lib"
)

func main() {

	// Set PAGER_LOG to a path to log to a file. For example:
	//
	//     export PAGER_LOG=debug.log
	//
	// This becomes handy when debugging stuff since you can't debug to stdout
	// because the UI is occupying it!
	path := os.Getenv("PAGER_LOG")
	if path != "" {
		f, err := tea.LogToFile(path, "pager")
		if err != nil {
			fmt.Printf("Could not open file %s: %v", path, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	var m lib.Model

	p := tea.NewProgram(&m)

	// Use the full size of the terminal in its "alternate screen buffer"
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	// We also turn on mouse support so we can track the mouse wheel
	p.EnableMouseCellMotion()
	defer p.DisableMouseCellMotion()

	if err := p.Start(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
}
