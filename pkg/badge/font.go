package badge

import "fmt"

const fallbackCharCode = 64 // @

func computeTextWidth(text string, fontSize int, fontFamily string) (int, error) {
	textArray := []rune(text)
	textWidth := 0

	var charWidthTable []int
	switch fontFamily {
	case "Verdana":
		switch fontSize {
		case 9:
			charWidthTable = verdana9CharWidths[:]
		case 11:
			charWidthTable = verdana11CharWidths[:]
		default:
			return 0, fmt.Errorf("unsupported font size: %d", fontSize)
		}
	default:
		return 0, fmt.Errorf("unsupported font family: %s", fontFamily)
	}

	charWidthTableSize := len(charWidthTable)
	for _, character := range textArray {
		charCode := int(character)
		if charCode > charWidthTableSize {
			charCode = fallbackCharCode
		}
		textWidth += charWidthTable[charCode]
	}

	return textWidth, nil
}
