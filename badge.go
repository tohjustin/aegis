package main

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr"
)

const (
	badgeTemplateDirectory = "./assets/badge-templates"
)

var badgeColors = map[string]string{
	"blue":        "#007ec6",
	"brightgreen": "#4c1",
	"green":       "#97CA00",
	"yellowgreen": "#a4a61d",
	"yellow":      "#dfb317",
	"orange":      "#fe7d37",
	"red":         "#e05d44",
	"default":     "#e05d44",
	"grey":        "#555",
	"gray":        "#555",
	"lightgrey":   "#9f9f9f",
	"lightgray":   "#9f9f9f",
}

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

func isValidHexColor(str string) bool {
	hexColorPattern := regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	matched := hexColorPattern.FindStringSubmatch(str)
	return matched != nil
}

func parseHexColor(str string) string {
	if color := "#" + str; isValidHexColor(color) {
		return color
	}

	if color, ok := badgeColors[str]; ok {
		return color
	}

	return badgeColors["default"]
}

func newBadge(badgeType, subject, status, color string) (badge, error) {
	var svgBadge badge

	switch badgeType {
	case "semaphore":
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
	case "classic":
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

func generateBadge(badgeType, subject, status, color string) (string, error) {
	svgBadge, err := newBadge(badgeType, subject, status, color)
	if err != nil {
		return "", err
	}

	badgeTemplate := packr.NewBox(badgeTemplateDirectory).String(svgBadge.TemplateFilename)
	t := template.New(svgBadge.TemplateFilename)
	t.Funcs(template.FuncMap{
		"add":      func(a, b int) int { return a + b },
		"multiply": func(a, b int) int { return a * b },
	})
	t.Parse(badgeTemplate)

	var svgBuffer bytes.Buffer
	err = t.Execute(&svgBuffer, svgBadge)

	return svgBuffer.String(), err
}
