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

var (
	colorWhite  = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	colorBlack  = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	colorRed    = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	colorYellow = color.RGBA{R: 255, G: 255, B: 0, A: 255}
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

///////////////////////////////////////////////////////////////////////////////

func JoinBWR(mode BlendMode, bw, rw image.Image) image.Image {
	result := image.NewRGBA(bw.Bounds())
	width := result.Bounds().Dx()
	height := result.Bounds().Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			black := false
			red := false

			if v, _, _, _ := bw.At(x, y).RGBA(); v == 0 {
				black = true
			}
			if v, _, _, _ := rw.At(x, y).RGBA(); v == 0 {
				red = true
			}

			c := BlendColors(mode, black, red, false)

			switch c {
			case BlendModeB:
				result.Set(x, y, colorBlack)
			case BlendModeR:
				result.Set(x, y, colorRed)
			default:
				result.Set(x, y, colorWhite)
			}
		}
	}

	return result
}

func JoinBWRY(mode BlendMode, bw, rw, yw image.Image) image.Image {
	result := image.NewRGBA(bw.Bounds())
	width := result.Bounds().Dx()
	height := result.Bounds().Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			black := false
			red := false
			yellow := false

			if v, _, _, _ := bw.At(x, y).RGBA(); v == 0 {
				black = true
			}
			if v, _, _, _ := rw.At(x, y).RGBA(); v == 0 {
				red = true
			}
			if v, _, _, _ := yw.At(x, y).RGBA(); v == 0 {
				yellow = true
			}

			c := BlendColors(mode, black, red, yellow)

			switch c {
			case BlendModeB:
				result.Set(x, y, colorBlack)
			case BlendModeR:
				result.Set(x, y, colorRed)
			case BlendModeY:
				result.Set(x, y, colorYellow)
			default:
				result.Set(x, y, colorWhite)
			}
		}
	}

	return result
}

///////////////////////////////////////////////////////////////////////////////

const (
	BlendModeW = 0
	BlendModeB = 1
	BlendModeR = 2
	BlendModeY = 3
)

type BlendMode [3]int

// 0 - White
// 1 - Black
// 2 - Red
// 3 - Yellow
func BlendColors(mode BlendMode, b, r, y bool) int {
	for i := range 3 {
		switch mode[i] {
		case BlendModeB:
			if b {
				return BlendModeB
			}
		case BlendModeR:
			if r {
				return BlendModeR
			}
		case BlendModeY:
			if y {
				return BlendModeY
			}
		default:
			return BlendModeW
		}
	}

	return BlendModeW
}

func StringToBlendMode(mode string) BlendMode {
	if len(mode) != 3 {
		return BlendMode{BlendModeB, BlendModeR, BlendModeY}
	}

	v := BlendMode{BlendModeW, BlendModeW, BlendModeW}

	for i := range 3 {
		switch mode[i] {
		case 'B':
			v[i] = BlendModeB
		case 'R':
			v[i] = BlendModeR
		case 'Y':
			v[i] = BlendModeY
		default:
			v[i] = BlendModeW
		}
	}

	return v
}
