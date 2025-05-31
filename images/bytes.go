package images

import "image"

const (
	BWRY_B = 0b00
	BWRY_W = 0b01
	BWRY_R = 0b11
	BWRY_Y = 0b10
)

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
	width := imgBW.Bounds().Dx()
	height := imgBW.Bounds().Dy()

	outputLength := (width * height) / 4
	output := make([]byte, outputLength)
	outputIdx := 0

	for y := 0; y < height; y++ {
		for xStart := 0; xStart < width; xStart += 4 {
			var val byte = 0
			for xPos := 0; xPos < 4; xPos++ {
				val <<= 2

				x := xStart + xPos

				black := false
				red := false
				yellow := false

				if v, _, _, _ := imgBW.At(x, y).RGBA(); v == 0 {
					black = true
				}
				if v, _, _, _ := imgRW.At(x, y).RGBA(); v == 0 {
					red = true
				}
				if v, _, _, _ := imgYW.At(x, y).RGBA(); v == 0 {
					yellow = true
				}

				c := BlendColors(blendMode, black, red, yellow)

				switch c {
				case BlendModeB:
					val += BWRY_B
				case BlendModeR:
					val += BWRY_R
				case BlendModeY:
					val += BWRY_Y
				default:
					val += BWRY_W
				}
			}

			output[outputIdx] = val
			outputIdx++
		}
	}

	return output
}
