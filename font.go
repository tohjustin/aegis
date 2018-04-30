package main

import (
	"errors"
	"math"

	"github.com/gobuffalo/packr"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

func computeTextWidth(text string, fontSize int, fontFamily string) (int, error) {
	box := packr.NewBox("./assets/fonts")
	fontBinary := box.Bytes(fontFamily + ".ttf")
	f, err := truetype.Parse(fontBinary)
	if err != nil {
		return 0.0, errors.New("computeTextWidth: Unable to parse \"" + fontFamily + ".ttf\"")
	}

	textArray := []rune(text)
	textWidth := fixed.Int26_6(0)
	fUnitsPerEm := fixed.Int26_6(f.FUnitsPerEm())

	for i := range textArray {
		i0 := f.Index(textArray[i])
		horizontalMetric := f.HMetric(fUnitsPerEm, i0)
		textWidth += horizontalMetric.AdvanceWidth

		if i < len(textArray)-1 {
			i1 := f.Index(textArray[i+1])
			kerning := f.Kern(fUnitsPerEm, i0, i1)
			textWidth += kerning
		}
	}

	return int(math.Round(float64(textWidth) / float64(fUnitsPerEm) * float64(fontSize))), nil
}
