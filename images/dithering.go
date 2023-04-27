package images

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

type DitheringMultipliers [][]float64

func Dithering(img image.Image, transformation PixelTransformation, multipliers DitheringMultipliers) image.Image {
	result := image.NewRGBA(img.Bounds())
	width := result.Bounds().Dx()
	height := result.Bounds().Dy()
	draw.Draw(result, result.Bounds(), img, image.Point{X: 0, Y: 0}, draw.Src)

	var errors [][][]float64
	for x := 0; x < width; x++ {
		var values [][]float64
		for y := 0; y < height; y++ {
			values = append(values, []float64{0, 0, 0})
		}
		errors = append(errors, values)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := result.RGBAAt(x, y)

			r := max(0, min(255, int(math.Ceil(float64(c.R)+errors[x][y][0]))))
			g := max(0, min(255, int(math.Ceil(float64(c.G)+errors[x][y][1]))))
			b := max(0, min(255, int(math.Ceil(float64(c.B)+errors[x][y][2]))))
			gray := transformation.Transform(r, g, b)

			transformedColor := 0.0
			if gray < transformation.GetThreshold() {
				result.Set(x, y, color.Black)
			} else {
				result.Set(x, y, color.White)
				transformedColor = 255.0
			}

			redError := float64(gray) - transformedColor
			greenError := float64(gray) - transformedColor
			blueError := float64(gray) - transformedColor

			for ky := 0; ky < len(multipliers); ky++ {
				for kx := 0; kx < len(multipliers[ky]); kx++ {
					nx := x + kx - 2
					if nx < 0 || nx >= width {
						continue
					}

					ny := y + ky
					if ny < 0 || ny >= height {
						continue
					}

					errors[nx][ny][0] += redError * multipliers[ky][kx]
					errors[nx][ny][1] += greenError * multipliers[ky][kx]
					errors[nx][ny][2] += blueError * multipliers[ky][kx]
				}
			}
		}
	}

	return result
}

///////////////////////////////////////////////////////////////////////////////
//multipliers

var DitheringFloydSteinberg = [][]float64{
	{0.0, 0.0, 0.0, 7.0 / 16.0, 0.0},
	{0.0, 0.0, 5.0 / 16.0, 1.0 / 16.0, 0.0},
	{0.0, 0.0, 0.0, 0.0, 0.0},
}

var DitheringJarvisJudiceNinke = [][]float64{
	{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
	{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
	{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
}

var DitheringStucki = [][]float64{
	{0.0, 0.0, 0.0, 8.0 / 42.0, 4.0 / 42.0},
	{2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0},
	{1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0},
}

var DitheringAtkinson = [][]float64{
	{0.0, 0.0, 0.0, 1.0 / 8.0, 1.0 / 8.0},
	{0.0, 1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0.0},
	{0.0, 0.0, 1.0 / 8.0, 0.0, 0.0},
}

var DitheringBurkes = [][]float64{
	{0.0, 0.0, 0.0, 8.0 / 32.0, 4.0 / 32.0},
	{2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0},
	{0.0, 0.0, 0.0, 0.0, 0.0},
}

var DitheringSierra = [][]float64{
	{0.0, 0.0, 0.0, 5.0 / 32.0, 3.0 / 32.0},
	{2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0},
	{0.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0.0},
}

func GetDitheringAlgorithm(name string) DitheringMultipliers {
	switch name {
	case "floyd_steinberg":
		return DitheringFloydSteinberg
	case "jarvis_judice_ninke":
		return DitheringJarvisJudiceNinke
	case "atkinson":
		return DitheringAtkinson
	case "burkes":
		return DitheringBurkes
	case "stucki":
		return DitheringStucki
	case "sierra":
		return DitheringSierra
	default:
		return DitheringStucki
	}
}
