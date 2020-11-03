package lib

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Scrollable interface {
	Len() int
	Offset(int) Drawer
	Updater
}

type LazyScroll struct {
	Offset     int
	scrollable Scrollable
}

func NewLazyScroll(s Scrollable) *LazyScroll {
	return &LazyScroll{
		scrollable: s,
	}
}

func (s *LazyScroll) Update(msg tea.Msg) tea.Cmd {
	if dir, ok := GetDirection(msg); ok {
		switch dir {
		case Down:
			s.Offset = clamp(s.Offset+1, 0, s.scrollable.Len())
		case Up:
			s.Offset = clamp(s.Offset-1, 0, s.scrollable.Len())
		}
	}

	return s.scrollable.Update(msg)
}

func (s *LazyScroll) Drawer() Drawer {
	return s.scrollable.Offset(s.Offset)
}

func clamp(v, floor, ceil int) int {
	return min(ceil, max(floor, v))
}
