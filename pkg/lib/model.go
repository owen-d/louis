package lib

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grafana/loki/pkg/logcli/client"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/prometheus/prometheus/pkg/labels"
)

type Model struct {
	client client.Client
	panes  panes
	params Params
}

func (m *Model) Init() tea.Cmd {
	m.panes.separator = MergableSep{
		Sep: " │ ",
	}

	m.params = DefaultParams
	m.panes.params.Component = NoopUpdater{Content(m.params.Content())}
	m.panes.labels.Component = NoopUpdater{Content("")}
	m.panes.logs.Component = NoopUpdater{Content("")}
	m.panes.help = DefaultHelp()

	m.client = &client.DefaultClient{
		Address:  os.Getenv("LOKI_ADDR"),
		Username: os.Getenv("LOKI_USERNAME"),
		Password: os.Getenv("LOKI_PASSWORD"),
	}
	return checkServer(m.client, m.params)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Ctrl+c exits
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			return m, tea.Quit
		}
	}

	if cmd := m.panes.Update(msg); cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var b strings.Builder
	d := m.panes.Drawer()

	for i := 0; i < m.panes.Height(); i++ {
		str := quickRender(m.panes.Width(), d)
		b.WriteString(str)
		d.Advance()
	}

	return b.String()
}

// Hilarious we don't have type for this that's not bound to the ast.
// Mimic 2/3 of a label matcher :)
type Filter struct {
	Type  labels.MatchType
	Match string
}

func (f Filter) String() (res string) {
	switch f.Type {
	case labels.MatchEqual:
		res = "|="
	case labels.MatchRegexp:
		res = "|~"
	case labels.MatchNotEqual:
		res = "!="
	case labels.MatchNotRegexp:
		res = "!~"
	}

	return res + fmt.Sprintf(` "%s"`, f.Match)

}

type Params struct {
	Matchers     []labels.Matcher
	Filters      []Filter
	Since, Until time.Duration
	Direction    logproto.Direction
	Limit        int

	// internals
}

func (p Params) Content() Content {
	var b strings.Builder
	for _, m := range p.Matchers {
		b.WriteString(m.String() + "\n")
	}
	for _, f := range p.Filters {
		b.WriteString(f.String() + "\n")
	}
	b.WriteString(fmt.Sprintf("since: %s\n", p.Since.String()))
	b.WriteString(fmt.Sprintf("until: %s\n", p.Until.String()))
	b.WriteString(fmt.Sprintf("direction: %s\n", p.Direction.String()))
	b.WriteString(fmt.Sprintf("limit: %d\n", p.Limit))

	return Content(b.String())
}

var DefaultParams = Params{
	Matchers: []labels.Matcher{
		{
			Type:  labels.MatchRegexp,
			Name:  "job",
			Value: "tns-demo/.*",
		},
	},
	Filters: []Filter{
		{
			Type:  labels.MatchNotEqual,
			Match: "/metrics",
		},
	},
	Since:     time.Hour,
	Until:     0,
	Direction: logproto.BACKWARD,
	Limit:     200,
}

func (p Params) Query() string {
	mStrs := make([]string, 0, len(p.Matchers))
	for _, m := range p.Matchers {
		mStrs = append(mStrs, m.String())
	}

	var fStr strings.Builder
	for _, f := range p.Filters {
		fStr.WriteString(f.String())
	}
	return fmt.Sprintf("{%s}%s", strings.Join(mStrs, ","), fStr.String())
}
