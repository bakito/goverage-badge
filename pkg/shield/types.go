package shield

import (
	"fmt"
	"sort"
)

var (
	defaultColor    = "brightgreen"
	defaultSeverity = SeverityMap{
		"red":         0,
		"orange":      20,
		"yellow":      30,
		"yellowgreen": 40,
		"green":       55,
		"brighgreen":  70,
	}
)

// Badge https://shields.io/endpoint
type Badge struct {
	// Required. Always the number 1.
	SchemaVersion int `json:"schemaVersion" yaml:"schemaVersion"`
	// Required. The left text, or the empty string to omit the left side of the badge. This can be overridden by the query string.
	Label string `json:"label" yaml:"label"`
	// Required. Can't be empty. The right text.
	Message string `json:"message" yaml:"message"`
	// Default: lightgrey. The right color. Supports the eight named colors above, as well as hex, rgb, rgba, hsl, hsla and css named colors. This can be overridden by the query string.
	Color string `json:"color" yaml:"color"`
	// Default: grey. The left color. This can be overridden by the query string.
	LabelColor string `json:"labelColor,omitempty" yaml:"labelColor,omitempty"`
	// Default: false. true to treat this as an error badge. This prevents the user from overriding the color. In the future it may affect cache behavior.
	IsError bool `json:"isError,omitempty" yaml:"isError,omitempty"`
	// Default: none. One of the named logos supported by Shields or simple-icons. Can be overridden by the query string.
	NamedLogo string `json:"namedLogo,omitempty" yaml:"namedLogo,omitempty"`
	// Default: none. An SVG string containing a custom logo.
	LogoSvg string `json:"logoSvg,omitempty" yaml:"logoSvg,omitempty"`
	// Default: none. Same meaning as the query string. Can be overridden by the query string.
	LogoColor string `json:"logoColor,omitempty" yaml:"logoColor,omitempty"`
	// Default: none. Same meaning as the query string. Can be overridden by the query string.
	LogoWidth string `json:"logoWidth,omitempty" yaml:"logoWidth,omitempty"`
	// Default: none. Same meaning as the query string. Can be overridden by the query string.
	LogoPosition string `json:"logoPosition,omitempty" yaml:"logoPosition,omitempty"`
	// Default: flat. The default template to use. Can be overridden by the query string.
	Syle string `json:"syle,omitempty" yaml:"syle,omitempty"`
	// Default: 300, min 300. Set the HTTP cache lifetime in seconds, which should be respected by the Shields' CDN and downstream users. Values below 300 will be ignored. This lets you tune performance and traffic vs. responsiveness. The value you specify can be overridden by the user via the query string, but only to a longer value.
	CacheSecondsDefault int `json:"cacheSecondsDefault,omitempty" yaml:"cacheSecondsDefault,omitempty"`
}

// SeverityMap map color to severity
type SeverityMap map[string]int

type severity struct {
	from  int
	color string
}

func (sm SeverityMap) color(coverage float64) string {

	var s []severity
	for k, v := range sm {
		s = append(s, severity{color: k, from: v})
	}

	sort.Slice(s, func(i, j int) bool {
		if s[i].from != s[j].from {
			return s[i].from > s[j].from
		}

		// get priority from defaultSeverity if from is equal
		si, iOK := defaultSeverity[s[i].color]
		sj, jOK := defaultSeverity[s[j].color]

		if iOK && jOK {
			return si < sj
		}

		return iOK
	})

	for _, r := range s {
		if coverage >= float64(r.from) {
			return r.color
		}
	}

	return ""
}

// Setup the badge
func (b *Badge) Setup(label string, coverage float64, color string, severity *SeverityMap) {
	b.SchemaVersion = 1
	b.Label = label
	if b.Label == "" {
		b.Label = "coverage"
	}
	b.Message = fmt.Sprintf("%.f%%", coverage)

	if color != "" {
		b.Color = color
	} else {
		if severity != nil && len(*severity) != 0 {
			b.Color = severity.color(coverage)
		} else {
			b.Color = defaultSeverity.color(coverage)
		}
	}

	if b.Color == "" {
		b.Color = defaultColor
	}
}
