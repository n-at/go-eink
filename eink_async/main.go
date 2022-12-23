package eink_async

import (
	"errors"
	"go.bug.st/serial"
)

const (
	PortBaudRate = 115200
	PortStopBits = serial.OneStopBit
	PortParity   = serial.NoParity
	PortDataBits = 8

	CR = 0x0d
	LF = 0x0a

	ImageWidth  = 800
	ImageHeight = 480

	DisplayModel = 0xc4 //IL075U - black and white, 7.5 inch
	DisplayRed   = 0

	ReadDataPause  = 50
	WriteDataPause = 1000
)

var (
	IdleTimeout = 30
)

func New(portName string, imageData []byte) (*EInkDisplay, error) {
	if err := testPort(portName); err != nil {
		return nil, err
	}
	if !imageDataValid(imageData) {
		return nil, errors.New("image data wrong length")
	}

	display := &EInkDisplay{
		portName:  portName,
		imageData: imageData,
	}

	return display, nil
}
