package lib

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grafana/loki/pkg/loghttp"
	"github.com/muesli/reflow/ansi"
)

type Pane int

const (
	ParamsPane Pane = iota
	LabelsPane
	LogsPane

	MinPane = ParamsPane
	MaxPane = LogsPane

	GoldenRatio = 1.618
)

func (p Pane) String() string {
	switch p {
	case LabelsPane:
		return "labels"
	case LogsPane:
		return "logs"
	default:
		return "params"
	}
}

func (p Pane) Next() Pane {
	n := p + 1
	if n > MaxPane {
		n = MinPane
	}
	return n
}

func (p Pane) Prev() Pane {
	n := p - 1
	if n < MinPane {
		n = MaxPane
	}
	return n
}

type panes struct {
	totals               tea.WindowSizeMsg
	focusPane            Pane
	separator            MergableSep
	params, labels, logs Viewport
	help                 HelpPane

	// data
	streams loghttp.Streams
}

func (v *panes) Height() int {
	return v.totals.Height
}
func (v *panes) Width() int {
	return v.totals.Width
}

func (v *panes) focused() *Viewport {
	switch v.focusPane {
	case LabelsPane:
		return &v.labels
	case LogsPane:
		return &v.logs
	default:
		return &v.params
	}
}

func (v *panes) Update(msg tea.Msg) tea.Cmd {
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

	case *loghttp.QueryResponse:
		v.streams = msg.Data.Result.(loghttp.Streams)
	}

	cmd := v.focused().Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (v *panes) selected() (main *Viewport, secondaries []*Viewport) {
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
func (v *panes) Size(msg tea.WindowSizeMsg) {
	v.totals = msg
	width := msg.Width - v.separator.Width()*2

	v.help.height = 3
	v.help.width = v.totals.Width

	// height = msgHeight - header height - footer margin - footer height
	height := msg.Height - 3 - 1 - v.help.height

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

func (v *panes) header() string {
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

func (v *panes) Drawer() Drawer {
	return ExactWidthDrawer{
		NewVMerge(
			NewHeightDrawer(3, Content(v.header()).Drawer()),
			NewHeightDrawer(
				v.params.Height(),
				CrossMerge{
					v.params.Drawer(),
					NewLogData(v.streams, v.labels.Width(), v.logs.Width(), v.separator).Drawer(),
				}.Intersperse(v.separator),
			),
			NewHeightDrawer(
				1,
				Content(RPadWith("", '─', v.totals.Width)).Drawer(),
			),
			v.help.Drawer(),
		),
	}
}
