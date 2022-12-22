package eink

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
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
)

func Print(portName string, imageData []byte) {
	if len(imageData) != ImageWidth*ImageHeight/8 {
		log.Errorf("image data length mismatch")
		return
	}

	if err := testPort(portName); err != nil {
		log.Errorf("unable to test port: %s", err)
		return
	}

	mode := &serial.Mode{
		BaudRate: PortBaudRate,
		DataBits: PortDataBits,
		Parity:   PortParity,
		StopBits: PortStopBits,
	}

	log.Info("open port...")
	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Errorf("unable to open port %s: %s", portName, err)
		return
	}

	log.Info("set port RTS...")
	if err := port.SetRTS(true); err != nil {
		log.Errorf("unable to set RTS: %s", err)
		return
	}

	log.Info("set port read timeout...")
	if err := port.SetReadTimeout(serial.NoTimeout); err != nil {
		log.Errorf("unable to reset read timeout: %s", err)
	}

	//handshake

	log.Info("handshake...")
	if err := handshake(port); err != nil {
		log.Errorf("unable to handshake: %s", err)
		return
	} else {
		log.Info("handshake ok")
	}

	//print image

	chunkIdx := 0
	buf := make([]byte, 5000)
	for chunkStart := 0; chunkStart < len(imageData); chunkStart += 4096 {
		chunkLength := min(4096, len(imageData)-chunkStart)
		for chunkPos := 0; chunkPos < chunkLength; chunkPos++ {
			buf[chunkPos] = imageData[chunkStart+chunkPos]
		}
		chunk := buf[:chunkLength]

		log.Infof("write chunk #%d (%d bytes)...", chunkIdx, len(chunk))
		if err := writePortData(port, chunk); err != nil {
			log.Errorf("unable to write chunk: %s", err)
			return
		}
		if err := writePortData(port, []byte{CR, LF}); err != nil {
			log.Errorf("unable to write \\r\\n after chunk: %s", err)
			return
		}

		log.Debugf("read data after chunk #%d...", chunkIdx)
		if d, err := readPortData(port); err != nil {
			log.Errorf("unable to read data: %s", err)
			return
		} else {
			log.Debugf("read %d bytes", len(d))
		}
		if d, err := readPortData(port); err != nil {
			log.Errorf("unable to read data: %s", err)
			return
		} else {
			log.Debugf("read %d bytes", len(d))
		}

		time.Sleep(200 * time.Millisecond)

		chunkIdx++
	}

	//refresh screen

	log.Info("refresh screen...")
	if err := writePortData(port, []byte{CR, LF}); err != nil {
		log.Errorf("unable to write \\r\\n after all data sent: %s", err)
	}

	time.Sleep(5 * time.Second)

	log.Debugf("read remaining data...")
	if d, err := readPortData(port); err != nil {
		log.Errorf("unable to read data: %s", err)
		return
	} else {
		log.Debugf("read %d bytes", len(d))
	}

	log.Info("done")
}

func readPortData(port serial.Port) ([]byte, error) {
	outputBuffer := make([]byte, 1024)

	bytesRead, err := port.Read(outputBuffer)
	if err != nil {
		return nil, err
	}

	return outputBuffer[:bytesRead], nil
}

func writePortData(port serial.Port, data []byte) error {
	if _, err := port.Write(data); err != nil {
		return err
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////

func testPort(portName string) error {
	port, err := serial.Open(portName, &serial.Mode{})
	if err != nil {
		return err
	}
	if err := port.Close(); err != nil {
		return err
	}
	return nil
}

func handshake(port serial.Port) error {
	request := make([]byte, 12)
	request[0] = 0xaa
	request[1] = 0x55
	request[2] = 0xe1
	request[3] = ((ImageWidth * ImageHeight) / 8) / 256
	request[4] = ((ImageWidth * ImageHeight) / 8) % 256
	request[5] = 0
	request[6] = DisplayModel
	request[7] = DisplayRed

	sum := 0
	for i := 0; i < 8; i++ {
		sum += int(request[i])
	}

	request[8] = byte(sum % 256)
	request[9] = 0xff
	request[10] = CR
	request[11] = LF

	if _, err := port.Write(request); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	buf := make([]byte, 1024)
	count, err := port.Read(buf)
	if err != nil {
		return err
	}
	if count != 10 {
		return errors.New("read too few bytes")
	}

	if buf[0] != 0xa0 {
		return errors.New("1-st byte mismatch")
	}
	if buf[1] != 0x50 {
		return errors.New("2-nd byte mismatch")
	}
	if buf[2] != 0xf1 {
		return errors.New("3-rd byte mismatch (connection)")
	}
	if buf[9] != 0xff {
		return errors.New("10-th byte mismatch (FF)")
	}

	sum = 0
	for i := 0; i < 8; i++ {
		sum += int(buf[i])
	}
	if buf[8] != byte(sum%256) {
		return errors.New("8-th byte mismatch (checksum)")
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

///////////////////////////////////////////////////////////////////////////////

func EnumerateDevices() {
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatalf("unable to get ports list: %s", err)
	}
	if len(ports) == 0 {
		log.Fatalf("no serial ports found")
	}
	for _, port := range ports {
		log.Infof("found port: %s", port)
	}
}

func EnumerateDevicesExtended() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatalf("unable to get detailed ports list: %s", err)
	}
	if len(ports) == 0 {
		log.Fatalf("no serial ports found")
	}
	for _, port := range ports {
		log.Infof("found port: name=%s, product=%s, usb=%v, vendorId=%s, productId=%s, serialNumber=%s", port.Name, port.Product, port.IsUSB, port.VID, port.PID, port.SerialNumber)
	}
}
