package images

import (
	"image"
	"image/color"
)

type AlignValue int

const (
	AlignTopLeft = iota
	AlignTopMiddle
	AlignTopRight
	AlignMiddleLeft
	AlignMiddle
	AlignMiddleRight
	AlignBottomLeft
	AlignBottomMiddle
	AlignBottomRight
)

func Align(img image.Image, imageWidth, imageHeight int, align AlignValue) image.Image {
	size := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: imageWidth, Y: imageHeight},
	}
	newImage := image.NewRGBA(size)

	for x := 0; x < imageWidth; x++ {
		for y := 0; y < imageHeight; y++ {
			newImage.Set(x, y, color.White)
		}
	}

	outputImageWidth := img.Bounds().Size().X
	outputImageHeight := img.Bounds().Size().X

	outputMiddleX := max(imageWidth/2-outputImageWidth/2, 0)
	outputMiddleY := max(imageHeight/2-outputImageHeight/2, 0)

	outputMaxX := max(imageWidth-outputImageWidth, 0)
	outputMaxY := max(imageHeight-outputImageHeight, 0)

	offsetX := 0
	offsetY := 0

	switch align {
	case AlignTopLeft:
		offsetX = 0
		offsetY = 0

	case AlignTopMiddle:
		offsetX = outputMiddleX
		offsetY = 0

	case AlignTopRight:
		offsetX = outputMaxX
		offsetY = 0

	case AlignMiddleLeft:
		offsetX = 0
		offsetY = outputMiddleY

	case AlignMiddle:
		offsetX = outputMiddleX
		offsetY = outputMiddleY

	case AlignMiddleRight:
		offsetX = outputMaxX
		offsetY = outputMiddleY

	case AlignBottomLeft:
		offsetX = 0
		offsetY = outputMaxY

	case AlignBottomMiddle:
		offsetX = outputMiddleX
		offsetY = outputMaxY

	case AlignBottomRight:
		offsetX = outputMaxX
		offsetY = outputMaxY
	}

	for x := 0; x < outputImageWidth; x++ {
		for y := 0; y < outputImageHeight; y++ {
			newX := offsetX + x
			if newX >= imageWidth {
				continue
			}

			newY := offsetY + y
			if newY >= imageHeight {
				continue
			}

			newImage.Set(newX, newY, img.At(x, y))
		}
	}

	return newImage
}

func GetAlign(name string) AlignValue {
	switch name {
	case "top-left":
		return AlignTopLeft
	case "top-middle":
		return AlignTopMiddle
	case "top-right":
		return AlignTopRight
	case "middle-left":
		return AlignMiddleLeft
	case "middle":
		return AlignMiddle
	case "middle-right":
		return AlignMiddleRight
	case "bottom-left":
		return AlignBottomLeft
	case "bottom-middle":
		return AlignBottomMiddle
	case "bottom-right":
		return AlignBottomRight
	default:
		return AlignMiddle
	}
}
