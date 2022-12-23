package main

import (
	log "github.com/sirupsen/logrus"
	"go-eink/eink"
	"go-eink/images"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	//log.SetLevel(log.InfoLevel)

	img, err := images.Open("./assets/img.jpg")
	if err != nil {
		log.Fatalf("unable to open image: %s", err)
	}

	img = images.Resize(img, eink.ImageWidth, eink.ImageHeight, false)
	img = images.Align(img, eink.ImageWidth, eink.ImageHeight, images.AlignMiddle)
	img = images.Dithering(img, images.DitheringFloydSteinberg, 128)
	imageData := images.ToImageData(img)

	eink.EnumerateDevicesExtended()
	if err := eink.Print("/dev/cu.usbserial-14140", imageData); err != nil {
		log.Errorf("unable to print image: %s", err)
	}

	//if err := images.Save(img, "./assets/output.png"); err != nil {
	//	log.Fatalf("unable to save image: %s", err)
	//}
}
