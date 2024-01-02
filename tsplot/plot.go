package tsplot

import (
	"fmt"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/api/metric"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

// findMaxFromPoints finds the maximum value of a list of metric points of the GAPI MetricDescriptor Value type.
func findMaxFromPoints(gapiValueType metric.MetricDescriptor_ValueType, points []*monitoringpb.Point) (float64, error) {
	var maximum float64

	switch gapiValueType {
	case metric.MetricDescriptor_INT64:
		maximum = findMaxFromInt64Data(points)
	case metric.MetricDescriptor_DOUBLE:
		maximum = findMaxFromFloat64Data(points)
	default:
		return 0, fmt.Errorf("cannot operate on value type %q", gapiValueType.String())
	}

	return maximum, nil

}

// findMaxFromFloat64Data is a helper function to findMaxFromPoints but is specific to the data being of type float64.
// The other helper function findMaxFromInt64Data is also required due to the methods that which Google allows the consumer
// to fetch the data.
func findMaxFromFloat64Data(points []*monitoringpb.Point) float64 {
	var maximum float64
	for _, v := range points {
		cur := v.GetValue().GetDoubleValue()
		if cur > maximum {
			maximum = cur
		}
	}
	return maximum
}

// findMaxFromInt64Data is a helper function to findMaxFromPoints but is specific to the data being of type float64.
// The other helper function findMaxFromFloat64Data is also required due to the methods that which Google allows the consumer
// to fetch the data.
func findMaxFromInt64Data(points []*monitoringpb.Point) float64 {
	var maximum int64
	for _, v := range points {
		cur := v.GetValue().GetInt64Value()
		if cur > maximum {
			maximum = cur
		}
	}
	return float64(maximum)
}

// NewPlotFromTimeSeries creates a plot from a single time series.
func NewPlotFromTimeSeries(ts *monitoringpb.TimeSeries, opts ...PlotOption) (*plot.Plot, error) {
	points := ts.GetPoints()
	YMax, err := findMaxFromPoints(ts.GetValueType(), points)
	if err != nil {
		return nil, err
	}

	opts = append(opts, WithLineFromPoints(points))
	p := plot.New()
	for _, opt := range opts {
		opt(p)
	}

	p.Y.Max = YMax + (.1 * YMax)
	return p, nil
}

// NewPlotFromTimeSeriesIterator creates a plot from multiple time series.
func NewPlotFromTimeSeriesIterator(tsi *monitoring.TimeSeriesIterator, legendKey string, opts ...PlotOption) (*plot.Plot, error) {

	var yMax float64
	p := plot.New()
	for {
		timeSeries, err := tsi.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, err
		}

		points := timeSeries.GetPoints()

		// Find the maximum datapoint between the different time series data.
		// Use this to scale the Y Axis.
		curMax, _ := findMaxFromPoints(timeSeries.GetValueType(), points)
		if curMax > yMax {
			yMax = curMax
		}

		// add colored line to plot
		lineColor := getUnusedColor()
		applyLine := WithColoredLineFromPoints(timeSeries.GetPoints(), lineColor)
		applyLine(p)

		// add to legend
		if legendKey != "" {
			legendEntry, _ := plotter.NewPolygon()
			legendEntry.Color = lineColor
			p.Legend.Left = true
			p.Legend.Top = true
			p.Legend.Add(timeSeries.GetMetric().GetLabels()[legendKey], legendEntry)
		}
	}

	// set Y Axis scale
	p.Y.Max = yMax + (.1 * yMax)

	for _, opt := range opts {
		opt(p)
	}

	resetUsedColors()

	return p, nil
}

// createLine creates a line from data points.
func createLine(dataPoints []*monitoringpb.Point) (*plotter.Line, error) {
	var XYs plotter.XYs
	for _, point := range dataPoints {
		x := point.GetInterval().GetEndTime().GetSeconds()
		y := point.GetValue().GetDoubleValue() // todo: This breaks if the value type is an int64
		XYs = append(XYs, plotter.XY{
			X: float64(x),
			Y: y,
		})
	}

	return plotter.NewLine(XYs)
}
