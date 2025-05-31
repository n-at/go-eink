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

	DisplayModel        = 0xc4 //IL075U(R), GDP075FU1 - BW, BWR, BWRY, 7.5 inch
	DisplayModeByteBWR  = 0x01
	DisplayModeByteBWRY = 0x04
	DisplayModeByteNone = 0x00

	DeviceModeBW   = "bw"
	DeviceModeBWR  = "bwr"
	DeviceModeBWRY = "bwry"
)

var (
	WriteDataPause     = 1000
	ScreenRefreshPause = 5000
	ReadDeviceOutput   = false
)

///////////////////////////////////////////////////////////////////////////////

func PrintBW(portName string, imageData []byte) error {
	if !imageDataValid(imageData) {
		return errors.New("image data length mismatch")
	}

	imageData = prepareImageDataBW(imageData)

	return printImage(portName, DeviceModeBW, imageData)
}

func PrintBWR(portName string, imageDataBW, imageDataRW []byte) error {
	if !imageDataValid(imageDataBW) {
		return errors.New("BW image data length mismatch")
	}
	if !imageDataValid(imageDataRW) {
		return errors.New("RW image data length mismatch")
	}

	imageData := prepareImageDataBWR(imageDataBW, imageDataRW)

	return printImage(portName, DeviceModeBWR, imageData)
}

func PrintBWRY(portName string, imageData []byte) error {
	if !imageDataBWRYValid(imageData) {
		return errors.New("BWRY image data length mismatch")
	}
	return printImage(portName, DeviceModeBWRY, imageData)
}

///////////////////////////////////////////////////////////////////////////////

func preparePort(portName string) (serial.Port, error) {
	log.Debug("test port")
	if err := testPort(portName); err != nil {
		return nil, fmt.Errorf("unable to test port: %s", err)
	}

	log.Debug("open port")
	port, err := serial.Open(portName, portMode())
	if err != nil {
		return nil, fmt.Errorf("unable to open port %s: %s", portName, err)
	}

	//setup port

	log.Debug("set port RTS")
	if err := port.SetRTS(true); err != nil {
		return nil, fmt.Errorf("unable to set RTS: %s", err)
	}

	log.Debug("set port read timeout to unlimited")
	if err := port.SetReadTimeout(serial.NoTimeout); err != nil {
		return nil, fmt.Errorf("unable to set read timeout: %s", err)
	}

	return port, nil
}

func handshake(port serial.Port, displayModel byte, deviceMode string) error {
	log.Debug("send handshake request")
	if _, err := port.Write(handshakeRequest(displayModel, deviceMode)); err != nil {
		return fmt.Errorf("unable to send handshake request: %s", err)
	}

	time.Sleep(time.Duration(WriteDataPause) * time.Millisecond)

	log.Debug("read handshake response")
	buf, err := readPortData(port)
	if err != nil {
		return fmt.Errorf("unable to read handshake response: %s", err)
	}

	log.Debugf("handshake response: %s", printable(buf))

	if err := validateHandshakeResponse(buf); err != nil {
		return fmt.Errorf("unable to validate handshake response: %s", err)
	}

	return nil
}

func printImage(portName string, deviceMode string, imageData []byte) error {
	//open port

	port, err := preparePort(portName)
	if err != nil {
		return err
	}

	//handshake

	log.Debug("handshake")
	if err := handshake(port, DisplayModel, deviceMode); err != nil {
		return fmt.Errorf("unable to handshake: %s", err)
	} else {
		log.Info("handshake ok")
	}

	//print

	return printImageImpl(port, imageData)
}

func printImageImpl(port serial.Port, imageData []byte) error {
	chunkIdx := 0

	for chunkStart := 0; chunkStart < len(imageData); chunkStart += 4096 {
		chunkLength := min(4096, len(imageData)-chunkStart)
		chunk := imageData[chunkStart : chunkStart+chunkLength]

		log.Debugf("write chunk #%d (%d bytes)", chunkIdx, len(chunk))
		if err := writePortData(port, chunk); err != nil {
			return fmt.Errorf("unable to write chunk: %s", err)
		}

		log.Debugf("write CRLF after chunk #%d", chunkIdx)
		if err := writePortData(port, []byte{CR, LF}); err != nil {
			return fmt.Errorf("unable to write CRLF after chunk #%d: %s", chunkIdx, err)
		}

		if ReadDeviceOutput {
			log.Debugf("read data after chunk #%d (1-st line)", chunkIdx)
			if _, err := readPortData(port); err != nil {
				return fmt.Errorf("unable to read data: %s", err)
			}

			log.Debugf("read data after chunk #%d (2-nd line)", chunkIdx)
			if _, err := readPortData(port); err != nil {
				return fmt.Errorf("unable to read data: %s", err)
			}
		}

		time.Sleep(time.Duration(WriteDataPause) * time.Millisecond)

		chunkIdx++
	}

	log.Debug("draining output buffer...")
	if err := port.Drain(); err != nil {
		return fmt.Errorf("unable to drain output buffer: %s", err)
	}

	log.Info("waiting for screen to refresh")
	time.Sleep(time.Duration(ScreenRefreshPause) * time.Millisecond)

	if ReadDeviceOutput {
		log.Debugf("read remaining data")
		remaining, err := readPortData(port)
		if err != nil {
			return fmt.Errorf("unable to read data: %s", err)
		}

		bytesReceived, err := extractReceivedBytes(remaining)
		if err != nil {
			return err
		}

		log.Debugf("bytes received: %d", bytesReceived)
		if bytesReceived != len(imageData) {
			return errors.New("received incorrect number of bytes from display")
		}
	}

	log.Debugf("reset input buffer...")
	if err := port.ResetInputBuffer(); err != nil {
		return fmt.Errorf("unable to reset input buffer: %s", err)
	}

	log.Debugf("reset output buffer...")
	if err := port.ResetOutputBuffer(); err != nil {
		return fmt.Errorf("unable to reset output buffer: %s", err)
	}

	return nil
}
