package lib

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grafana/loki/pkg/logcli/client"
)

func checkServer(c client.Client, p Params) func() tea.Msg {
	return func() tea.Msg {

		resp, err := c.QueryRange(
			p.Query(),
			p.Limit,
			time.Now().Add(-p.Since),
			time.Now().Add(-p.Until),
			p.Direction,
			0,
			0,
			true,
		)

		if err != nil {
			return err
		}
		return resp
	}
}
