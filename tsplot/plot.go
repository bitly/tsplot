package tsplot

import (
	"errors"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// TimeSeries is a map representation of unique names to time series.
type TimeSeries map[string][]*monitoringpb.Point

// Plot builds and returns a *plot.Plot. If the TimeSeriesGroup field
// in TimeSeriesPlot is empty, a nil plot is returned as well as an error.
// This is also the case if an error is encountered building the line from the XY coordinates.
func (ts TimeSeries) Plot(opts ...PlotOption) (*plot.Plot, error) {

	if len(ts) == 0 {
		return nil, errors.New("no data to plot")
	}

	p := plot.New()

	// default high contrast
	ApplyDefaultHighContrast(p)

	// user overrides
	for _, opt := range opts {
		opt(p)
	}

	lineColors := DefaultColors_HighContrast.LineColors

	// create a unique line for each time series
	// each line should have a unique color and entry
	// in the legend.
	// current limit = 4
	limit := len(lineColors) - 1
	for name, series := range ts {
		if limit < 0 {
			break
		}
		line, err := createLine(series)
		if err != nil {
			return nil, err
		}

		// width of the line
		line.Width = vg.Points(2)

		// color the line
		line.Color = lineColors[limit]

		// add to legend
		p.Legend.Add(name, line)

		// add to chart
		p.Add(line)
		limit--
	}

	return p, nil
}

// createLine creates a line from data points.
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
