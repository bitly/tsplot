package tsplot

import (
	"image/color"

	"golang.org/x/image/colornames"
)

type ColorPalette struct {
	Foreground color.Color
	Background color.Color
	GridColor  color.Color
	LineColors []color.Color
}

var (
	DefaultColors_HighContrast = ColorPalette{
		Foreground: colornames.White,
		Background: colornames.Black,
		GridColor:  colornames.Darkgrey,
		LineColors: []color.Color{
			colornames.Lightblue,
			colornames.Yellow,
			colornames.Orange,
			colornames.Fuchsia,
		},
	}
)
