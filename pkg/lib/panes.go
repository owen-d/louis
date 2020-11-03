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
	totals    tea.WindowSizeMsg
	separator MergableSep
	params    Content
	data      *LazyScroll
	help      HelpPane

	focusPane   Pane
	paneHeight  int // how tall each pane should be
	paramsWidth int
}

func (p *panes) Init(sep MergableSep, params Content, data *LogData) {
	p.separator = sep
	p.params = params
	p.data = NewLazyScroll(data)
	p.help = DefaultHelp()
}

func (v *panes) Height() int {
	return v.totals.Height
}
func (v *panes) Width() int {
	return v.totals.Width
}

func (v *panes) focused() Updater {
	switch v.focusPane {
	case LabelsPane, LogsPane:
		return v.data
	default:
		// params currently isn't a component (TODO)
		return nil
	}
}

func (v *panes) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		logDataWidths := v.Size(msg)
		if cmd := v.data.Update(logDataWidths); cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			v.focusPane = v.focusPane.Next()
			if cmd := v.data.Update(v.Size(v.totals)); cmd != nil {
				cmds = append(cmds, cmd)
			}
		case "p":
			v.focusPane = v.focusPane.Prev()
			if cmd := v.data.Update(v.Size(v.totals)); cmd != nil {
				cmds = append(cmds, cmd)
			}
		}

	case *loghttp.QueryResponse:
		// always load data, even when it's not focused.
		if cmd := v.data.Update(msg); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if focused := v.focused(); focused != nil {
		cmd := v.focused().Update(msg)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

// Size sets pane sizes (primary & secondaries) based on the golden ratio.
func (v *panes) Size(msg tea.WindowSizeMsg) *LogDataWidths {
	v.totals = msg
	width := msg.Width - v.separator.Width()*2

	v.help.height = 3
	v.help.width = v.totals.Width

	// height = msgHeight - header height - footer margin - footer height
	height := msg.Height - 3 - 1 - v.help.height

	v.paneHeight = height
	primaryWidth := int(float64(width) / GoldenRatio)
	secondaryWidth := (width - primaryWidth) / 2

	var paramsWidth, labelsWidth, logsWidth int
	switch v.focusPane {
	case LabelsPane:
		labelsWidth = primaryWidth
		paramsWidth, logsWidth = secondaryWidth, secondaryWidth
	case LogsPane:
		logsWidth = primaryWidth
		paramsWidth, labelsWidth = secondaryWidth, secondaryWidth
	default:
		paramsWidth = primaryWidth
		labelsWidth, logsWidth = secondaryWidth, secondaryWidth
	}

	v.paramsWidth = paramsWidth

	return &LogDataWidths{
		LabelsWidth: labelsWidth,
		LogsWidth:   logsWidth,
	}
}

func (v *panes) header() string {
	pane := v.focusPane
	width := v.totals.Width
	var start int

	switch pane {
	case LabelsPane:
		start = v.paramsWidth + v.separator.Width()
	case LogsPane:
		start = v.paramsWidth*2 + v.separator.Width()*2 // all non-primary panes have the same size
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
				v.paneHeight,
				CrossMerge{
					NewWidthDrawer(v.paramsWidth, v.params.Drawer()),
					NewWidthDrawer(v.totals.Width-v.paramsWidth, v.data.Drawer()),
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
