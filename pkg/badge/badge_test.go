package badge

import (
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name     string
	input    Params
	expected Params
}

var testCases = (func() []testCase {
	result := []testCase{}
	for _, testStyle := range append(SupportedStyles[:], "") {
		testNamePrefix := string(testStyle)
		if len(testNamePrefix) > 0 {
			testNamePrefix = strings.ToUpper(testNamePrefix[:1]) + testNamePrefix[1:]
		}
		testSubject, testStatus := "testSubject", "testStatus"
		expectedStyle, expectedSubject, expectedStatus := testStyle,
			testSubject, testStatus

		// default style case
		if testStyle == "" {
			expectedStyle = defaultStyle
		}

		// Semaphore style badges converts text to uppercase
		if testStyle == SemaphoreCIStyle {
			expectedSubject, expectedStatus = "TESTSUBJECT", "TESTSTATUS"
		}

		result = append(result, []testCase{
			{
				name: testNamePrefix + "Badge",
				input: Params{
					Style:   testStyle,
					Subject: testSubject,
					Status:  testStatus,
					Color:   "#f7b137",
					Icon:    "solid/star",
				},
				expected: Params{
					Style:   expectedStyle,
					Subject: expectedSubject,
					Status:  expectedStatus,
					Color:   "#f7b137",
					Icon:    "solid/star",
				},
			},
			{
				name:     testNamePrefix + "BadgeWithDefaultParams",
				input:    Params{},
				expected: Params{Style: defaultStyle, Color: defaultColor},
			},
			{
				name:     testNamePrefix + "BadgeWithColorName1",
				input:    Params{Style: testStyle, Color: "red"},
				expected: Params{Style: expectedStyle, Color: "red"},
			},
			{
				name:     testNamePrefix + "BadgeWithColorName2",
				input:    Params{Style: testStyle, Color: "RED"},
				expected: Params{Style: expectedStyle, Color: "red"},
			},
			{
				name:     testNamePrefix + "BadgeWithInvalidColorName",
				input:    Params{Style: testStyle, Color: "rainbow"},
				expected: Params{Style: expectedStyle, Color: defaultColor},
			},
			{
				name:     testNamePrefix + "BadgeWithShortHexColorCode1",
				input:    Params{Style: testStyle, Color: "#abc"},
				expected: Params{Style: expectedStyle, Color: "#abc"},
			},
			{
				name:     testNamePrefix + "BadgeWithShortHexColorCode2",
				input:    Params{Style: testStyle, Color: "#ABC"},
				expected: Params{Style: expectedStyle, Color: "#abc"},
			},
			{
				name:     testNamePrefix + "BadgeWithHexColorCode1",
				input:    Params{Style: testStyle, Color: "#f7b137"},
				expected: Params{Style: expectedStyle, Color: "#f7b137"},
			},
			{
				name:     testNamePrefix + "BadgeWithHexColorCode2",
				input:    Params{Style: testStyle, Color: "#F7B137"},
				expected: Params{Style: expectedStyle, Color: "#f7b137"},
			},
			{
				name:     testNamePrefix + "BadgeWithHexColorCode2",
				input:    Params{Style: testStyle, Color: "#F7B137"},
				expected: Params{Style: expectedStyle, Color: "#f7b137"},
			},
			{
				name:     testNamePrefix + "BadgeWithInvalidColorCode",
				input:    Params{Style: testStyle, Color: "#f7b137aa"},
				expected: Params{Style: expectedStyle, Color: defaultColor},
			},
			{
				name:     testNamePrefix + "BadgeWithIcon",
				input:    Params{Style: testStyle, Icon: "solid/star"},
				expected: Params{Style: expectedStyle, Color: defaultColor, Icon: "solid/star"},
			},
			{
				name:     testNamePrefix + "BadgeWithInvalidIcon1",
				input:    Params{Style: testStyle, Icon: "solid/STAR"},
				expected: Params{Style: expectedStyle, Color: defaultColor, Icon: ""},
			},
			{
				name:     testNamePrefix + "BadgeWithInvalidIcon2",
				input:    Params{Style: testStyle, Icon: "invalid-icon"},
				expected: Params{Style: expectedStyle, Color: defaultColor},
			},
		}...)
	}

	return result
})()

func TestSnapshotBadgeCreate(t *testing.T) {
	t.Parallel()

	for _, spec := range testCases {
		t.Run(spec.name, func(t *testing.T) {
			result, err := Create(&spec.input)
			if err != nil {
				t.Error(err)
			}

			cupaloy.SnapshotT(t, result)
		})
	}
}

func TestBadgeCreate(t *testing.T) {
	t.Parallel()

	for _, spec := range testCases {
		t.Run(spec.name, func(t *testing.T) {
			newBadge, err := Create(&spec.input)
			if err != nil {
				t.Fatal(err)
			}

			newBadgeParams, err := ExtractParams(newBadge)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, spec.expected, *newBadgeParams)
		})
	}
}
