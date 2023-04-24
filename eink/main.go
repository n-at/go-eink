package eink

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"strconv"
	"strings"
	"time"
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

	DeviceModeBW  = "bw"
	DeviceModeBWR = "bwr"
)

var (
	WriteDataPause     = 1000
	ScreenRefreshPause = 5000
)

func Print(portName string, mode string, imageDataBW []byte, imageDataRed []byte) error {
	//TODO print red image

	//prepare data

	if !imageDataValid(imageDataBW) {
		return errors.New("image data length mismatch")
	}

	imageDataBW = prepareImageData(imageDataBW)

	//open port

	log.Debug("test port")
	if err := testPort(portName); err != nil {
		return errors.New(fmt.Sprintf("unable to test port: %s", err))
	}

	log.Debug("open port")
	port, err := serial.Open(portName, portMode())
	if err != nil {
		return errors.New(fmt.Sprintf("unable to open port %s: %s", portName, err))
	}

	//setup port

	log.Debug("set port RTS")
	if err := port.SetRTS(true); err != nil {
		return errors.New(fmt.Sprintf("unable to set RTS: %s", err))
	}

	log.Debug("set port read timeout")
	if err := port.SetReadTimeout(serial.NoTimeout); err != nil {
		return errors.New(fmt.Sprintf("unable to reset read timeout: %s", err))
	}

	//handshake

	log.Debug("handshake")
	if err := handshake(port); err != nil {
		return errors.New(fmt.Sprintf("unable to handshake: %s", err))
	} else {
		log.Info("handshake ok")
	}

	//print image

	for {
		ok, err := printImage(port, imageDataBW)
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

func handshake(port serial.Port) error {
	log.Debug("send handshake request")
	if _, err := port.Write(handshakeRequest()); err != nil {
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

	parts := strings.Split(string(remaining), "=")
	if len(parts) != 2 {
		return false, errors.New("malformed remaining data received")
	}

	bytesReceived, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, errors.New(fmt.Sprintf("unable to read bytes received count: %s", err))
	}

	log.Debugf("bytes received: %d", bytesReceived)
	if bytesReceived != (ImageHeight*ImageWidth)/8 {
		return false, nil
	}

	return true, nil
}
