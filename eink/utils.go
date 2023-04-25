package eink

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"strconv"
	"strings"
)

///////////////////////////////////////////////////////////////////////////////

func handshakeRequest(displayModel, displayRed byte) []byte {
	request := make([]byte, 12)
	request[0] = 0xaa
	request[1] = 0x55
	request[2] = 0xe1
	request[3] = ((ImageWidth * ImageHeight) / 8) / 256
	request[4] = ((ImageWidth * ImageHeight) / 8) % 256
	request[5] = 0
	request[6] = displayModel
	request[7] = displayRed

	sum := 0
	for i := 0; i < 8; i++ {
		sum += int(request[i])
	}

	request[8] = byte(sum % 256)
	request[9] = 0xff
	request[10] = CR
	request[11] = LF

	return request
}

func validateHandshakeResponse(response []byte) error {
	if len(response) != 10 {
		return errors.New("wrong length")
	}

	if response[0] != 0xa0 {
		return errors.New("1-st byte mismatch")
	}
	if response[1] != 0x50 {
		return errors.New("2-nd byte mismatch")
	}
	if response[2] != 0xf1 {
		return errors.New("3-rd byte mismatch (connection)")
	}
	if response[9] != 0xff {
		return errors.New("10-th byte mismatch (FF)")
	}

	sum := 0
	for i := 0; i < 8; i++ {
		sum += int(response[i])
	}
	if response[8] != byte(sum%256) {
		return errors.New("8-th byte mismatch (checksum)")
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func testPort(portName string) error {
	port, err := serial.Open(portName, portMode())
	if err != nil {
		return err
	}
	if err := port.Close(); err != nil {
		return err
	}
	return nil
}

func portMode() *serial.Mode {
	return &serial.Mode{
		BaudRate: PortBaudRate,
		DataBits: PortDataBits,
		Parity:   PortParity,
		StopBits: PortStopBits,
	}
}

func readPortData(port serial.Port) ([]byte, error) {
	buf := make([]byte, 1024)

	count, err := port.Read(buf)
	if err != nil {
		return nil, err
	}

	log.Debugf("read %d bytes: \"%s\"", count, printable(buf[:count]))

	return buf[:count], nil
}

func writePortData(port serial.Port, data []byte) error {
	if _, err := port.Write(data); err != nil {
		return err
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////

func imageDataValid(imageData []byte) bool {
	return len(imageData) == (ImageHeight*ImageWidth)/8
}

func prepareImageDataBW(imageData []byte) []byte {
	prepared := make([]byte, len(imageData))

	for idx := 0; idx < len(imageData); idx++ {
		if imageData[idx] == 13 {
			prepared[idx] = 12
		} else {
			prepared[idx] = imageData[idx]
		}
	}

	return prepared
}

func prepareImageDataBWR(imageDataBW, imageDataRW []byte) []byte {
	prepared := make([]byte, len(imageDataBW)+len(imageDataRW))
	offset := 0

	for idx := 0; idx < len(imageDataBW); idx++ {
		if imageDataBW[idx] == 13 {
			prepared[idx+offset] = 12
		} else {
			prepared[idx+offset] = imageDataBW[idx]
		}
	}

	offset += len(imageDataBW)

	for idx := 0; idx < len(imageDataRW); idx++ {
		if imageDataRW[idx] == 13 {
			prepared[idx+offset] = 12
		} else {
			prepared[idx+offset] = imageDataRW[idx]
		}
	}

	return prepared
}

func extractReceivedBytes(data []byte) (int, error) {
	parts := strings.Split(string(data), "=")
	if len(parts) != 2 {
		return 0, errors.New("malformed data")
	}

	bytesReceived, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New(fmt.Sprintf("unable to read bytes data count: %s", err))
	}

	return bytesReceived, nil
}

func printable(buf []byte) string {
	sb := strings.Builder{}

	for _, c := range buf {
		if c >= 0x20 && c < 0x7F {
			sb.WriteRune(rune(c))
		} else {
			sb.WriteRune('.')
		}
	}

	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
