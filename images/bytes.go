package images

import (
	"image"
	"image/color"
)

const (
	BWRY_B = 0b00
	BWRY_W = 0b01
	BWRY_R = 0b11
	BWRY_Y = 0b10
)

///////////////////////////////////////////////////////////////////////////////

func ToImageDataBW(img image.Image) []byte {
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	outputLength := (width * height) / 8
	output := make([]byte, outputLength)
	for i := range outputLength {
		output[i] = 0
	}

	current := 0
	bitNum := 0

	for y := range height {
		for x := range width {
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

	return prepareImageDataBW(output)
}

func prepareImageDataBW(imageData []byte) []byte {
	prepared := make([]byte, len(imageData))

	for idx := range imageData {
		if imageData[idx] == 13 {
			prepared[idx] = 12
		} else {
			prepared[idx] = imageData[idx]
		}
	}

	return prepared
}

///////////////////////////////////////////////////////////////////////////////

func ToImageDataBWR(blendMode BlendMode, imgBW, imgRW image.Image) []byte {
	resultBW := image.NewRGBA(imgBW.Bounds())
	resultRW := image.NewRGBA(imgRW.Bounds())
	width := resultBW.Bounds().Dx()
	height := resultBW.Bounds().Dy()

	colorBlack := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	colorWhite := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	for y := range height {
		for x := range width {
			black := false
			red := false

			if v, _, _, _ := imgBW.At(x, y).RGBA(); v == 0 {
				black = true
			}
			if v, _, _, _ := imgRW.At(x, y).RGBA(); v == 0 {
				red = true
			}

			c := BlendColors(blendMode, black, red, false)

			switch c {
			case BlendModeB:
				resultBW.Set(x, y, colorBlack)
				resultRW.Set(x, y, colorWhite)
			case BlendModeR:
				resultBW.Set(x, y, colorWhite)
				resultRW.Set(x, y, colorBlack)
			default:
				resultBW.SetRGBA(x, y, colorWhite)
				resultRW.SetRGBA(x, y, colorWhite)
			}
		}
	}

	bitsBW := ToImageDataBW(resultBW)
	bitsRW := ToImageDataBW(resultRW)

	return prepareImageDataBWR(bitsBW, bitsRW)
}

func prepareImageDataBWR(imageDataBW, imageDataRW []byte) []byte {
	prepared := make([]byte, len(imageDataBW)+len(imageDataRW))
	offset := 0

	for idx := range imageDataBW {
		if imageDataBW[idx] == 13 {
			prepared[idx+offset] = 12
		} else {
			prepared[idx+offset] = imageDataBW[idx]
		}
	}

	offset += len(imageDataBW)

	for idx := range imageDataRW {
		if imageDataRW[idx] == 13 {
			prepared[idx+offset] = 12
		} else {
			prepared[idx+offset] = imageDataRW[idx]
		}
	}

	return prepared
}

///////////////////////////////////////////////////////////////////////////////

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

	for i := range output {
		if output[i] == 13 {
			output[i] = 12
		}
	}

	return output
}
