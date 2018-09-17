// Package badge provides functions for generating SVG badges.
package badge

//go:generate packr

import (
	"bytes"
	"encoding/base64"
	"regexp"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr"
)

// Style determines the type of badge to generate
type Style string

// List of supported badge styles
const (
	ClassicStyle   Style = "classic"
	FlatStyle      Style = "flat"
	PlasticStyle   Style = "plastic"
	SemaphoreStyle Style = "semaphore"
)

// Options holds badge parameters
type Options struct {
	// CSS color names (up to CSS Color Module Level 3) or HEX values
	// (eg. "coral", "1bacbf", "fff")
	Color string

	// Font Awesome Free 5.2.0 by @fontawesome - https://fontawesome.com
	// (eg. "brands/docker", "regular/credit-card", "solid/anchor")
	Icon string

	Style Style
}

// badge holds dimensions used in generating SVG badge
type badge struct {
	Color        string
	FontFamily   string
	FontSize     int
	PaddingInner int
	PaddingOuter int
	TotalWidth   int

	TemplateFilename string

	Status          string
	StatusFontColor string
	StatusOffset    int
	StatusTextWidth int
	StatusWidth     int

	Subject          string
	SubjectFontColor string
	SubjectOffset    int
	SubjectTextWidth int
	SubjectWidth     int

	IconBase64Str string
	IconOffset    int
}

// minifySVG minifies the SVG by removing newline, tab characters
func minifySVG(svg string) string {
	return regexp.MustCompile(`[\n\t]`).ReplaceAllString(svg, "")
}

// new computes SVG dimensions based on the given badge parameters & stores into `badge` data object
func new(subject, status string, options *Options) (badge, error) {
	badgeOptions := options
	if badgeOptions == nil {
		badgeOptions = &Options{}
	}

	var newBadge badge
	switch badgeOptions.Style {
	case FlatStyle:
		newBadge = badge{
			Color:            parseColor(badgeOptions.Color),
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
	case PlasticStyle:
		newBadge = badge{
			Color:            parseColor(badgeOptions.Color),
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
	case SemaphoreStyle:
		newBadge = badge{
			Color:            parseColor(badgeOptions.Color),
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
	case ClassicStyle:
		fallthrough
	default:
		newBadge = badge{
			Color:            parseColor(badgeOptions.Color),
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

	subjectTextWidth, err := computeTextWidth(newBadge.Subject, newBadge.FontSize,
		newBadge.FontFamily)
	if err != nil {
		return newBadge, err
	}

	statusTextWidth, err := computeTextWidth(newBadge.Status, newBadge.FontSize,
		newBadge.FontFamily)
	if err != nil {
		return newBadge, err
	}

	if badgeOptions.Icon != "" {
		svgIcon, err := packr.NewBox("./assets/icons/").MustString(badgeOptions.Icon + ".svg")
		if err == nil {
			// Set SVG icon color to match `newBadge.SubjectFontColor`,
			// include font-awesome license into the base64-encoded result
			modifiedSvgIcon := "<svg fill=\"" + newBadge.SubjectFontColor + "\"" + svgIcon[len("<svg"):]
			newBadge.IconBase64Str = base64.StdEncoding.EncodeToString([]byte(modifiedSvgIcon))
			newBadge.IconOffset = 3 + 13 // IconPadding + IconSize
		}
	}

	newBadge.SubjectOffset = newBadge.PaddingOuter + newBadge.IconOffset
	newBadge.SubjectTextWidth = subjectTextWidth
	newBadge.SubjectWidth = newBadge.SubjectOffset + subjectTextWidth + newBadge.PaddingInner

	newBadge.StatusOffset = newBadge.SubjectWidth + newBadge.PaddingInner
	newBadge.StatusTextWidth = statusTextWidth
	newBadge.StatusWidth = newBadge.PaddingInner + statusTextWidth + newBadge.PaddingOuter

	newBadge.TotalWidth = newBadge.SubjectWidth + newBadge.StatusWidth

	return newBadge, err
}

// Create generates a SVG badge
func Create(subject, status string, options *Options) (string, error) {
	newBadge, err := new(subject, status, options)
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

// CreateUnsafe generates a SVG badge (unsafe version does not return errors)
func CreateUnsafe(subject, status string, options *Options) string {
	generatedBadge, _ := Create(subject, status, options)
	return generatedBadge
}
