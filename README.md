# image to e-ink

Output an arbitrary image on
[GoodDisplay IL075U](https://www.good-display.com/product/404.html) (black and white)
[GoodDisplay IL075RU](https://www.good-display.com/product/418.html) (black, white and red)
[GoodDisplay GDP075FU1](https://www.good-display.com/product/640.html) (black, white, red and yellow)
e-ink displays.

* IL075U works only in black-white mode.
* IL075RU can work in black-white-red and black-white modes.
* GDP075FU1 works only in black-white-red-yellow mode.

Tested on Raspberry Pi 4 (arm64), macOS Ventura (amd64, arm64), Ubuntu 24.04 (amd64).

## Build

Go 1.24+ required.

```bash
go build -a -o app
```

## Usage

Run with `-help` flag to get all options:

```txt
  -device string
    	device name, required, can be obtained with -list flag
  -device-mode string
    	device mode, one of: bw (black and white for IL075U, IL075RU), bwr (black, white and red for IL075RU), bwry (black, white, red and yellow for GDP075FU1) (default "bw")
  -eink-read-device-output
    	read data sent by device (NOTICE: in some cases output may be inconsistent)
  -eink-screen-refresh-pause int
    	pause for screen refresh (ms) (default 5000)
  -eink-write-data-pause int
    	pause between image chunk writing (ms) (default 1000)
  -image string
    	path to image to print, required
  -image-align string
    	image alignment, one of: top-left, top-middle, top-right, middle-left, middle, middle-right, bottom-left, bottom-middle, bottom-right (default "middle")
  -image-blend-mode string
    	combination of letters {B, R, Y} defines order of blending result image from black, red, and yellow components, from top layer to bottom (default "BYR")
  -image-dithering-algo string
    	dithering algorithm for black and white, one of: floyd_steinberg, jarvis_judice_ninke, atkinson, burkes, stucki, sierra (default "floyd_steinberg")
  -image-dithering-threshold int
    	dithering threshold, 0..256 (default 128)
  -image-enlarge
    	enlarge image to fit screen
  -image-red-dithering-algo string
    	dithering algorithm for red color, same values as -image-dithering-algo (default "sierra")
  -image-red-dithering-threshold int
    	red dithering threshold 0..256 (default 128)
  -image-red-hue-threshold int
    	hue threshold for red image (degrees) 0..360 (default 25)
  -image-yellow-dithering-algo string
    	dithering algorithm for yellow color, same values as -image-dithering-algo (default "stucki")
  -image-yellow-dithering-threshold int
    	yellow dithering threshold 0..256 (default 180)
  -image-yellow-hue-threshold int
    	hue threshold for yellow image (degrees) 0..360 (default 25)
  -list
    	show available devices and exit
  -output string
    	output result to file and exit
  -verbose
    	show extended output
```

## Linux USB permissions

```bash
sudo nano /etc/udev/rules.d/50-myusb.rules
```

Add line:

```txt
SUBSYSTEMS=="usb", ATTRS{idVendor}=="1a86", ATTRS{idProduct}=="7523", GROUP="users", MODE="0666"
```

`idVendor` and `idProduct` can be found in output of `lsusb -vvv`.
Or run program with `-list` flag to list available ports with USB device info.

Reboot machine.

## Uses

* [sirupsen/logrus](https://github.com/sirupsen/logrus) - MIT
* [bugst/go-serial](https://github.com/bugst/go-serial) - BSD-3-Clause
* [nfnt/resize](https://github.com/nfnt/resize) - ISC
