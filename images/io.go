package images

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
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
