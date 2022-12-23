package main

import (
	log "github.com/sirupsen/logrus"
	"go-eink/eink"
	"math/rand"
	"os"
	"time"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	//log.SetLevel(log.InfoLevel)
}

func main() {
	rnd := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	image := make([]byte, (800*480)/8)
	for i := 0; i < len(image); i++ {
		image[i] = byte(rnd.Intn(256))
		//image[i] = byte(i % 256)
		//image[i] = 0
	}

	eink.EnumerateDevicesExtended()
	eink.Print("/dev/cu.usbserial-14140", image)
}
