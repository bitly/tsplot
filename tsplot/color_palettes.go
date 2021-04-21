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
			color.RGBA{
				R: 0xe6,
				G: 0x9f,
				B: 0x00,
				A: 0xff,
			},
			color.RGBA{
				R: 0x56,
				G: 0xb4,
				B: 0xe9,
				A: 0xff,
			},
			color.RGBA{
				R: 0x00,
				G: 0x9e,
				B: 0x73,
				A: 0xff,
			},
			color.RGBA{
				R: 0xf0,
				G: 0xe4,
				B: 0x42,
				A: 0xff,
			},
			color.RGBA{
				R: 0x00,
				G: 0x72,
				B: 0xb2,
				A: 0x00,
			},
			color.RGBA{
				R: 0xd5,
				G: 0x5e,
				B: 0x00,
				A: 0xff,
			},
			color.RGBA{
				R: 0xcc,
				G: 0x79,
				B: 0xa7,
				A: 0xff,
			},
		},
	}
)
