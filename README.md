# image to e-ink

Output an arbitrary image on 
[GoodDisplay IL075U](https://www.good-display.com/product/404.html) (black and white)
[GoodDisplay IL075RU](https://www.good-display.com/product/418.html) (black, white and red)
e-ink displays.

IL075RU can work in both modes - color and monochrome.

Tested on Raspberry Pi 4 (arm64), macOS Ventura (amd64, arm64). 

## Build

Go 1.19+ required.

```bash
go build -a -o app
```

## Usage

Run with `-help` flag to get all options:

```txt
-device string
    device name, required, can be obtained with -list flag
-device-mode string
    device mode, bw - black and white (IL075U, IL075RU), bwr - black, white and red (IL075RU) (default "bw")
-eink-screen-refresh-pause int
    pause for screen refresh (ms) (default 5000)
-eink-write-data-pause int
    pause between image chunk writing (ms) (default 100)
-image string
    path to image to print, required
-image-align string
    image alignment, one of: top-left, top-middle, top-right, middle-left, middle, middle-right, bottom-left, bottom-middle, bottom-right (default "middle")
-image-dithering-algo string
    dithering algorithm, one of: floyd_steinberg, jarvis_judice_ninke, atkinson, burkes, stucki, sierra (default "floyd_steinberg")
-image-dithering-threshold int
    dithering threshold, 0..256 (default 128)
-image-enlarge
    enlarge image to fit screen
-list
    show available devices and exit
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
