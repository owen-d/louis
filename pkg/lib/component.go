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

func (NoopDrawer) Draw(n int) Renderables { return nil }
func (NoopDrawer) Advance()               {}
func (NoopDrawer) Done() bool             { return true }

type DrawableFunc func() Drawer

func (fn DrawableFunc) Drawer() Drawer { return fn() }
