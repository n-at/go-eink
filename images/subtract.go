package images

import (
	"image"
	"image/color"
	"image/draw"
)

const (
	SubtractNone  = "none"
	SubtractRed   = "red"
	SubtractBlack = "black"
)

func Subtract(original, sub image.Image) image.Image {
	result := image.NewRGBA(original.Bounds())
	width := result.Bounds().Dx()
	height := result.Bounds().Dy()
	draw.Draw(result, result.Bounds(), original, image.Point{X: 0, Y: 0}, draw.Src)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			originalC := result.RGBAAt(x, y)
			subR, _, _, _ := sub.At(x, y).RGBA()
			if originalC.R == 0 && subR == 0 {
				result.Set(x, y, color.White)
			}
		}
	}

	return result
}
