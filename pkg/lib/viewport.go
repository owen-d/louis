package lib

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/ansi"
)

type Viewport struct {
	ModelWidth, ModelHeight, YPosition int
	Component
}

func (v *Viewport) Width() int { return v.ModelWidth }

func (v *Viewport) Height() int { return v.ModelHeight }

func (v Viewport) Drawer() *ViewportDrawer {
	return &ViewportDrawer{Viewport: v, Drawer: v.Component.Drawer()}
}

type ViewportDrawer struct {
	Viewport // embed for height/width methods
	Drawer
}

type viewports struct {
	totals               tea.WindowSizeMsg
	ready                bool
	focusPane            Pane
	separator            MergableSep
	params, labels, logs Viewport
	help                 HelpPane
}

func (v *viewports) focused() *Viewport {
	switch v.focusPane {
	case LabelsPane:
		return &v.labels
	case LogsPane:
		return &v.logs
	default:
		return &v.params
	}
}

func (v *viewports) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.Size(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			v.focusPane = v.focusPane.Next()
			v.Size(v.totals)
		case "p":
			v.focusPane = v.focusPane.Prev()
			v.Size(v.totals)
		}
	}

	cmd := v.focused().Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (v *viewports) selected() (main *Viewport, secondaries []*Viewport) {
	switch v.focusPane {
	case LabelsPane:
		return &v.labels, []*Viewport{&v.params, &v.logs}
	case LogsPane:
		return &v.logs, []*Viewport{&v.params, &v.labels}
	// ParamsPane is the default
	default:
		return &v.params, []*Viewport{&v.labels, &v.logs}
	}
}

// Size sets pane sizes (primary & secondaries) based on the golden ratio.
func (v *viewports) Size(msg tea.WindowSizeMsg) {
	v.totals = msg
	width := msg.Width - v.separator.Width()*2
	if !v.ready {
		v.ready = true
	}

	v.help.Height = 4
	v.help.Width = v.totals.Width

	withoutHeaders := msg.Height - 3
	withoutHelp := withoutHeaders - v.help.Height

	height := withoutHelp

	v.params.ModelHeight = height
	v.labels.ModelHeight = height
	v.logs.ModelHeight = height

	primary := int(float64(width) / GoldenRatio)
	secondary := (width - primary) / 2
	main, secondaries := v.selected()
	main.ModelWidth = primary
	for _, s := range secondaries {
		s.ModelWidth = secondary
	}
}

func (v *viewports) header() string {
	pane := v.focusPane
	width := v.totals.Width
	var start int

	switch pane {
	case LabelsPane:
		start = v.params.Width() + v.separator.Width()
	case LogsPane:
		start = v.params.Width()*2 + v.separator.Width()*2 // all non-primary panes have the same size
	}

	headerTopFrame := "╭─────────────╮"
	headerBotFrame := "╰─────────────╯"
	headerTop := ExactWidth(LPad(headerTopFrame, start+ansi.PrintableRuneWidth(headerTopFrame)), width)
	headerBot := ExactWidth(LPad(headerBotFrame, start+ansi.PrintableRuneWidth(headerBotFrame)), width)

	lConnector := "│"
	if start > 0 {
		lConnector = "┤"
	}
	headerMid := lConnector + CenterTo(pane.String(), ansi.PrintableRuneWidth(headerTopFrame)-2) + "├"
	headerMid = LPadWith(headerMid, '─', start+ansi.PrintableRuneWidth(headerMid))
	headerMid = RPadWith(headerMid, '─', width)

	return strings.Join([]string{headerTop, headerMid, headerBot}, "\n")
}

func (v *viewports) View() string {
	if !v.ready {
		return "\n  Initializing..."
	}

	merger := CrossMerge{
		v.params.Drawer(),
		v.separator,
		v.labels.Drawer(),
		v.separator,
		v.logs.Drawer(),
	}

	return strings.Join(
		[]string{
			v.header(),
			merger.View(),
			v.help.View(),
		},
		"\n",
	)
}
