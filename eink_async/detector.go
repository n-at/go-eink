package eink_async

import (
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func EnumerateDevices() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Errorf("unable to get serial ports list: %s", err)
		return
	}
	if len(ports) == 0 {
		log.Errorf("no serial ports found")
		return
	}
	for _, port := range ports {
		log.Infof("found serial port: %s", port)
	}
}

func EnumerateDevicesExtended() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Errorf("unable to get detailed serial ports list: %s", err)
		return
	}
	if len(ports) == 0 {
		log.Fatalf("no serial ports found")
		return
	}
	for _, port := range ports {
		if !port.IsUSB {
			continue
		}
		log.Infof("found serial port: name=\"%s\", vendorId=\"%s\", productId=\"%s\"", port.Name, port.VID, port.PID)
	}
}
