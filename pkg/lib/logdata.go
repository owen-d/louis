package lib

import (
	"sort"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/grafana/loki/pkg/loghttp"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	"github.com/prometheus/prometheus/pkg/labels"
)

type streamEntry struct {
	fp uint64
	loghttp.Entry
}

type streamEntries []streamEntry

// Len is the number of elements in the collection.
func (xs streamEntries) Len() int {
	return len(xs)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (xs streamEntries) Less(i, j int) bool {
	return !xs[i].Entry.Timestamp.After(xs[j].Entry.Timestamp)
}

// Swap swaps the elements with indexes i and j.
func (xs streamEntries) Swap(i, j int) { xs[i], xs[j] = xs[j], xs[i] }

type LogDataWidths struct {
	LabelsWidth, LogsWidth int
}

type LogData struct {
	coloredLabels map[uint64]*Overlay
	entries       streamEntries
	sep           MergableSep
	LogDataWidths
}

func NewLogData(streams loghttp.Streams, widths LogDataWidths, sep MergableSep) *LogData {
	result := LogData{
		coloredLabels: make(map[uint64]*Overlay),
		sep:           sep,
		LogDataWidths: widths,
	}

	result.SetStreams(streams)
	return &result
}

func (d *LogData) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case *loghttp.QueryResponse:
		d.SetStreams(msg.Data.Result.(loghttp.Streams))
	case *LogDataWidths:
		d.LogDataWidths = *msg
	}
	return nil
}

func (d *LogData) SetStreams(streams loghttp.Streams) {
	colorMap := make(map[string]termenv.Color)
	for _, stream := range streams {
		for name := range stream.Labels {
			colorMap[name] = nil
		}
	}

	colors, err := gamut.Generate(len(colorMap), gamut.PastelGenerator{})
	if err != nil {
	}

	var i int
	for name, _ := range colorMap {
		colorMap[name] = profile.FromColor(colors[i])
		i++
	}

	for _, stream := range streams {
		var o Overlay
		ls := labels.FromMap(stream.Labels.Map())
		sort.Sort(ls)

		o.Add("{", nil)
		for i, l := range ls {
			if i > 0 {
				o.Add(", ", nil)
			}
			o.Add(l.Name, colorMap[l.Name])
			quoted := strconv.Quote(l.Value)
			o.Add("="+quoted, nil)
		}
		o.Add("}\n", nil)
		fp := ls.Hash()
		d.coloredLabels[fp] = &o

		for _, entry := range stream.Entries {
			d.entries = append(d.entries, streamEntry{
				fp:    fp,
				Entry: entry,
			})
		}

	}

	sort.Sort(d.entries)
}

func (d *LogData) SetWidths(labelsWidth, logsWidth int) {
	d.LabelsWidth = labelsWidth
	d.LogsWidth = logsWidth
}

func (d *LogData) Len() int {
	return d.entries.Len()
}

func (d *LogData) Offset(n int) Drawer {
	drawer := &logDataDrawer{
		LogData: *d,
	}
	drawer.entries = drawer.entries[n:]
	return drawer
}

type logDataDrawer struct {
	i     int
	cache *struct {
		labels Drawer
		line   Drawer
	}
	LogData
}

func (d *logDataDrawer) Done() bool {
	return d.i >= len(d.entries)
}

func (d *logDataDrawer) Width() int { return d.LabelsWidth + d.LogsWidth + d.sep.Width() }

func (d *logDataDrawer) Advance() {
	if d.cache != nil && (!d.cache.labels.Done() || !d.cache.line.Done()) {
		d.cache.labels.Advance()
		d.cache.line.Advance()
		return
	}
	d.i++
	d.cache = nil
}

func (d *logDataDrawer) Draw(n int) (results Renderables) {
	if d.Done() {
		return nil
	}

	if d.cache == nil {

		entry := d.entries[d.i]
		labelsOverlay := d.coloredLabels[entry.fp]
		var lineOverlay Overlay
		lineOverlay.Add(entry.Line, nil)

		d.cache = &struct {
			labels Drawer
			line   Drawer
		}{
			labels: labelsOverlay.Drawer(),
			line:   lineOverlay.Drawer(),
		}
	}

	return CrossMerge{
		NewWidthDrawer(d.LabelsWidth, d.cache.labels),
		d.sep,
		NewWidthDrawer(d.LogsWidth, d.cache.line),
	}.Draw(n)
}
