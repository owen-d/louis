package lib

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
)

type Direction int

const (
	Left Direction = iota
	Down
	Up
	Right
)

func GetDirection(msg tea.Msg) (Direction, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			return Left, nil
		case "down", "j":
			return Down, nil
		case "up", "k":
			return Up, nil
		case "right", "l":
			return Right, nil
		}
		return 0, errors.Errorf("not directional key: %s", msg.String())
	}
	return 0, errors.Errorf("not KeyMsg: %T", msg)
}
