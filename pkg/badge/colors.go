package badge

import "regexp"

// BadgeColors A map between color names & their corresponding hex values
var BadgeColors = map[string]string{
	"blue":        "#007ec6",
	"brightgreen": "#4c1",
	"green":       "#97CA00",
	"yellowgreen": "#a4a61d",
	"yellow":      "#dfb317",
	"orange":      "#fe7d37",
	"red":         "#e05d44",
	"default":     "#e05d44",
	"grey":        "#555",
	"gray":        "#555",
	"lightgrey":   "#9f9f9f",
	"lightgray":   "#9f9f9f",
}

func isValidHexColor(str string) bool {
	hexColorPattern := regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	matched := hexColorPattern.FindStringSubmatch(str)
	return matched != nil
}

func parseHexColor(str string) string {
	if color := "#" + str; isValidHexColor(color) {
		return color
	}

	if color, ok := BadgeColors[str]; ok {
		return color
	}

	return BadgeColors["default"]
}
