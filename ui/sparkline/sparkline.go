package sparkline

import (
	"fmt"
	"strings"
)

// Props configures the sparkline SVG.
type Props struct {
	Data        []float64 // data points to plot
	Width       int       // SVG width in pixels (default 60)
	Height      int       // SVG height in pixels (default 24)
	Color       string    // stroke/dot color (default "currentColor")
	StrokeWidth float64   // line width (default 1.5)
	ShowDot     bool      // show dot at last data point (default true)
	DotRadius   float64   // radius of end dot (default 2.5)
	Class       string
}

// sparklineData holds pre-computed SVG values for the template.
type sparklineData struct {
	Points      string
	LastDotCX   string
	LastDotCY   string
	Width       string
	Height      string
	Color       string
	StrokeWidth string
	DotRadius   string
	ShowDot     bool
	Class       string
}

func computeData(p Props) sparklineData {
	if p.Width <= 0 {
		p.Width = 60
	}
	if p.Height <= 0 {
		p.Height = 24
	}
	if p.Color == "" {
		p.Color = "currentColor"
	}
	if p.StrokeWidth <= 0 {
		p.StrokeWidth = 1.5
	}
	if p.DotRadius <= 0 {
		p.DotRadius = 2.5
	}

	d := sparklineData{
		Width:       fmt.Sprintf("%d", p.Width),
		Height:      fmt.Sprintf("%d", p.Height),
		Color:       p.Color,
		StrokeWidth: fmt.Sprintf("%.1f", p.StrokeWidth),
		DotRadius:   fmt.Sprintf("%.1f", p.DotRadius),
		ShowDot:     p.ShowDot,
		Class:       p.Class,
	}

	if len(p.Data) < 2 {
		return d
	}

	min, max := p.Data[0], p.Data[0]
	for _, v := range p.Data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	rng := max - min
	if rng == 0 {
		rng = 1
	}

	parts := make([]string, len(p.Data))
	var lastX, lastY float64
	for i, v := range p.Data {
		x := float64(i) / float64(len(p.Data)-1) * float64(p.Width)
		y := float64(p.Height) - (v-min)/rng*float64(p.Height)
		parts[i] = fmt.Sprintf("%.1f,%.1f", x, y)
		lastX = x
		lastY = y
	}

	d.Points = strings.Join(parts, " ")
	d.LastDotCX = fmt.Sprintf("%.1f", lastX)
	d.LastDotCY = fmt.Sprintf("%.1f", lastY)

	return d
}
