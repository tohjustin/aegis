package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var hexColorTests = []struct {
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
	for _, testcase := range hexColorTests {
		t.Run(testcase.input, func(t *testing.T) {
			assert.Equal(t, isValidHexColor(testcase.input), testcase.output)
		})
	}
}

func TestParseHexColor(t *testing.T) {
	for colorName, hexValue := range badgeColors {
		assert.Equal(t, hexValue, parseHexColor(colorName))
	}

	assert.Equal(t, "#f35f35", parseHexColor("f35f35"))
	assert.Equal(t, badgeColors["default"], parseHexColor("f35f"))
}
