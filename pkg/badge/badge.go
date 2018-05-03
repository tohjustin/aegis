// Package badge provides functions for creating SVG badges
package badge

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr"
)

// Supported badge types
const (
	Classic   = "classic"
	Flat      = "flat"
	Plastic   = "plastic"
	Semaphore = "semaphore"
)

type badge struct {
	Subject          string
	Status           string
	Color            string
	InnerPadding     int
	OuterPadding     int
	FontSize         int
	FontFamily       string
	TemplateFilename string
	StatusTextWidth  int
	StatusWidth      int
	SubjectTextWidth int
	SubjectWidth     int
}

// minifySVG minifies SVG by removing newline & tab characters
func minifySVG(svg string) string {
	result := svg
	result = strings.Replace(result, "\n", "", -1)
	result = strings.Replace(result, "\t", "", -1)
	return result
}

func newBadge(badgeStyle, subject, status, color string) (badge, error) {
	var svgBadge badge

	switch badgeStyle {
	case Flat:
		svgBadge = badge{
			Color:            parseHexColor(color),
			Status:           status,
			Subject:          subject,
			InnerPadding:     4,
			OuterPadding:     6,
			FontSize:         11,
			FontFamily:       "Verdana",
			TemplateFilename: "flat.tmpl",
		}
	case Plastic:
		svgBadge = badge{
			Color:            parseHexColor(color),
			Status:           status,
			Subject:          subject,
			InnerPadding:     4,
			OuterPadding:     6,
			FontSize:         11,
			FontFamily:       "Verdana",
			TemplateFilename: "plastic.tmpl",
		}
	case Semaphore:
		svgBadge = badge{
			Color:            parseHexColor(color),
			Status:           strings.ToUpper(status),
			Subject:          strings.ToUpper(subject),
			InnerPadding:     10,
			OuterPadding:     10,
			FontSize:         10,
			FontFamily:       "Verdana",
			TemplateFilename: "semaphore.tmpl",
		}
	case Classic:
		fallthrough
	default:
		svgBadge = badge{
			Color:            parseHexColor(color),
			Status:           status,
			Subject:          subject,
			InnerPadding:     4,
			OuterPadding:     6,
			FontSize:         11,
			FontFamily:       "Verdana",
			TemplateFilename: "classic.tmpl",
		}
	}

	subjectTextWidth, err := computeTextWidth(svgBadge.Subject, svgBadge.FontSize, svgBadge.FontFamily)
	if err != nil {
		return svgBadge, err
	}

	statusTextWidth, err := computeTextWidth(svgBadge.Status, svgBadge.FontSize, svgBadge.FontFamily)
	if err != nil {
		return svgBadge, err
	}

	svgBadge.SubjectWidth = subjectTextWidth + svgBadge.OuterPadding + svgBadge.InnerPadding
	svgBadge.StatusWidth = statusTextWidth + svgBadge.OuterPadding + svgBadge.InnerPadding
	svgBadge.SubjectTextWidth = subjectTextWidth
	svgBadge.StatusTextWidth = statusTextWidth

	return svgBadge, nil
}

// GenerateSVG generates a SVG badge based on the provided arguments
// (badgeStyle, subject, status, color)
func GenerateSVG(badgeStyle, subject, status, color string) (string, error) {
	newBadge, err := newBadge(badgeStyle, subject, status, color)
	if err != nil {
		return "", err
	}

	badgeSVGTemplate := packr.NewBox("./assets/badge-templates").String(newBadge.TemplateFilename)
	t := template.New(newBadge.TemplateFilename)
	t.Funcs(template.FuncMap{
		"add":      func(a, b int) int { return a + b },
		"multiply": func(a, b int) int { return a * b },
	})
	t.Parse(badgeSVGTemplate)

	var buf bytes.Buffer
	err = t.Execute(&buf, newBadge)
	if err != nil {
		return "", err
	}

	return minifySVG(buf.String()), nil
}
