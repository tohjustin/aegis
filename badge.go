package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr"
)

const (
	badgeTemplateDirectory = "./assets/badge-templates"
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

func newBadge(badgeType, subject, status, color string) (badge, error) {
	var svgBadge badge

	switch badgeType {
	case "semaphore":
		svgBadge = badge{
			Color:            "#" + color,
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
			Color:            "#" + color,
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
