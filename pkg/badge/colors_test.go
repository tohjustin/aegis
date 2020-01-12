package badge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var hexColorTestCases = []struct {
	input    string
	expected bool
}{
	{"000", true},
	{"f35", true},
	{"F35", true},
	{"fff", true},
	{"FFF", true},
	{"000000", true},
	{"f35f35", true},
	{"F35F35", true},
	{"F35f35", true},
	{"ffffff", true},
	{"FFFFFF", true},
	{"FFffFF", true},
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
	{"red", false},
	{"#ffffffaa", false},
	{"#FFFFFFAA", false},
}

var cssColorTestCases = []struct {
	input    string
	expected bool
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
			assert.Equal(t, isValidHexColor(testCase.input), testCase.expected)
		})
	}
}

func TestIsValidCSSColorName(t *testing.T) {
	t.Parallel()

	for _, testCase := range cssColorTestCases {
		t.Run(testCase.input, func(t *testing.T) {
			assert.Equal(t, isValidCSSColorName(testCase.input), testCase.expected)
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

	assert.Equal(t, "#f7b137", parseColor("#F7B137"))
	assert.Equal(t, "#f7b137", parseColor("F7B137"))
	assert.Equal(t, "#f7b137", parseColor("#f7b137"))
	assert.Equal(t, "#f7b137", parseColor("f7b137"))
	assert.Equal(t, "#f7b", parseColor("#F7B"))
	assert.Equal(t, "#f7b", parseColor("F7B"))
	assert.Equal(t, "#f7b", parseColor("#f7b"))
	assert.Equal(t, "#f7b", parseColor("f7b"))
	assert.Equal(t, "", parseColor("#f7b137aa"))
	assert.Equal(t, "", parseColor("#f7baa"))
	assert.Equal(t, "", parseColor("#f7b1"))
}
