package badge

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/gobuffalo/packr"
)

var minifySVGTestCases = map[string]string{
	"classicSVGTemplate":   "classic.tmpl",
	"flatSVGTemplate":      "flat.tmpl",
	"plasticSVGTemplate":   "plastic.tmpl",
	"semaphoreSVGTemplate": "semaphore.tmpl",
}

var badgeTestCases = map[string][]string{
	"ClassicBadgeWithColorName":      {"classic", "testSubject", "testStatus", "red", ""},
	"ClassicBadgeWithHexColorCode":   {"classic", "testSubject", "testStatus", "abc", ""},
	"FlatBadgeWithColorName":         {"flat", "testSubject", "testStatus", "blue", ""},
	"FlatBadgeWithHexColorCode":      {"flat", "testSubject", "testStatus", "abcdef", ""},
	"SemaphoreBadgeWithColorName":    {"semaphore", "testSubject", "testStatus", "yellow", ""},
	"SemaphoreBadgeWithHexColorCode": {"semaphore", "testSubject", "testStatus", "abcdef", ""},
	"PlasticBadgeWithColorName":      {"plastic", "testSubject", "testStatus", "green", ""},
	"PlasticBadgeWithHexColorCode":   {"plastic", "testSubject", "testStatus", "abcdef", ""},
}

func TestMinifySVG(t *testing.T) {
	for testName, templateName := range minifySVGTestCases {
		t.Run(testName, func(t *testing.T) {
			badgeSVGTemplate := packr.NewBox("./assets/badge-templates").String(templateName)
			result := minifySVG(badgeSVGTemplate)
			cupaloy.SnapshotT(t, result)
		})
	}
}
func TestNewBadge(t *testing.T) {
	for testName, args := range badgeTestCases {
		t.Run(testName, func(t *testing.T) {
			badgeStyle, subject, status, color, icon := args[0], args[1], args[2], args[3], args[4]
			result, _ := newBadge(badgeStyle, subject, status, color, icon)
			cupaloy.SnapshotT(t, result)
		})
	}
}

func TestGenerateSVG(t *testing.T) {
	for testName, args := range badgeTestCases {
		t.Run(testName, func(t *testing.T) {
			badgeStyle, subject, status, color, icon := args[0], args[1], args[2], args[3], args[4]
			result, _ := GenerateSVG(badgeStyle, subject, status, color, icon)
			cupaloy.SnapshotT(t, result)
		})
	}
}
