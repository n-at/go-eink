# eink

Output image on GoodDisplay IL075U e-ink display

## Linux permissions

```bash
sudo nano /etc/udev/rules.d/50-myusb.rules
```

Add line:

```txt
SUBSYSTEMS=="usb", ATTRS{idVendor}=="1a86", ATTRS{idProduct}=="7523", GROUP="users", MODE="0666"
```

`idVendor` and `idProduct` can be found in output of `lsusb -vvv`.
Or run program without parameters to list available ports with USB device info.

Reboot machine.
