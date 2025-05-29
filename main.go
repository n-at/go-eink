package main

import (
	"flag"
	"go-eink/eink"
	"go-eink/images"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	verbose := flag.Bool("verbose", false, "show extended output")
	list := flag.Bool("list", false, "show available devices and exit")
	output := flag.String("output", "", "output result to file and exit")
	deviceName := flag.String("device", "", "device name, required, can be obtained with -list flag")
	deviceMode := flag.String("device-mode", "bw", "device mode, one of: bw (black and white for IL075U, IL075RU, GDP075FU1), bwr (black, white and red for IL075RU, GDP075FU1), bwry (black, white, red and yellow for GDP075FU1)")
	imagePath := flag.String("image", "", "path to image to print, required")
	imageEnlarge := flag.Bool("image-enlarge", false, "enlarge image to fit screen")
	imageAlign := flag.String("image-align", "middle", "image alignment, one of: top-left, top-middle, top-right, middle-left, middle, middle-right, bottom-left, bottom-middle, bottom-right")
	imageDitheringAlgorithm := flag.String("image-dithering-algo", "floyd_steinberg", "dithering algorithm, one of: floyd_steinberg, jarvis_judice_ninke, atkinson, burkes, stucki, sierra")
	imageDitheringThreshold := flag.Int("image-dithering-threshold", 128, "dithering threshold, 0..256")
	imageRedDitheringAlgorithm := flag.String("image-red-dithering-algo", "atkinson", "dithering algorithm for red pixels, same values as -image-dithering-algo")
	imageRedDitheringThreshold := flag.Int("image-red-dithering-threshold", 128, "red dithering threshold 0..256")
	imageRedHueThreshold := flag.Int("image-red-hue-threshold", 50, "hue threshold for red image (degrees)")
	imageRedSaturationThreshold := flag.Int("image-red-saturation-threshold", 0, "saturation threshold for red image (%)")
	imageRedLightnessThreshold := flag.Int("image-red-lightness-threshold", 80, "lightness threshold for red image (%)")
	imageYellowDitheringAlgorithm := flag.String("image-yellow-dithering-algo", "floyd_steinberg", "dithering algorithm for yellow pixels, same values as -image-dithering-algo")
	imageYellowDitheringThreshold := flag.Int("image-yellow-dithering-threshold", 128, "yellow dithering threshold 0..256")
	imageYellowHueThreshold := flag.Int("image-yellow-hue-threshold", 50, "hue threshold for yellow image (degrees)")
	imageYellowSaturationThreshold := flag.Int("image-yellow-saturation-threshold", 0, "saturation threshold for yellow image (%)")
	imageYellowLightnessThreshold := flag.Int("image-yellow-lightness-threshold", 90, "lightness threshold for yellow image (%)")
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
		Threshold:              *imageRedDitheringThreshold,
		RedHueThreshold:        *imageRedHueThreshold,
		RedSaturationThreshold: *imageRedSaturationThreshold,
		RedLightnessThreshold:  *imageRedLightnessThreshold,
	}
	imgRW := images.Dithering(img, transformRW, images.GetDitheringAlgorithm(*imageRedDitheringAlgorithm))

	transformYW := &images.PixelTransformationYellow{
		Threshold:                 *imageYellowDitheringThreshold,
		YellowHueThreshold:        *imageYellowHueThreshold,
		YellowSaturationThreshold: *imageYellowSaturationThreshold,
		YellowLightnessThreshold:  *imageYellowLightnessThreshold,
	}
	imgYW := images.Dithering(img, transformYW, images.GetDitheringAlgorithm(*imageYellowDitheringAlgorithm))

	//output?

	if len(*output) > 0 {
		if *deviceMode == eink.DeviceModeBWR {
			imgBW = images.JoinBWR(imgBW, imgRW)
		}
		if *deviceMode == eink.DeviceModeBWRY {
			imgBW = images.JoinBWRY(imgBW, imgRW, imgYW)
		}
		if err := images.Save(imgBW, *output); err != nil {
			log.Fatalf("unable to save image: %s", err)
		}
		return
	}

	//print

	imageDataBW := images.ToImageData(imgBW)
	imageDataRW := images.ToImageData(imgRW)
	imageDataYW := images.ToImageData(imgYW)

	if len(*deviceName) == 0 {
		log.Fatal("device required")
	}

	if *deviceMode == eink.DeviceModeBW {
		if err := eink.PrintBW(*deviceName, imageDataBW); err != nil {
			log.Fatalf("unable to print BW image: %s", err)
		}
	} else if *deviceMode == eink.DeviceModeBWR {
		if err := eink.PrintBWR(*deviceName, imageDataBW, imageDataRW); err != nil {
			log.Fatalf("unable to print BWR image: %s", err)
		}
	} else if *deviceMode == eink.DeviceModeBWRY {
		if err := eink.PrintBWRY(*deviceName, imageDataBW, imageDataRW, imageDataYW); err != nil {
			log.Fatalf("unable to print BWRY image: %s", err)
		}
	} else {
		log.Fatalf("unknown device-mode: %s", *deviceMode)
	}
}
