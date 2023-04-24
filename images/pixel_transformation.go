package images

import "math"

type PixelTransformation interface {
	Transform(r, g, b int) int
}

///////////////////////////////////////////////////////////////////////////////

type PixelTransformationGrayscale struct {
}

func (c *PixelTransformationGrayscale) Transform(r, g, b int) int {
	return int(math.Ceil(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)))
}

///////////////////////////////////////////////////////////////////////////////

const (
	redHueThreshold        = 15
	redSaturationThreshold = 0.75
	redLightensThreshold   = 0.70
)

type PixelTransformationRed struct {
}

func (c *PixelTransformationRed) Transform(r, g, b int) int {
	var h, s, l float64

	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	//https://github.com/gerow/go-color/blob/master/color.go
	max := math.Max(rf, math.Max(gf, bf))
	min := math.Min(rf, math.Min(gf, bf))

	l = (max + min) / 2.0

	d := max - min
	if d == 0 {
		return 255 //gray
	}

	if l < 0.5 {
		s = d / (max + min)
	} else {
		s = d / (2 - max - min)
	}

	r2 := (((max - rf) / 6) + (d / 2)) / d
	g2 := (((max - gf) / 6) + (d / 2)) / d
	b2 := (((max - bf) / 6) + (d / 2)) / d
	switch {
	case rf == max:
		h = b2 - g2
	case gf == max:
		h = (1.0 / 3.0) + r2 - b2
	case bf == max:
		h = (2.0 / 3.0) + g2 - r2
	}

	switch {
	case h < 0:
		h += 1
	case h > 1:
		h -= 1
	}

	h *= 360.0

	if h > redHueThreshold && h < 360-redHueThreshold {
		return 255 //not in the red part of the hue circle
	}
	if s < redSaturationThreshold || l > redLightensThreshold {
		return 255
	}

	return 255 - int(math.Ceil(255*s))
}
