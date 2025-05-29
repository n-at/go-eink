package images

import (
	"math"
)

type PixelTransformation interface {
	Transform(r, g, b int) int
	GetThreshold() int
}

///////////////////////////////////////////////////////////////////////////////

type PixelTransformationGrayscale struct {
	Threshold int
}

func (c *PixelTransformationGrayscale) Transform(r, g, b int) int {
	return int(math.Ceil(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)))
}

func (c *PixelTransformationGrayscale) GetThreshold() int {
	return c.Threshold
}

///////////////////////////////////////////////////////////////////////////////

type PixelTransformationRed struct {
	Threshold              int
	RedHueThreshold        int
	RedSaturationThreshold int
	RedLightnessThreshold  int
}

func (c *PixelTransformationRed) Transform(r, g, b int) int {
	h, s, l := RgbToHsl(r, g, b)
	if h == 0 && s == 0 {
		return 255 //gray
	}

	//red: h = 0 or 360

	if h > float64(c.RedHueThreshold)/2.0 && h < 360-float64(c.RedHueThreshold)/2.0 {
		return 255 //not in the red part of the hue circle
	}
	if s < float64(c.RedSaturationThreshold) {
		return 255
	}
	if l > float64(c.RedLightnessThreshold) {
		return 255
	}

	return int(255 * s / 100.0)
}

func (c *PixelTransformationRed) GetThreshold() int {
	return c.Threshold
}

///////////////////////////////////////////////////////////////////////////////

type PixelTransformationYellow struct {
	Threshold                 int
	YellowHueThreshold        int
	YellowSaturationThreshold int
	YellowLightnessThreshold  int
}

func (c *PixelTransformationYellow) Transform(r, g, b int) int {
	h, s, l := RgbToHsl(r, g, b)
	if h == 0 && s == 0 {
		return 255 //gray
	}

	//yellow: h = 50
	//https://hslpicker.com/#ffff00

	if h > (50.0+float64(c.YellowHueThreshold)/2.0) || h < (50.0-float64(c.YellowHueThreshold)/2.0) {
		return 255
	}
	if s < float64(c.YellowSaturationThreshold) {
		return 255
	}
	if l > float64(c.YellowLightnessThreshold) {
		return 255
	}

	return int(255 * s / 100.0)
}

func (c *PixelTransformationYellow) GetThreshold() int {
	return c.Threshold
}

///////////////////////////////////////////////////////////////////////////////

func RgbToHsl(r, g, b int) (float64, float64, float64) {
	r = min(255, max(0, r))
	g = min(255, max(0, g))
	b = min(255, max(0, b))

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
		return 0, 0, l * 100.0 //gray
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
	return h * 360.0, s * 100.0, l * 100.0
}
