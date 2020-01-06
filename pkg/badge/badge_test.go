package badge

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

type badgeTestCase struct {
	subject string
	status  string
	options Options
}

var badgeTestCases = map[string]badgeTestCase{
	"ClassicBadgeWithColorName": {
		subject: "testSubject",
		status:  "testStatus",
		options: Options{Color: "red", Style: ClassicStyle},
	},
	"ClassicBadgeWithHexColorCode": {
		subject: "testSubject",
		status:  "testStatus",
		options: Options{Color: "abc", Style: ClassicStyle},
	},
	"FlatBadgeWithColorName": {
		subject: "testSubject",
		status:  "testStatus",
		options: Options{Color: "blue", Style: FlatStyle},
	},
	"FlatBadgeWithHexColorCode": {
		subject: "testSubject",
		status:  "testStatus",
		options: Options{Color: "abcdef", Style: FlatStyle},
	},
	"SemaphoreBadgeWithColorName": {
		subject: "testSubject",
		status:  "testStatus",
		options: Options{Color: "yellow", Style: SemaphoreStyle},
	},
	"SemaphoreBadgeWithHexColorCode": {
		subject: "testSubject",
		status:  "testStatus",
		options: Options{Color: "abcdef", Style: SemaphoreStyle},
	},
	"PlasticBadgeWithColorName": {
		subject: "testSubject",
		status:  "testStatus",
		options: Options{Color: "green", Style: PlasticStyle},
	},
	"PlasticBadgeWithHexColorCode": {
		subject: "testSubject",
		status:  "testStatus",
		options: Options{Color: "abcdef", Style: PlasticStyle},
	},
}

func TestBadgeNew(t *testing.T) {
	t.Parallel()

	for testName, testParams := range badgeTestCases {
		t.Run(testName, func(t *testing.T) {
			result, _ := new(testParams.subject, testParams.status, &testParams.options)
			cupaloy.SnapshotT(t, result)
		})
	}
}

func TestBadgeCreate(t *testing.T) {
	t.Parallel()

	for testName, testParams := range badgeTestCases {
		t.Run(testName, func(t *testing.T) {
			result, _ := Create(testParams.subject, testParams.status, &testParams.options)
			cupaloy.SnapshotT(t, result)
		})
	}
}
