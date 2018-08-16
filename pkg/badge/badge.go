// Package badge provides function(s) for creating SVG badges.
//
// Badge designs are based on Shields.IO's specification found in
// https://github.com/badges/shields/blob/master/spec/SPECIFICATION.md
//
package badge

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr"
)

// Supported badge styles
const (
	BadgeStyleClassic   = "classic"
	BadgeStyleFlat      = "flat"
	BadgeStylePlastic   = "plastic"
	BadgeStyleSemaphore = "semaphore"
)

type badge struct {
	Color            string
	FontFamily       string
	FontSize         int
	PaddingInner     int
	PaddingOuter     int
	Status           string
	StatusFontColor  string
	StatusOffset     int
	StatusTextWidth  int
	StatusWidth      int
	Subject          string
	SubjectFontColor string
	SubjectOffset    int
	SubjectTextWidth int
	SubjectWidth     int
	TemplateFilename string
	TotalWidth       int
}

// minifies SVG by removing newline & tab characters
func minifySVG(svg string) string {
	result := svg
	result = strings.Replace(result, "\n", "", -1)
	result = strings.Replace(result, "\t", "", -1)
	return result
}

func newBadge(style, subject, status, color string) (badge, error) {
	var svgBadge badge

	switch style {
	case BadgeStyleFlat:
		svgBadge = badge{
			Color:            parseColor(color),
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           status,
			StatusFontColor:  "#fff",
			Subject:          subject,
			SubjectFontColor: "#fff",
			TemplateFilename: "flat.tmpl",
		}
	case BadgeStylePlastic:
		svgBadge = badge{
			Color:            parseColor(color),
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           status,
			StatusFontColor:  "#fff",
			Subject:          subject,
			SubjectFontColor: "#fff",
			TemplateFilename: "plastic.tmpl",
		}
	case BadgeStyleSemaphore:
		svgBadge = badge{
			Color:            parseColor(color),
			FontFamily:       "Verdana",
			FontSize:         9,
			PaddingInner:     10,
			PaddingOuter:     10,
			Status:           strings.ToUpper(status),
			StatusFontColor:  "#fff",
			Subject:          strings.ToUpper(subject),
			SubjectFontColor: "#888",
			TemplateFilename: "semaphore.tmpl",
		}
	case BadgeStyleClassic:
		fallthrough
	default:
		svgBadge = badge{
			Color:            parseColor(color),
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           status,
			StatusFontColor:  "#fff",
			Subject:          subject,
			SubjectFontColor: "#fff",
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

	svgBadge.SubjectOffset = svgBadge.PaddingOuter
	svgBadge.SubjectTextWidth = subjectTextWidth
	svgBadge.SubjectWidth = svgBadge.SubjectOffset + subjectTextWidth + svgBadge.PaddingInner

	svgBadge.StatusOffset = svgBadge.SubjectWidth + svgBadge.PaddingInner
	svgBadge.StatusTextWidth = statusTextWidth
	svgBadge.StatusWidth = svgBadge.PaddingInner + statusTextWidth + svgBadge.PaddingOuter

	svgBadge.TotalWidth = svgBadge.SubjectWidth + svgBadge.StatusWidth

	return svgBadge, nil
}

// GenerateSVG returns a string representation of the generated SVG badge
//
func GenerateSVG(style, subject, status, color string) (string, error) {
	newBadge, err := newBadge(style, subject, status, color)
	if err != nil {
		return "", err
	}

	badgeSVGTemplate := packr.NewBox("./assets/badge-templates").String(newBadge.TemplateFilename)
	t := template.New(newBadge.TemplateFilename)
	t.Parse(badgeSVGTemplate)

	var buf bytes.Buffer
	err = t.Execute(&buf, newBadge)
	if err != nil {
		return "", err
	}

	return minifySVG(buf.String()), nil
}
