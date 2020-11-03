package lib

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Direction int

const (
	Left Direction = iota
	Down
	Up
	Right
)

func GetDirection(msg tea.Msg) (Direction, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			return Left, true
		case "down", "j":
			return Down, true
		case "up", "k":
			return Up, true
		case "right", "l":
			return Right, true
		}
	}
	return 0, false
}
