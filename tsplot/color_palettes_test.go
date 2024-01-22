package tsplot

import (
	"golang.org/x/image/colornames"
	"image/color"
	"reflect"
	"testing"
)

func Test_getUnusedColor(t *testing.T) {
	tests := []struct {
		name    string
		setUsed func()
		want    color.RGBA
	}{
		{
			name: "basic success, unused color is navy",
			setUsed: func() {
				usedColors = map[string]color.RGBA{
					"brown":           color.RGBA{0xa5, 0x2a, 0x2a, 0xff},
					"crimson":         color.RGBA{0xdc, 0x14, 0x3c, 0xff},
					"darkkhaki":       color.RGBA{0xbd, 0xb7, 0x6b, 0xff},
					"deepskyblue":     color.RGBA{0x00, 0xbf, 0xff, 0xff},
					"goldenrod":       color.RGBA{0xda, 0xa5, 0x20, 0xff},
					"gray":            color.RGBA{0x80, 0x80, 0x80, 0xff},
					"green":           color.RGBA{0x00, 0x80, 0x00, 0xff},
					"limegreen":       color.RGBA{0x32, 0xcd, 0x32, 0xff},
					"magenta":         color.RGBA{0xff, 0x00, 0xff, 0xff},
					"mediumturquoise": color.RGBA{0x48, 0xd1, 0xcc, 0xff},
					"orangered":       color.RGBA{0xff, 0x45, 0x00, 0xff},
					"purple":          color.RGBA{0x80, 0x00, 0x80, 0xff},
					"royalblue":       color.RGBA{0x41, 0x69, 0xe1, 0xff},
					"violet":          color.RGBA{0xee, 0x82, 0xee, 0xff},
				}
			},
			want: availableColors["navy"],
		},
		{
			name: "no unused color available, default line color to black",
			setUsed: func() {
				usedColors = map[string]color.RGBA{
					"navy":            color.RGBA{0x00, 0x00, 0x80, 0xff},
					"brown":           color.RGBA{0xa5, 0x2a, 0x2a, 0xff},
					"crimson":         color.RGBA{0xdc, 0x14, 0x3c, 0xff},
					"darkkhaki":       color.RGBA{0xbd, 0xb7, 0x6b, 0xff},
					"deepskyblue":     color.RGBA{0x00, 0xbf, 0xff, 0xff},
					"goldenrod":       color.RGBA{0xda, 0xa5, 0x20, 0xff},
					"gray":            color.RGBA{0x80, 0x80, 0x80, 0xff},
					"green":           color.RGBA{0x00, 0x80, 0x00, 0xff},
					"limegreen":       color.RGBA{0x32, 0xcd, 0x32, 0xff},
					"magenta":         color.RGBA{0xff, 0x00, 0xff, 0xff},
					"mediumturquoise": color.RGBA{0x48, 0xd1, 0xcc, 0xff},
					"orangered":       color.RGBA{0xff, 0x45, 0x00, 0xff},
					"purple":          color.RGBA{0x80, 0x00, 0x80, 0xff},
					"royalblue":       color.RGBA{0x41, 0x69, 0xe1, 0xff},
					"violet":          color.RGBA{0xee, 0x82, 0xee, 0xff},
				}
			},
			want: colornames.Black,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setUsed()
			if got := getUnusedColor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUnusedColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
