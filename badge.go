package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr"
)

type badge struct {
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

func generateBadge(svgBadge badge) (string, error) {
	scale := 10

	subjectTextWidth, err := computeTextWidth(svgBadge.Subject, svgBadge.FontSize, svgBadge.FontFamily)
	if err != nil {
		return "", err
	}

	statusTextWidth, err := computeTextWidth(svgBadge.Status, svgBadge.FontSize, svgBadge.FontFamily)
	if err != nil {
		return "", err
	}

	// Multiply dimensions by factor of 10 (& later scaling down by 10 in SVG template) to avoid
	// using float64's (eg. 375 instead of 37.5)
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
	svgBadge := badge{
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

func generateSemaphoreBadge(subject string, status string, color string) (string, error) {
	// Use larger font-size of 10 ("semaphore.tmpl" uses font-size of 9) to
	// increase letter-spacing of subject & status text
	svgBadge := badge{
		Color:            color,
		Status:           strings.ToUpper(status),
		Subject:          strings.ToUpper(subject),
		CenterPadding:    10,
		SidePadding:      10,
		FontSize:         10,
		FontFamily:       "Verdana",
		TemplateFilename: "semaphore.tmpl",
	}

	return generateBadge(svgBadge)
}
