package bmp

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/utils"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bmp085"
)

type BMPProducer struct {
	sensor     *bmp085.BMP085
	outputChan chan<- events.Event
	done       <-chan struct{}
}

func (producer *BMPProducer) Initialize(config conf.Configuration) {

	bus := embd.NewI2CBus(config.Bmpi2c)
	producer.sensor = bmp085.New(bus)

}

func (producer *BMPProducer) Start(done <-chan struct{}, EventOutput chan<- events.Event) {
	producer.outputChan = EventOutput
	producer.done = done
	log.Println("Initializing BMP085 polling")
	go pollBMP085(producer, 10*time.Second)
}

func pollBMP085(producer *BMPProducer, every time.Duration) {
	tick := time.Tick(every)
	for {
		select {
		case <-producer.done:
			{
				log.Println("pollBMP085 Terminated")
				return
			}
		default:
			log.Println("Starting BMP085 polling")
			temperature, err := producer.sensor.Temperature()
			utils.LogErrorandExit("Cannot get temperature", err)

			pressure, err := producer.sensor.Pressure()
			utils.LogErrorandExit("Cannot get pressure", err)

			altitude, err := producer.sensor.Altitude()
			utils.LogErrorandExit("Cannot get altitude", err)

			tempStr := strconv.FormatFloat(temperature, 'f', 2, 64)
			pressStr := strconv.Itoa(pressure)
			altStr := strconv.FormatFloat(altitude, 'f', 2, 64)

			finalMessage := fmt.Sprintf("BMP Reports: Temperature:%v Pressure:%v Altitude:%v", tempStr, pressStr, altStr)
			finalEvent := events.Event{finalMessage, "bmp", producer, time.Now(), events.PRIORITY_LOW}

			producer.outputChan <- finalEvent
			log.Println("BMP085 polling done")
			<-tick
		}

	}
}
