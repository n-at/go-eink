package images

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
)

const (
	SubtractNone  = "none"
	SubtractRed   = "red"
	SubtractBlack = "black"
)

func Open(path string) (image.Image, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func Save(img image.Image, path string) error {
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	if err := png.Encode(writer, img); err != nil {
		return err
	}

	return nil
}

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

func Join(black, red image.Image) image.Image {
	result := image.NewRGBA(black.Bounds())
	width := result.Bounds().Dx()
	height := result.Bounds().Dy()
	draw.Draw(result, result.Bounds(), black, image.Point{X: 0, Y: 0}, draw.Src)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, _, _, _ := red.At(x, y).RGBA()
			if r == 0 {
				result.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			}
		}
	}

	return result
}
