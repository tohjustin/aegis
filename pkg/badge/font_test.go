package badge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeasureTextWidth(t *testing.T) {
	t.Parallel()

	expectedLength := 153
	width, err := computeTextWidth("Lorem ipsum dolor sit amet", 11, "Verdana")
	if err != nil {
		t.Errorf("Fail to measure text width: %s", err)
	}

	assert.Equal(t, width, expectedLength)
}
func TestMeasureTextWidthWithMissingFont(t *testing.T) {
	t.Parallel()

	expectedError := "unable to parse \"UNKNOWN_FONT.ttf\""
	_, err := computeTextWidth("Lorem ipsum dolor sit amet", 11, "UNKNOWN_FONT")
	if err != nil {
		assert.Equal(t, err.Error(), expectedError)
	} else {
		t.Errorf("Should return an error")
	}
}
