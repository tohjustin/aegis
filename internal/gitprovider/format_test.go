package gitprovider

import (
	"strconv"
	"testing"
)

func TestFormatCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
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

	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.input), func(t *testing.T) {
			t.Parallel()
			got := FormatCount(tt.input)
			if got != tt.expected {
				t.Errorf("FormatCount(%d) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
