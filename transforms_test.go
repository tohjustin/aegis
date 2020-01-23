package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatIntegerWithMetricPrefix(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{11, "11"},
		{112, "112"},
		{1122, "1.12k"},
		{11223, "11.2k"},
		{112233, "112k"},
		{1122334, "1.12M"},
		{11223344, "11.2M"},
		{112233445, "112M"},
		{1122334455, "1.12G"},
		{11223344556, "11.2G"},
		{112233445566, "112G"},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.input), func(t *testing.T) {
			assert.Equal(t, testCase.expected, formatIntegerWithMetricPrefix(testCase.input))
		})
	}
}
