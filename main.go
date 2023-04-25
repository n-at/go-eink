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
	deviceMode := flag.String("device-mode", "bw", "device mode, bw - black and white (IL075U, IL075RU), bwr - black, white and red (IL075RU)")
	imagePath := flag.String("image", "", "path to image to print, required")
	imageEnlarge := flag.Bool("image-enlarge", false, "enlarge image to fit screen")
	imageAlign := flag.String("image-align", "middle", "image alignment, one of: top-left, top-middle, top-right, middle-left, middle, middle-right, bottom-left, bottom-middle, bottom-right")
	imageDitheringAlgorithm := flag.String("image-dithering-algo", "floyd_steinberg", "dithering algorithm, one of: floyd_steinberg, jarvis_judice_ninke, atkinson, burkes, stucki, sierra")
	imageDitheringThreshold := flag.Int("image-dithering-threshold", 128, "dithering threshold, 0..256")
	einkWriteDataPause := flag.Int("eink-write-data-pause", 100, "pause between image chunk writing (ms)")
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

	if len(*deviceName) == 0 {
		log.Fatal("device required")
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

	imgBW := images.Dithering(img, &images.PixelTransformationGrayscale{}, images.GetDitheringAlgorithm(*imageDitheringAlgorithm), *imageDitheringThreshold)
	imageDataBW := images.ToImageData(imgBW)

	imgRW := images.Dithering(img, &images.PixelTransformationRed{}, images.GetDitheringAlgorithm(*imageDitheringAlgorithm), *imageDitheringThreshold)
	imageDataRW := images.ToImageData(imgRW)

	//print

	if *deviceMode == eink.DeviceModeBW {
		if err := eink.PrintBW(*deviceName, imageDataBW); err != nil {
			log.Fatalf("unable to print image: %s", err)
		}
	} else if *deviceMode == eink.DeviceModeBWR {
		if err := eink.PrintBWR(*deviceName, imageDataBW, imageDataRW); err != nil {
			log.Fatalf("imable to print image: %s", err)
		}
	} else {
		log.Fatalf("unknown device-mode: %s", *deviceMode)
	}
}
