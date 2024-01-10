package tsplot

import (
	"image/color"

	"golang.org/x/image/colornames"
)

var usedColors = make(map[string]color.RGBA)

// simple colors (subset of golang.org/x/image/colornames)
var availableColors = map[string]color.RGBA{
	"navy":            color.RGBA{0x00, 0x00, 0x80, 0xff}, // rgb(0, 0, 128)
	"brown":           color.RGBA{0xa5, 0x2a, 0x2a, 0xff}, // rgb(165, 42, 42)
	"crimson":         color.RGBA{0xdc, 0x14, 0x3c, 0xff}, // rgb(220, 20, 60)
	"darkkhaki":       color.RGBA{0xbd, 0xb7, 0x6b, 0xff}, // rgb(189, 183, 107)
	"deepskyblue":     color.RGBA{0x00, 0xbf, 0xff, 0xff}, //(rgb: 0, 191, 255),
	"goldenrod":       color.RGBA{0xda, 0xa5, 0x20, 0xff}, // rgb(218, 165, 32)
	"gray":            color.RGBA{0x80, 0x80, 0x80, 0xff}, // rgb(128, 128, 128)
	"green":           color.RGBA{0x00, 0x80, 0x00, 0xff}, //(rgb: 0, 128, 0),
	"limegreen":       color.RGBA{0x32, 0xcd, 0x32, 0xff}, // rgb(50, 205, 50)
	"magenta":         color.RGBA{0xff, 0x00, 0xff, 0xff}, //(rgb: 255, 0, 255),
	"mediumturquoise": color.RGBA{0x48, 0xd1, 0xcc, 0xff}, // rgb(72, 209, 204)
	"orangered":       color.RGBA{0xff, 0x45, 0x00, 0xff}, // rgb(255, 69, 0)
	"purple":          color.RGBA{0x80, 0x00, 0x80, 0xff}, // rgb(128, 0, 128)
	"royalblue":       color.RGBA{0x41, 0x69, 0xe1, 0xff}, // rgb(65, 105, 225)
	"violet":          color.RGBA{0xee, 0x82, 0xee, 0xff}, //(rgb: 238, 130, 238),
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

func resetUsedColors() {
	usedColors = make(map[string]color.RGBA)
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
