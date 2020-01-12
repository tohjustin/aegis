// Package badge provides functions for generating SVG badges.
package badge

//go:generate go run gen.go

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"
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

// SupportedStyles contains a list of all supported badge styles
var SupportedStyles = [...]Style{ClassicStyle, FlatStyle, PlasticStyle, SemaphoreStyle}

// Params holds badge parameters
type Params struct {
	// Subject determines the subject text of the badge.
	Subject string
	// Status determines the status text of the badge.
	Status string
	// Color determines the highlight color of the badge.
	// Valid color values includes CSS color names (up to CSS Color Module Level 3) or HEX values (eg. "coral", "#1bacbf", "1bacbf", "fff", "#fff")
	Color string
	// Icon determines whether the badge should include icons or not (eg. "brands/docker", "regular/credit-card", "solid/anchor")
	Icon string
	// Style determines the visual style of the badge
	Style Style
}

// badgeDimensions holds dimensions required for generating SVG badge
type badgeDimensions struct {
	Style        Style
	Template     *template.Template
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

	IconLabel     string
	IconBase64Str string
	IconOffset    int
}

const defaultColor = "#1bacbf"
const defaultStyle = ClassicStyle

// generateBadge converts badge parameters into dimensions for generating SVG badge
func generateBadge(params *Params) (*badgeDimensions, error) {
	badgeParams := params
	if badgeParams == nil {
		badgeParams = &Params{}
	}
	badgeColor := parseColor(badgeParams.Color)
	if badgeColor == "" {
		badgeColor = defaultColor
	}
	badgeStyle := badgeParams.Style
	if badgeStyle == Style("") {
		badgeStyle = defaultStyle
	}

	var newBadge badgeDimensions
	switch badgeStyle {
	case FlatStyle:
		newBadge = badgeDimensions{
			Style:            FlatStyle,
			Template:         badgeTemplates[FlatStyle],
			Color:            badgeColor,
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           badgeParams.Status,
			StatusFontColor:  "#fff",
			Subject:          badgeParams.Subject,
			SubjectFontColor: "#fff",
		}
	case PlasticStyle:
		newBadge = badgeDimensions{
			Style:            PlasticStyle,
			Template:         badgeTemplates[PlasticStyle],
			Color:            badgeColor,
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           badgeParams.Status,
			StatusFontColor:  "#fff",
			Subject:          badgeParams.Subject,
			SubjectFontColor: "#fff",
		}
	case SemaphoreStyle:
		newBadge = badgeDimensions{
			Style:            SemaphoreStyle,
			Template:         badgeTemplates[SemaphoreStyle],
			Color:            badgeColor,
			FontFamily:       "Verdana",
			FontSize:         9,
			PaddingInner:     10,
			PaddingOuter:     10,
			Status:           strings.ToUpper(badgeParams.Status),
			StatusFontColor:  "#fff",
			Subject:          strings.ToUpper(badgeParams.Subject),
			SubjectFontColor: "#888",
		}
	case ClassicStyle:
		fallthrough
	default:
		newBadge = badgeDimensions{
			Style:            ClassicStyle,
			Template:         badgeTemplates[ClassicStyle],
			Color:            badgeColor,
			FontFamily:       "Verdana",
			FontSize:         11,
			PaddingInner:     4,
			PaddingOuter:     6,
			Status:           badgeParams.Status,
			StatusFontColor:  "#fff",
			Subject:          badgeParams.Subject,
			SubjectFontColor: "#fff",
		}
	}

	subjectTextWidth, err := computeTextWidth(newBadge.Subject, newBadge.FontSize,
		newBadge.FontFamily)
	if err != nil {
		return nil, err
	}

	statusTextWidth, err := computeTextWidth(newBadge.Status, newBadge.FontSize,
		newBadge.FontFamily)
	if err != nil {
		return nil, err
	}

	if badgeParams.Icon != "" {
		svgIcon, ok := fontAwesomeIcons[badgeParams.Icon]
		if ok {
			// Encode icon into a base64 string
			modifiedSvgIcon := "<svg fill=\"" + newBadge.SubjectFontColor + "\"" + svgIcon[len("<svg"):]
			newBadge.IconLabel = badgeParams.Icon
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

	if newBadge.Template == nil {
		return nil, fmt.Errorf("Badge template does not exist: %s", params.Style)
	}

	return &newBadge, nil
}

// Create generates a SVG badge
func Create(params *Params) (string, error) {
	newBadge, err := generateBadge(params)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = newBadge.Template.Execute(&buf, newBadge); err != nil {
		return "", err
	}

	return buf.String(), nil
}
