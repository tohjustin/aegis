package badge

import (
	"regexp"
)

const defaultColor = "#1bacbf"

var cssColorNames = map[string]struct{}{
	"aliceblue":            struct{}{},
	"antiquewhite":         struct{}{},
	"aqua":                 struct{}{},
	"aquamarine":           struct{}{},
	"azure":                struct{}{},
	"beige":                struct{}{},
	"bisque":               struct{}{},
	"black":                struct{}{},
	"blanchedalmond":       struct{}{},
	"blue":                 struct{}{},
	"blueviolet":           struct{}{},
	"brown":                struct{}{},
	"burlywood":            struct{}{},
	"cadetblue":            struct{}{},
	"chartreuse":           struct{}{},
	"chocolate":            struct{}{},
	"coral":                struct{}{},
	"cornflowerblue":       struct{}{},
	"cornsilk":             struct{}{},
	"crimson":              struct{}{},
	"cyan":                 struct{}{},
	"darkblue":             struct{}{},
	"darkcyan":             struct{}{},
	"darkgoldenrod":        struct{}{},
	"darkgray":             struct{}{},
	"darkgreen":            struct{}{},
	"darkgrey":             struct{}{},
	"darkkhaki":            struct{}{},
	"darkmagenta":          struct{}{},
	"darkolivegreen":       struct{}{},
	"darkorange":           struct{}{},
	"darkorchid":           struct{}{},
	"darkred":              struct{}{},
	"darksalmon":           struct{}{},
	"darkseagreen":         struct{}{},
	"darkslateblue":        struct{}{},
	"darkslategray":        struct{}{},
	"darkslategrey":        struct{}{},
	"darkturquoise":        struct{}{},
	"darkviolet":           struct{}{},
	"deeppink":             struct{}{},
	"deepskyblue":          struct{}{},
	"dimgray":              struct{}{},
	"dimgrey":              struct{}{},
	"dodgerblue":           struct{}{},
	"firebrick":            struct{}{},
	"floralwhite":          struct{}{},
	"forestgreen":          struct{}{},
	"fuchsia":              struct{}{},
	"gainsboro":            struct{}{},
	"ghostwhite":           struct{}{},
	"gold":                 struct{}{},
	"goldenrod":            struct{}{},
	"gray":                 struct{}{},
	"green":                struct{}{},
	"greenyellow":          struct{}{},
	"grey":                 struct{}{},
	"honeydew":             struct{}{},
	"hotpink":              struct{}{},
	"indianred":            struct{}{},
	"indigo":               struct{}{},
	"ivory":                struct{}{},
	"khaki":                struct{}{},
	"lavender":             struct{}{},
	"lavenderblush":        struct{}{},
	"lawngreen":            struct{}{},
	"lemonchiffon":         struct{}{},
	"lightblue":            struct{}{},
	"lightcoral":           struct{}{},
	"lightcyan":            struct{}{},
	"lightgoldenrodyellow": struct{}{},
	"lightgray":            struct{}{},
	"lightgreen":           struct{}{},
	"lightgrey":            struct{}{},
	"lightpink":            struct{}{},
	"lightsalmon":          struct{}{},
	"lightseagreen":        struct{}{},
	"lightskyblue":         struct{}{},
	"lightslategray":       struct{}{},
	"lightslategrey":       struct{}{},
	"lightsteelblue":       struct{}{},
	"lightyellow":          struct{}{},
	"lime":                 struct{}{},
	"limegreen":            struct{}{},
	"linen":                struct{}{},
	"magenta":              struct{}{},
	"maroon":               struct{}{},
	"mediumaquamarine":     struct{}{},
	"mediumblue":           struct{}{},
	"mediumorchid":         struct{}{},
	"mediumpurple":         struct{}{},
	"mediumseagreen":       struct{}{},
	"mediumslateblue":      struct{}{},
	"mediumspringgreen":    struct{}{},
	"mediumturquoise":      struct{}{},
	"mediumvioletred":      struct{}{},
	"midnightblue":         struct{}{},
	"mintcream":            struct{}{},
	"mistyrose":            struct{}{},
	"moccasin":             struct{}{},
	"navajowhite":          struct{}{},
	"navy":                 struct{}{},
	"oldlace":              struct{}{},
	"olive":                struct{}{},
	"olivedrab":            struct{}{},
	"orange":               struct{}{},
	"orangered":            struct{}{},
	"orchid":               struct{}{},
	"palegoldenrod":        struct{}{},
	"palegreen":            struct{}{},
	"paleturquoise":        struct{}{},
	"palevioletred":        struct{}{},
	"papayawhip":           struct{}{},
	"peachpuff":            struct{}{},
	"peru":                 struct{}{},
	"pink":                 struct{}{},
	"plum":                 struct{}{},
	"powderblue":           struct{}{},
	"purple":               struct{}{},
	"red":                  struct{}{},
	"rosybrown":            struct{}{},
	"royalblue":            struct{}{},
	"saddlebrown":          struct{}{},
	"salmon":               struct{}{},
	"sandybrown":           struct{}{},
	"seagreen":             struct{}{},
	"seashell":             struct{}{},
	"sienna":               struct{}{},
	"silver":               struct{}{},
	"skyblue":              struct{}{},
	"slateblue":            struct{}{},
	"slategray":            struct{}{},
	"slategrey":            struct{}{},
	"snow":                 struct{}{},
	"springgreen":          struct{}{},
	"steelblue":            struct{}{},
	"tan":                  struct{}{},
	"teal":                 struct{}{},
	"thistle":              struct{}{},
	"tomato":               struct{}{},
	"turquoise":            struct{}{},
	"violet":               struct{}{},
	"wheat":                struct{}{},
	"white":                struct{}{},
	"whitesmoke":           struct{}{},
	"yellow":               struct{}{},
	"yellowgreen":          struct{}{},
}

func isValidHexColor(str string) bool {
	hexColorPattern := regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
	matched := hexColorPattern.FindStringSubmatch(str)
	return matched != nil
}

func isValidCSSColorName(str string) bool {
	_, ok := cssColorNames[str]
	return ok
}

func parseColor(str string) string {
	if color := "#" + str; isValidHexColor(color) {
		return color
	}

	if isValidCSSColorName(str) {
		return str
	}

	return defaultColor
}
