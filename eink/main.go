package eink

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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

	DisplayModel = 0xc4 //IL075U(R) - black and white (or black, white, red), 7.5 inch

	DeviceModeBW   = "bw"
	DeviceModeBWR  = "bwr"
	DeviceModeBWRY = "bwry"
)

var (
	WriteDataPause     = 1000
	ScreenRefreshPause = 5000
)

///////////////////////////////////////////////////////////////////////////////

func PrintBW(portName string, imageData []byte) error {
	//prepare data

	if !imageDataValid(imageData) {
		return errors.New("image data length mismatch")
	}

	imageData = prepareImageDataBW(imageData)

	//open port

	port, err := preparePort(portName)
	if err != nil {
		return err
	}

	//handshake

	log.Debug("handshake")
	if err := handshake(port, DisplayModel, 0); err != nil {
		return errors.New(fmt.Sprintf("unable to handshake: %s", err))
	} else {
		log.Info("handshake ok")
	}

	//print image

	return printImageCycle(port, imageData)
}

func PrintBWR(portName string, imageDataBW, imageDataRW []byte) error {
	//prepare data

	if !imageDataValid(imageDataBW) {
		return errors.New("BW image data length mismatch")
	}
	if !imageDataValid(imageDataRW) {
		return errors.New("RW image data length mismatch")
	}

	imageData := prepareImageDataBWR(imageDataBW, imageDataRW)

	//open port

	port, err := preparePort(portName)
	if err != nil {
		return err
	}

	//handshake

	log.Debug("handshake")
	if err := handshake(port, DisplayModel, 1); err != nil {
		return errors.New(fmt.Sprintf("unable to handshake: %s", err))
	} else {
		log.Info("handshake ok")
	}

	//print image

	return printImageCycle(port, imageData)
}

func PrintBWRY(portName string, imageDataBW, imageDataRW, imageDataYW []byte) error {
	return errors.New("not implemented") //TODO
}

///////////////////////////////////////////////////////////////////////////////

func preparePort(portName string) (serial.Port, error) {
	log.Debug("test port")
	if err := testPort(portName); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to test port: %s", err))
	}

	log.Debug("open port")
	port, err := serial.Open(portName, portMode())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to open port %s: %s", portName, err))
	}

	//setup port

	log.Debug("set port RTS")
	if err := port.SetRTS(true); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to set RTS: %s", err))
	}

	log.Debug("set port read timeout")
	if err := port.SetReadTimeout(serial.NoTimeout); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to reset read timeout: %s", err))
	}

	return port, nil
}

func handshake(port serial.Port, displayModel, displayRed byte) error {
	log.Debug("send handshake request")
	if _, err := port.Write(handshakeRequest(displayModel, displayRed)); err != nil {
		return errors.New(fmt.Sprintf("unable to send handshake request: %s", err))
	}

	time.Sleep(time.Duration(WriteDataPause) * time.Millisecond)

	log.Debug("read handshake response")
	buf := make([]byte, 1024)
	count, err := port.Read(buf)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to read handshake response: %s", err))
	}

	log.Debugf("handshake response: %s", printable(buf[:count]))

	if err := validateHandshakeResponse(buf[:count]); err != nil {
		return errors.New(fmt.Sprintf("unable to validate handshake response: %s", err))
	}

	return nil
}

func printImageCycle(port serial.Port, imageData []byte) error {
	for {
		ok, err := printImage(port, imageData)
		if err != nil {
			return err
		}
		if ok {
			log.Info("done")
			return nil
		} else {
			log.Info("print failed, retrying")
		}
	}
}

func printImage(port serial.Port, imageData []byte) (bool, error) {
	chunkIdx := 0

	for chunkStart := 0; chunkStart < len(imageData); chunkStart += 4096 {
		chunkLength := min(4096, len(imageData)-chunkStart)
		chunk := imageData[chunkStart : chunkStart+chunkLength]

		log.Debugf("write chunk #%d (%d bytes)", chunkIdx, len(chunk))
		if err := writePortData(port, chunk); err != nil {
			return false, errors.New(fmt.Sprintf("unable to write chunk: %s", err))
		}

		log.Debugf("write CRLF after chunk #%d", chunkIdx)
		if err := writePortData(port, []byte{CR, LF}); err != nil {
			return false, errors.New(fmt.Sprintf("unable to write CRLF after chunk #%d: %s", chunkIdx, err))
		}

		log.Debugf("read data after chunk #%d (1-st line)", chunkIdx)
		if _, err := readPortData(port); err != nil {
			return false, errors.New(fmt.Sprintf("unable to read data: %s", err))
		}

		log.Debugf("read data after chunk #%d (2-nd line)", chunkIdx)
		if _, err := readPortData(port); err != nil {
			return false, errors.New(fmt.Sprintf("unable to read data: %s", err))
		}

		time.Sleep(time.Duration(WriteDataPause) * time.Millisecond)

		chunkIdx++
	}

	log.Info("wait for screen to refresh")
	time.Sleep(time.Duration(ScreenRefreshPause) * time.Millisecond)

	log.Debugf("read remaining data")
	remaining, err := readPortData(port)
	if err != nil {
		return false, errors.New(fmt.Sprintf("unable to read data: %s", err))
	}

	bytesReceived, err := extractReceivedBytes(remaining)
	if err != nil {
		return false, err
	}

	log.Debugf("bytes received: %d", bytesReceived)
	if bytesReceived != len(imageData) {
		return false, nil
	}

	return true, nil
}
