package images

import "image"

func ToImageData(img image.Image) []byte {
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	outputLength := (width * height) / 8
	output := make([]byte, outputLength)
	for i := 0; i < outputLength; i++ {
		output[i] = 0
	}

	current := 0
	bitNum := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			if r/257 > 127 {
				output[current] |= 1
			}

			bitNum++

			if bitNum < 8 {
				output[current] <<= 1
			} else {
				current++
				bitNum = 0
			}
		}
	}

	return output
}

func ToImageDataBWRY(blendMode BlendMode, imgBW, imgRW, imgYW image.Image) []byte {
	return nil //TODO
}
