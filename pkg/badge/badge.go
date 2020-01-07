// Package badge provides functions for generating SVG badges.
package badge

//go:generate go run gen.go

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
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
	// Color determines the highlight color of the badge.
	// Valid color values includes CSS color names (up to CSS Color Module Level 3) or HEX values (eg. "coral", "1bacbf", "fff")
	Color string
	// Icon determines whether the badge should include icons or not (eg. "brands/docker", "regular/credit-card", "solid/anchor")
	Icon string
	// Style determines the visual style of the badge
	Style Style
}

// badge holds dimensions used in generating SVG badge
type badge struct {
	Style        Style
	Color        string
	FontFamily   string
	FontSize     int
	PaddingInner int
	PaddingOuter int
	TotalWidth   int

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
			Style:            FlatStyle,
			Color:            parseColor(badgeOptions.Color),
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           status,
			StatusFontColor:  "#fff",
			Subject:          subject,
			SubjectFontColor: "#fff",
		}
	case PlasticStyle:
		newBadge = badge{
			Style:            PlasticStyle,
			Color:            parseColor(badgeOptions.Color),
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           status,
			StatusFontColor:  "#fff",
			Subject:          subject,
			SubjectFontColor: "#fff",
		}
	case SemaphoreStyle:
		newBadge = badge{
			Style:            SemaphoreStyle,
			Color:            parseColor(badgeOptions.Color),
			FontFamily:       "Verdana",
			FontSize:         9,
			PaddingInner:     10,
			PaddingOuter:     10,
			Status:           strings.ToUpper(status),
			StatusFontColor:  "#fff",
			Subject:          strings.ToUpper(subject),
			SubjectFontColor: "#888",
		}
	case ClassicStyle:
		fallthrough
	default:
		newBadge = badge{
			Style:            ClassicStyle,
			Color:            parseColor(badgeOptions.Color),
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           status,
			StatusFontColor:  "#fff",
			Subject:          subject,
			SubjectFontColor: "#fff",
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
		svgIcon, ok := fontAwesomeIcons[badgeOptions.Icon]
		if ok {
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

	t, ok := badgeTemplates[newBadge.Style]
	if !ok {
		return "", fmt.Errorf("badge template does not exist: %s", options.Style)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, newBadge)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
