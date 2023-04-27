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
	output := flag.String("output", "", "output result to file and exit")
	deviceName := flag.String("device", "", "device name, required, can be obtained with -list flag")
	deviceMode := flag.String("device-mode", "bw", "device mode, one of: bw (black and white for IL075U, IL075RU), bwr (black, white and red for IL075RU)")
	imagePath := flag.String("image", "", "path to image to print, required")
	imageEnlarge := flag.Bool("image-enlarge", false, "enlarge image to fit screen")
	imageAlign := flag.String("image-align", "middle", "image alignment, one of: top-left, top-middle, top-right, middle-left, middle, middle-right, bottom-left, bottom-middle, bottom-right")
	imageDitheringAlgorithm := flag.String("image-dithering-algo", "floyd_steinberg", "dithering algorithm, one of: floyd_steinberg, jarvis_judice_ninke, atkinson, burkes, stucki, sierra")
	imageDitheringThreshold := flag.Int("image-dithering-threshold", 128, "dithering threshold, 0..256")
	imageRedHueThreshold := flag.Int("image-red-hue-threshold", 25, "hue threshold for red image (degrees)")
	imageRedSaturationThreshold := flag.Int("image-red-saturation-threshold", 40, "saturation threshold for red image (%)")
	imageRedLightnessThreshold := flag.Int("image-red-lightness-threshold", 80, "lightness threshold for red image (%)")
	imageSubtract := flag.String("image-subtract", "none", "subtract images, one of: none (keep both), black (subtract black-and-white from red), red (subtract red from black-and-white)")
	einkWriteDataPause := flag.Int("eink-write-data-pause", 300, "pause between image chunk writing (ms)")
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

	//prepare settings

	eink.WriteDataPause = *einkWriteDataPause
	eink.ScreenRefreshPause = *einkScreenRefreshPause

	//list devices

	if *list {
		eink.EnumerateDevicesExtended()
		return
	}

	//prepare image

	if len(*imagePath) == 0 {
		log.Fatal("image required")
	}
	img, err := images.Open(*imagePath)
	if err != nil {
		log.Fatalf("unable to open image: %s", err)
	}
	img = images.Resize(img, eink.ImageWidth, eink.ImageHeight, *imageEnlarge)
	img = images.Align(img, eink.ImageWidth, eink.ImageHeight, images.GetAlign(*imageAlign))

	transformBW := &images.PixelTransformationGrayscale{
		Threshold: *imageDitheringThreshold,
	}
	imgBW := images.Dithering(img, transformBW, images.GetDitheringAlgorithm(*imageDitheringAlgorithm))

	transformRW := &images.PixelTransformationRed{
		Threshold:              *imageDitheringThreshold,
		RedHueThreshold:        *imageRedHueThreshold,
		RedSaturationThreshold: *imageRedSaturationThreshold,
		RedLightnessThreshold:  *imageRedLightnessThreshold,
	}
	imgRW := images.Dithering(img, transformRW, images.GetDitheringAlgorithm(*imageDitheringAlgorithm))

	switch *imageSubtract {
	case images.SubtractBlack:
		imgRW = images.Subtract(imgRW, imgBW)
	case images.SubtractRed:
		imgBW = images.Subtract(imgBW, imgRW)
	}

	imageDataBW := images.ToImageData(imgBW)
	imageDataRW := images.ToImageData(imgRW)

	//output?

	if len(*output) > 0 {
		if *deviceMode == eink.DeviceModeBWR {
			imgBW = images.Join(imgBW, imgRW)
		}
		if err := images.Save(imgBW, *output); err != nil {
			log.Fatalf("unable to save image: %s", err)
		}
		return
	}

	//print

	if len(*deviceName) == 0 {
		log.Fatal("device required")
	}

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
