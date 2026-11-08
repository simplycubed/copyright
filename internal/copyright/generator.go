package copyright

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/simplycubed/copyright/internal/config"
)

// HeaderGenerator renders a copyright header from the configured template.
type HeaderGenerator struct {
	cfg *config.Config
}

// headerData is the value passed to the copyright template.
// SPDX and SPDXIdentifier are the same spdx_identifier value.
type headerData struct {
	Year           string
	Company        string
	Author         string
	SPDX           string
	SPDXIdentifier string
}

// NewHeaderGenerator returns a generator for cfg.
func NewHeaderGenerator(cfg *config.Config) *HeaderGenerator {
	return &HeaderGenerator{cfg: cfg}
}

// Generate renders the copyright template. Year "auto" and an empty year
// resolve through Config.GetYear.
func (g *HeaderGenerator) Generate() (string, error) {
	if g.cfg == nil {
		return "", fmt.Errorf("config cannot be nil")
	}

	text := g.cfg.Copyright.Template
	if text == "" {
		text = config.DefaultTemplate
	}

	tmpl, err := template.New("copyright").Parse(text)
	if err != nil {
		return "", fmt.Errorf("failed to parse copyright template: %w", err)
	}

	spdx := g.cfg.Copyright.SPDXIdentifier
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, headerData{
		Year:           g.cfg.GetYear(),
		Company:        g.cfg.Copyright.Company,
		Author:         g.cfg.Copyright.Author,
		SPDX:           spdx,
		SPDXIdentifier: spdx,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render copyright template: %w", err)
	}
	return buf.String(), nil
}
