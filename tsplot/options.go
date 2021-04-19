package tsplot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"image/color"
)

type PlotOption func(p *plot.Plot)

func WithForeground(color color.Color) PlotOption {
	return func(p *plot.Plot) {
		p.Title.TextStyle.Color = color
		p.X.Color = color
		p.X.Label.TextStyle.Color = color
		p.X.Tick.Color = color
		p.X.Tick.Label.Color = color
		p.Y.Tick.Label.Color = color
		p.Y.Color = color
		p.Y.Label.TextStyle.Color = color
		p.Y.Tick.Color = color
		p.Y.Tick.Label.Color = color
	}
}

func WithBackground(color color.Color) PlotOption {
	return func(p *plot.Plot) {
		p.BackgroundColor = color
	}
}

func WithTitle(title string) PlotOption {
	return func(p *plot.Plot) {
		p.Title.Text = title
	}
}

func WithGrid(color color.Color) PlotOption {
	return func(p *plot.Plot) {
		grid := plotter.NewGrid()
		grid.Horizontal.Color = color
		grid.Vertical.Color = color
		p.Add(grid)
	}
}

func WithLegend(color color.Color) PlotOption {
	return func(p *plot.Plot) {
		p.Legend.TextStyle.Color = color
	}
}

func WithXAxisName(name string) PlotOption {
	return func(p *plot.Plot) {
		p.X.Label.Text = name
	}
}

func WithYAxisName(name string) PlotOption {
	return func(p *plot.Plot) {
		p.Y.Label.Text = name
	}
}

func WithXTimeTicks(format string) PlotOption {
	return func(p *plot.Plot) {
		p.X.Tick.Marker = plot.TimeTicks{
			Format: format,
		}
	}
}

func ApplyDefaultHighContrast(p *plot.Plot) {
	opts := []PlotOption{
		WithBackground(DefaultColors_HighContrast.Background),
		WithForeground(DefaultColors_HighContrast.Foreground),
		WithLegend(DefaultColors_HighContrast.Foreground),
	}
	for _, opt := range opts {
		opt(p)
	}
}
