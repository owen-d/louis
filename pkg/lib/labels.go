package lib

import (
	"log"

	"github.com/grafana/loki/pkg/loghttp"
	"github.com/muesli/gamut"
	"github.com/muesli/termenv"
)

type ColoredLabels struct {
	colorMap map[string]termenv.Color
}

func NewColoredLabels(streams loghttp.Streams) *ColoredLabels {
	result := &ColoredLabels{
		colorMap: make(map[string]termenv.Color),
	}

	for _, stream := range streams {
		for name, _ := range stream.Labels {
			result.colorMap[name] = nil
		}
	}

	colors, err := gamut.Generate(len(result.colorMap), gamut.PastelGenerator{})
	if err != nil {
		log.Println("error generating colors:", err)
		return result
	}

	var i int
	for name, _ := range result.colorMap {
		result.colorMap[name] = profile.FromColor(colors[i])
		i++
	}

	return result
}
