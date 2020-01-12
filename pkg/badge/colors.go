package badge

import (
	"regexp"
	"strings"
)

var cssColorNames = map[string]struct{}{
	"aliceblue":            {},
	"antiquewhite":         {},
	"aqua":                 {},
	"aquamarine":           {},
	"azure":                {},
	"beige":                {},
	"bisque":               {},
	"black":                {},
	"blanchedalmond":       {},
	"blue":                 {},
	"blueviolet":           {},
	"brown":                {},
	"burlywood":            {},
	"cadetblue":            {},
	"chartreuse":           {},
	"chocolate":            {},
	"coral":                {},
	"cornflowerblue":       {},
	"cornsilk":             {},
	"crimson":              {},
	"cyan":                 {},
	"darkblue":             {},
	"darkcyan":             {},
	"darkgoldenrod":        {},
	"darkgray":             {},
	"darkgreen":            {},
	"darkgrey":             {},
	"darkkhaki":            {},
	"darkmagenta":          {},
	"darkolivegreen":       {},
	"darkorange":           {},
	"darkorchid":           {},
	"darkred":              {},
	"darksalmon":           {},
	"darkseagreen":         {},
	"darkslateblue":        {},
	"darkslategray":        {},
	"darkslategrey":        {},
	"darkturquoise":        {},
	"darkviolet":           {},
	"deeppink":             {},
	"deepskyblue":          {},
	"dimgray":              {},
	"dimgrey":              {},
	"dodgerblue":           {},
	"firebrick":            {},
	"floralwhite":          {},
	"forestgreen":          {},
	"fuchsia":              {},
	"gainsboro":            {},
	"ghostwhite":           {},
	"gold":                 {},
	"goldenrod":            {},
	"gray":                 {},
	"green":                {},
	"greenyellow":          {},
	"grey":                 {},
	"honeydew":             {},
	"hotpink":              {},
	"indianred":            {},
	"indigo":               {},
	"ivory":                {},
	"khaki":                {},
	"lavender":             {},
	"lavenderblush":        {},
	"lawngreen":            {},
	"lemonchiffon":         {},
	"lightblue":            {},
	"lightcoral":           {},
	"lightcyan":            {},
	"lightgoldenrodyellow": {},
	"lightgray":            {},
	"lightgreen":           {},
	"lightgrey":            {},
	"lightpink":            {},
	"lightsalmon":          {},
	"lightseagreen":        {},
	"lightskyblue":         {},
	"lightslategray":       {},
	"lightslategrey":       {},
	"lightsteelblue":       {},
	"lightyellow":          {},
	"lime":                 {},
	"limegreen":            {},
	"linen":                {},
	"magenta":              {},
	"maroon":               {},
	"mediumaquamarine":     {},
	"mediumblue":           {},
	"mediumorchid":         {},
	"mediumpurple":         {},
	"mediumseagreen":       {},
	"mediumslateblue":      {},
	"mediumspringgreen":    {},
	"mediumturquoise":      {},
	"mediumvioletred":      {},
	"midnightblue":         {},
	"mintcream":            {},
	"mistyrose":            {},
	"moccasin":             {},
	"navajowhite":          {},
	"navy":                 {},
	"oldlace":              {},
	"olive":                {},
	"olivedrab":            {},
	"orange":               {},
	"orangered":            {},
	"orchid":               {},
	"palegoldenrod":        {},
	"palegreen":            {},
	"paleturquoise":        {},
	"palevioletred":        {},
	"papayawhip":           {},
	"peachpuff":            {},
	"peru":                 {},
	"pink":                 {},
	"plum":                 {},
	"powderblue":           {},
	"purple":               {},
	"red":                  {},
	"rosybrown":            {},
	"royalblue":            {},
	"saddlebrown":          {},
	"salmon":               {},
	"sandybrown":           {},
	"seagreen":             {},
	"seashell":             {},
	"sienna":               {},
	"silver":               {},
	"skyblue":              {},
	"slateblue":            {},
	"slategray":            {},
	"slategrey":            {},
	"snow":                 {},
	"springgreen":          {},
	"steelblue":            {},
	"tan":                  {},
	"teal":                 {},
	"thistle":              {},
	"tomato":               {},
	"turquoise":            {},
	"violet":               {},
	"wheat":                {},
	"white":                {},
	"whitesmoke":           {},
	"yellow":               {},
	"yellowgreen":          {},
}

func isValidHexColor(str string) bool {
	hexColorPattern := regexp.MustCompile(`^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	matched := hexColorPattern.FindStringSubmatch(str)
	return matched != nil
}

func isValidCSSColorName(str string) bool {
	_, ok := cssColorNames[str]
	return ok
}

func parseColor(str string) string {
	lowercaseStr := strings.ToLower(str)

	if isValidHexColor(lowercaseStr) {
		if lowercaseStr[0] != '#' {
			return "#" + lowercaseStr
		}

		return lowercaseStr
	}

	if isValidCSSColorName(lowercaseStr) {
		return lowercaseStr
	}

	return ""
}
