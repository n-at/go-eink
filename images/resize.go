package images

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"math"
)

func Resize(img image.Image, newWidth, newHeight int, enlarge bool) image.Image {
	originalBounds := img.Bounds()
	width := originalBounds.Size().X
	height := originalBounds.Size().Y

	widthScale := float64(newWidth) / float64(width)
	heightScale := float64(newHeight) / float64(height)

	scale := math.Min(widthScale, heightScale)

	if scale > 1.0 && !enlarge {
		return img
	}

	scaledWidth := int(math.Floor(float64(width) * scale))
	scaledHeight := int(math.Floor(float64(height) * scale))

	scaledImage := resize.Resize(uint(scaledWidth), uint(scaledHeight), img, resize.Lanczos3)

	plainImage := image.NewRGBA(scaledImage.Bounds())
	for x := 0; x < plainImage.Bounds().Size().X; x++ {
		for y := 0; y < plainImage.Bounds().Size().Y; y++ {
			r, g, b, a := scaledImage.At(x, y).RGBA()
			plainImage.Set(x, y, color.RGBA{
				R: blendAlpha(r, a),
				G: blendAlpha(g, a),
				B: blendAlpha(b, a),
				A: 255,
			})
		}
	}

	return plainImage
}
