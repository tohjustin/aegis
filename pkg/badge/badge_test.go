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
	"ClassicBadgeWithColorName":      []string{"classic", "testSubject", "testStatus", "red"},
	"ClassicBadgeWithHexColorCode":   []string{"classic", "testSubject", "testStatus", "abc"},
	"FlatBadgeWithColorName":         []string{"flat", "testSubject", "testStatus", "blue"},
	"FlatBadgeWithHexColorCode":      []string{"flat", "testSubject", "testStatus", "abcdef"},
	"SemaphoreBadgeWithColorName":    []string{"semaphore", "testSubject", "testStatus", "yellow"},
	"SemaphoreBadgeWithHexColorCode": []string{"semaphore", "testSubject", "testStatus", "abcdef"},
	"PlasticBadgeWithColorName":      []string{"plastic", "testSubject", "testStatus", "green"},
	"PlasticBadgeWithHexColorCode":   []string{"plastic", "testSubject", "testStatus", "abcdef"},
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
			badgeStyle, subject, status, color := args[0], args[1], args[2], args[3]
			result, _ := newBadge(badgeStyle, subject, status, color)
			cupaloy.SnapshotT(t, result)
		})
	}
}

func TestGenerateSVG(t *testing.T) {
	for testName, args := range badgeTestCases {
		t.Run(testName, func(t *testing.T) {
			badgeStyle, subject, status, color := args[0], args[1], args[2], args[3]
			result, _ := GenerateSVG(badgeStyle, subject, status, color)
			cupaloy.SnapshotT(t, result)
		})
	}
}
