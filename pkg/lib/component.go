package lib

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Updater interface {
	Update(tea.Msg) tea.Cmd
}

type Component interface {
	Drawable
	Updater
}

type NoopUpdater struct {
	Drawable
}

func (NoopUpdater) Update(_ tea.Msg) tea.Cmd { return nil }

type NoopDrawable struct {
	Updater
}

func (NoopDrawable) Drawer() Drawer { return NoopDrawer{} }

type NoopDrawer struct{}

func (NoopDrawer) Draw(n int) string { return strings.Repeat(" ", n) }
func (NoopDrawer) Advance()          {}

type Mempty struct{}

func (Mempty) Init() tea.Cmd { return nil }

func (Mempty) View() string { return "" }
