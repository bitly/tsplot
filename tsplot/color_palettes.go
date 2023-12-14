package tsplot

import (
	"image/color"

	"golang.org/x/image/colornames"
)

var usedColors = make(map[string]color.RGBA)

// simple colors (subset of golang.org/x/image/colornames)
var availableColors = map[string]color.RGBA{
	"blue":    color.RGBA{0x00, 0x00, 0xff, 0xff}, // rgb(0, 0, 255)
	"brown":   color.RGBA{0xa5, 0x2a, 0x2a, 0xff}, // rgb(165, 42, 42)
	"orange":  color.RGBA{0xff, 0xa5, 0x00, 0xff}, // rgb(255, 165, 0)
	"hotpink": color.RGBA{0xff, 0x69, 0xb4, 0xff}, // rgb(255, 105, 180)
	"red":     color.RGBA{0xff, 0x00, 0x00, 0xff}, // rgb(255, 0, 0)
	"purple":  color.RGBA{0x80, 0x00, 0x80, 0xff}, // rgb(128, 0, 128)
	"yellow":  color.RGBA{0xff, 0xff, 0x00, 0xff}, // rgb(255, 255, 0)
	"green":   color.RGBA{0x00, 0x80, 0x00, 0xff}, // rgb(0, 128, 0)

}

func getUnusedColor() color.RGBA {
	for k, v := range availableColors {
		if _, ok := usedColors[k]; !ok {
			usedColors[k] = v
			return v
		}
	}
	return colornames.Black
}

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
