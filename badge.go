package main

import (
	"bytes"
	"text/template"

	"github.com/gobuffalo/packr"
)

type Badge struct {
	Subject           string
	Status            string
	Color             string
	CenterPadding     int
	SidePadding       int
	FontSize          int
	FontFamily        string
	TemplateFilename  string
	StatusTextOffset  int
	StatusTextWidth   int
	StatusWidth       int
	SubjectTextOffset int
	SubjectTextWidth  int
	SubjectWidth      int
	TotalWidth        int
}

func generateBadge(svgBadge Badge) (string, error) {
	scale := 10

	subjectTextWidth, err := computeTextWidth(svgBadge.Subject, svgBadge.FontSize, svgBadge.FontFamily)
	if err != nil {
		return "", err
	}

	statusTextWidth, err := computeTextWidth(svgBadge.Status, svgBadge.FontSize, svgBadge.FontFamily)
	if err != nil {
		return "", err
	}

	// Multiply dimensions by factor of 10 (& scaling down by 10 in SVG) to avoid using float64's (eg. 375 instead of 37.5)
	svgBadge.SubjectWidth = subjectTextWidth + svgBadge.SidePadding + svgBadge.CenterPadding
	svgBadge.StatusWidth = statusTextWidth + svgBadge.SidePadding + svgBadge.CenterPadding
	svgBadge.TotalWidth = svgBadge.SubjectWidth + svgBadge.StatusWidth
	svgBadge.SubjectTextWidth = subjectTextWidth * scale
	svgBadge.SubjectTextOffset = svgBadge.SidePadding * scale
	svgBadge.StatusTextWidth = statusTextWidth * scale
	svgBadge.StatusTextOffset = (svgBadge.SubjectWidth + svgBadge.CenterPadding) * scale

	var svgBuffer bytes.Buffer
	t := template.New("")
	t.Parse(packr.NewBox("./assets/badge-templates").String(svgBadge.TemplateFilename))
	t.Execute(&svgBuffer, svgBadge)

	return svgBuffer.String(), nil
}

func generateClassicBadge(subject string, status string, color string) (string, error) {
	svgBadge := Badge{
		Color:            color,
		Status:           status,
		Subject:          subject,
		CenterPadding:    4,
		SidePadding:      6,
		FontSize:         11,
		FontFamily:       "Verdana",
		TemplateFilename: "classic.tmpl",
	}

	return generateBadge(svgBadge)
}
