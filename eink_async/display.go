package eink_async

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"time"
)

type EInkDisplay struct {
	portName  string
	port      serial.Port
	imageData []byte
	success   chan bool
	errors    chan error
}

func (d *EInkDisplay) Print() error {
	log.Debug("open port")
	port, err := serial.Open(d.portName, portMode())
	if err != nil {
		return errors.New(fmt.Sprintf("unable to open port: %s", err))
	}
	defer port.Close()

	log.Debug("set port RTS")
	if err := port.SetRTS(true); err != nil {
		return errors.New(fmt.Sprintf("unable to set RTS: %s", err))
	}

	log.Debug("set port read timeout")
	if err := port.SetReadTimeout(serial.NoTimeout); err != nil {
		return errors.New(fmt.Sprintf("unable to reset read timeout: %s", err))
	}

	d.success = make(chan bool)
	d.errors = make(chan error)
	d.port = port

	//start read loop
	go d.read()

	//handshake
	log.Debug("reset output buffer before handshake")
	if err := port.ResetOutputBuffer(); err != nil {
		return errors.New(fmt.Sprintf("unable to reset otuput buffer: %s", err))
	}
	log.Debug("send handshake request")
	if _, err := port.Write(handshakeRequest()); err != nil {
		return errors.New(fmt.Sprintf("unable to send handshake request: %s", err))
	}

	timeout := time.After(time.Duration(IdleTimeout) * time.Second)

	//wait for response
	select {
	case <-d.success:
		return nil
	case err := <-d.errors:
		return err
	case <-timeout:
		return errors.New("idle timeout exceeded")
	}
}

func (d *EInkDisplay) read() {
	buf := make([]byte, 1024)

	for {
		count, err := d.port.Read(buf)
		if err != nil {
			d.errors <- errors.New(fmt.Sprintf("unable to read data: %s", err))
			break
		}

		time.Sleep(ReadDataPause * time.Millisecond)

		if err := d.port.ResetInputBuffer(); err != nil {
			d.errors <- errors.New(fmt.Sprintf("unable to reset input buffer: %s", err))
			break
		}

		if err := validateHandshakeResponse(buf[:count]); err != nil {
			log.Debugf("non-handshake data received: %s, %s", err, printable(buf[:count]))
		} else {
			log.Debug("handshake valid response received")
			go d.printImage()
		}
	}
}

func (d *EInkDisplay) printImage() {
	chunkIdx := 0

	for chunkStart := 0; chunkStart < len(d.imageData); chunkStart += 4096 {
		chunkLength := min(4096, len(d.imageData)-chunkStart)
		chunk := d.imageData[chunkStart : chunkStart+chunkLength]

		log.Debugf("write chunk #%d (%d bytes)...", chunkIdx, len(chunk))
		if _, err := d.port.Write(chunk); err != nil {
			d.errors <- errors.New(fmt.Sprintf("unable to write chunk #%d: %s", chunkIdx, err))
			return
		}

		log.Debugf("write \\r\\n after chunk #%d...", chunkIdx)
		if _, err := d.port.Write([]byte{CR, LF}); err != nil {
			d.errors <- errors.New(fmt.Sprintf("unable to write \\r\\n after chunk #%d: %s", chunkIdx, err))
			return
		}

		time.Sleep(WriteDataPause * time.Millisecond)

		chunkIdx++
	}

	log.Debugf("finalizing...")

	log.Debug("first \\r\\n")
	if _, err := d.port.Write([]byte{CR, LF}); err != nil {
		d.errors <- errors.New(fmt.Sprintf("unable to write final \\r\\n (1): %s", err))
		return
	}

	log.Debugf("first reset output buffer")
	if err := d.port.ResetOutputBuffer(); err != nil {
		d.errors <- errors.New(fmt.Sprintf("unable to reset output buffer (1): %s", err))
		return
	}

	log.Debugf("second \\r\\n")
	if _, err := d.port.Write([]byte{CR, LF}); err != nil {
		d.errors <- errors.New(fmt.Sprintf("unable to write final \\r\\n (2): %s", err))
		return
	}

	log.Debugf("second reset output buffer")
	if err := d.port.ResetOutputBuffer(); err != nil {
		d.errors <- errors.New(fmt.Sprintf("unable to reset output buffer (2): %s", err))
		return
	}

	time.Sleep(WriteDataPause * time.Millisecond)

	d.success <- true
}
