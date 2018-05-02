package badge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var hexColorTestCases = []struct {
	input  string
	output bool
}{
	{"000", false},
	{"f35", false},
	{"F35", false},
	{"fff", false},
	{"FFF", false},
	{"000000", false},
	{"f35f35", false},
	{"F35F35", false},
	{"F35f35", false},
	{"ffffff", false},
	{"FFFFFF", false},
	{"FFffFF", false},
	{"red", false},
	{"#ffffffaa", false},
	{"#FFFFFFAA", false},
	{"#000", true},
	{"#f35", true},
	{"#F35", true},
	{"#fff", true},
	{"#FFF", true},
	{"#000000", true},
	{"#f35f35", true},
	{"#F35F35", true},
	{"#F35f35", true},
	{"#ffffff", true},
	{"#FFFFFF", true},
	{"#FFffFF", true},
}

func TestIsValidHexColor(t *testing.T) {
	for _, testCase := range hexColorTestCases {
		t.Run(testCase.input, func(t *testing.T) {
			assert.Equal(t, isValidHexColor(testCase.input), testCase.output)
		})
	}
}

func TestParseHexColor(t *testing.T) {
	for input, output := range BadgeColors {
		result := parseHexColor(input)
		assert.Equal(t, result, output)
	}

	assert.Equal(t, "#f35f35", parseHexColor("f35f35"))
	assert.Equal(t, BadgeColors["default"], parseHexColor("f35f"))
}
