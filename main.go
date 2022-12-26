package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"go-eink/eink"
	"go-eink/images"
	"os"
)

func main() {
	verbose := flag.Bool("verbose", false, "show extended output")
	list := flag.Bool("list", false, "show available devices and exit")
	deviceName := flag.String("device", "", "device name, required, can be obtained with -list flag")
	imagePath := flag.String("image", "", "path to image to print, required")
	imageEnlarge := flag.Bool("image-enlarge", false, "enlarge image to fit screen")
	imageAlign := flag.String("image-align", "middle", "image alignment, one of: top-left, top-middle, top-right, middle-left, middle, middle-right, bottom-left, bottom-middle, bottom-right")
	imageDitheringAlgorithm := flag.String("image-dithering-algo", "floyd_steinberg", "dithering algorithm, one of: floyd_steinberg, jarvis_judice_ninke, atkinson, burkes, stucki, sierra")
	imageDitheringThreshold := flag.Int("image-dithering-threshold", 128, "dithering threshold, 0..256")
	einkWriteDataPause := flag.Int("eink-write-data-pause", 1000, "pause between image chunk writing (ms)")
	einkScreenRefreshPause := flag.Int("eink-screen-refresh-pause", 5000, "pause for screen refresh (ms)")
	flag.Parse()

	//prepare logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	eink.WriteDataPause = *einkWriteDataPause
	eink.ScreenRefreshPause = *einkScreenRefreshPause

	//list devices
	if *list {
		eink.EnumerateDevicesExtended()
		return
	}

	//prepare imagePath
	if len(*imagePath) == 0 {
		log.Fatal("image required")
	}
	img, err := images.Open(*imagePath)
	if err != nil {
		log.Fatalf("unable to open image: %s", err)
	}
	img = images.Resize(img, eink.ImageWidth, eink.ImageHeight, *imageEnlarge)
	img = images.Align(img, eink.ImageWidth, eink.ImageHeight, images.GetAlign(*imageAlign))
	img = images.Dithering(img, images.GetDitheringAlgorithm(*imageDitheringAlgorithm), *imageDitheringThreshold)
	imageData := images.ToImageData(img)

	//print
	if len(*deviceName) == 0 {
		log.Fatal("device required")
	}
	if err := eink.Print(*deviceName, imageData); err != nil {
		log.Errorf("unable to print image: %s", err)
	}
}
