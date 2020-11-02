package lib

import (
	"log"
	"sort"
	"strconv"

	"github.com/grafana/loki/pkg/loghttp"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
	"github.com/prometheus/prometheus/pkg/labels"
)

type LogData struct {
	streams loghttp.Streams
}

func NewLogData(streams loghttp.Streams) *LogData {
	result := &LogData{
		streams: streams,
	}

	return result
}

func (d *LogData) Drawer() Drawer {
	var o Overlay
	colorMap := make(map[string]termenv.Color)
	for _, stream := range d.streams {
		for name := range stream.Labels {
			colorMap[name] = nil
		}
	}

	colors, err := gamut.Generate(len(colorMap), gamut.PastelGenerator{})
	if err != nil {
		log.Println("error generating colors:", err)
		return o.Drawer()
	}

	var i int
	for name, _ := range colorMap {
		colorMap[name] = profile.FromColor(colors[i])
		i++
	}

	for _, stream := range d.streams {
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
	}
	return o.Drawer()
}
