package tsplot

import (
	"errors"
	"image/color"
	"time"

	"golang.org/x/image/colornames"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

type TimeSeries map[string][]*monitoringpb.Point

type TimeSeriesPlot struct {
	Name        string
	XAxisName   string
	YAxisName   string
	Description string
	TimeSeries  TimeSeries
}

// color palette
// high contrast; attempt to be colorblind friendly
var (
	// graph colors
	backgroundColor = colornames.Black    // background color of plot
	foregroundColor = colornames.White    // color used for foreground components like the axis and labels
	gridColor       = colornames.Darkgrey // grid overlay color

	// line colors - bright dissimilar colors
	lineColor_1 = colornames.Lightblue
	lineColor_2 = colornames.Yellow
	lineColor_3 = colornames.Orange
	lineColor_4 = colornames.Fuchsia

	palette = []color.RGBA{lineColor_1, lineColor_2, lineColor_3, lineColor_4}
)

// Create builds and returns a *plot.Plot. If the DataPoints field
// in TimeSeriesPlot is nil or empty, a nil plot is returned as well as an error.
// This is also the case if an error is encountered building the line from the XY coordinates.
func (tsp TimeSeriesPlot) Create() (*plot.Plot, error) {

	if len(tsp.TimeSeries) == 0 {
		return nil, errors.New("no data to plot")
	}

	p := plot.New()

	// aesthetics - high contrast
	p.Title.Text = tsp.Name
	p.Title.TextStyle.Color = foregroundColor
	p.BackgroundColor = backgroundColor
	grid := plotter.NewGrid()
	grid.Horizontal.Color = gridColor
	grid.Vertical.Color = gridColor
	p.Legend.TextStyle.Color = foregroundColor
	p.X.Color = foregroundColor
	p.Y.Color = foregroundColor
	p.X.Label.TextStyle.Color = foregroundColor
	p.Y.Label.TextStyle.Color = foregroundColor
	p.X.Tick.Color = foregroundColor
	p.Y.Tick.Color = foregroundColor
	p.X.Tick.Label.Color = foregroundColor
	p.Y.Tick.Label.Color = foregroundColor
	p.Add(grid)

	// time format in UTC
	xticks := plot.TimeTicks{
		Format: time.Kitchen,
	}

	p.X.Tick.Marker = xticks

	// if XAxisName is unset, default to UTC
	if tsp.XAxisName == "" {
		p.X.Label.Text = "UTC"
	}

	// create a unique line for each time series
	// each line should have a unique color and entry
	// in the legend.
	// current limit = 4
	limit := len(palette)-1
	for name, series := range tsp.TimeSeries {
		if limit < 0 {
			break
		}
		line, err := createLine(series)
		if err != nil {
			return nil, err
		}

		// color the line
		line.Color = palette[limit]

		// add to legend
		p.Legend.Add(name, line)

		// add to chart
		p.Add(line)
		limit--
	}

	return p, nil
}

func createLine(dataPoints []*monitoringpb.Point) (*plotter.Line, error) {
	var XYs plotter.XYs
	for _, point := range dataPoints {
		x := point.GetInterval().GetEndTime().GetSeconds()
		y := point.GetValue().GetDoubleValue()
		XYs = append(XYs, plotter.XY{
			X: float64(x),
			Y: y,
		})
	}

	return plotter.NewLine(XYs)
}
