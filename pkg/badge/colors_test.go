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

var cssColorTestCases = []struct {
	input  string
	output bool
}{
	{"red", true},            // CSS Level 1
	{"orange", true},         // CSS Level 2
	{"orchid", true},         // CSS Level 3
	{"rebeccapurple", false}, // CSS Level 4 (not supported)
}

func TestIsValidHexColor(t *testing.T) {
	t.Parallel()

	for _, testCase := range hexColorTestCases {
		t.Run(testCase.input, func(t *testing.T) {
			assert.Equal(t, isValidHexColor(testCase.input), testCase.output)
		})
	}
}

func TestIsValidCSSColorName(t *testing.T) {
	t.Parallel()

	for _, testCase := range cssColorTestCases {
		t.Run(testCase.input, func(t *testing.T) {
			assert.Equal(t, isValidCSSColorName(testCase.input), testCase.output)
		})
	}
}
func TestParseColor(t *testing.T) {
	t.Parallel()

	for cssColorName := range cssColorNames {
		t.Run("TestParseColor-"+cssColorName, func(t *testing.T) {
			result := parseColor(cssColorName)
			assert.Equal(t, result, cssColorName)
		})
	}

	assert.Equal(t, "#f35f35", parseColor("f35f35"))
	assert.Equal(t, defaultColor, parseColor("f35f"))
}
