package shield

import (
	"fmt"
	"sort"
)

var (
	severity = map[string]int{
		"red":         60,
		"orange":      50,
		"yellow":      40,
		"yellowgreen": 30,
		"green":       20,
		"brighgreen":  10,
	}
)

type Shield struct {
	SchemaVersion int          `json:"schemaVersion"`
	Label         string       `json:"label"`
	Message       string       `json:"message"`
	Color         string       `json:"color"`
	Ranges        []ColorRange `json:"ranges""`
}

type ColorRange struct {
	From  int    `json:"from""`
	Color string `json:"color""`
}

func (s *Shield) Update(label string, coverage float64) {
	s.SchemaVersion = 1
	s.Label = label
	if s.Label == "" {
		s.Label = "coverage"
	}
	s.Message = fmt.Sprintf("%2.f%%", coverage)

	sort.Slice(s.Ranges, func(i, j int) bool {
		if s.Ranges[i].From != s.Ranges[j].From {
			return s.Ranges[i].From > s.Ranges[j].From
		}
		si, iOK := severity[s.Ranges[i].Color]
		sj, jOK := severity[s.Ranges[j].Color]

		if iOK && jOK {
			return si > sj
		}

		return iOK
	})

	for _, r := range s.Ranges {
		if coverage > float64(r.From) {
			s.Color = r.Color
			return
		}
	}
	if s.Color == "" {
		s.Color = "brightgreen"
	}
}
