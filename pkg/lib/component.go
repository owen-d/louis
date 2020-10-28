package lib

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Updater interface {
	Update(tea.Msg) tea.Cmd
}

// Component is an intersection interface defining what I consider a "component". The idea
// is to use this interface to build higher order components.
type Component interface {
	Drawable
	Updater
}

// NoopUpdater allows defining a component from a Drawable.
type NoopUpdater struct {
	Drawable
}

func (NoopUpdater) Update(_ tea.Msg) tea.Cmd { return nil }

// NoopDrawable allows defining a component from an Updater.
type NoopDrawable struct {
	Updater
}

func (NoopDrawable) Drawer() Drawer { return NoopDrawer{} }

type NoopDrawer struct{}

func (NoopDrawer) Draw(n int) string { return "" }
func (NoopDrawer) Advance()          {}

// The empty component, affectionately named after the empty Monoid from Haskell.
// This technically satisfies the Component interface, but does nothing
type Mempty struct{}

func (Mempty) Update(_ tea.Msg) tea.Cmd { return nil }

func (Mempty) View() string { return "" }
