# image to e-ink

Output an arbitrary image on
[GoodDisplay IL075U](https://www.good-display.com/product/404.html) (black and white)
[GoodDisplay IL075RU](https://www.good-display.com/product/418.html) (black, white and red)
[GoodDisplay GDP075FU1](https://www.good-display.com/product/640.html) (black, white, red and yellow)
e-ink displays.

IL075RU can work in black-white-red and black-white modes.

GDP075FU1 works only in black-white-red-yellow mode.

Tested on Raspberry Pi 4 (arm64), macOS Ventura (amd64, arm64).

## Build

Go 1.24+ required.

```bash
go build -a -o app
```

## Usage

Run with `-help` flag to get all options:

```txt
TODO
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
