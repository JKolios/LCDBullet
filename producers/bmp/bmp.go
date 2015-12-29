package bmp

import (
	"fmt"
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
	outputChan chan events.Event
}

func (producer *BMPProducer) Initialize(config conf.Configuration) {

	bus := embd.NewI2CBus(config.Bmpi2c)
	producer.sensor = bmp085.New(bus)

}

func (producer *BMPProducer) Subscribe(producerChan chan events.Event) {
	producer.outputChan = producerChan
	go pollBMP085(producer, 10*time.Second)
}

func pollBMP085(producer *BMPProducer, every time.Duration) {
	for {
		temperature, err := producer.sensor.Temperature()
		utils.LogErrorandExit("Cannot get temperature", err)

		pressure, err := producer.sensor.Pressure()
		utils.LogErrorandExit("Cannot get pressure", err)

		altitude, err := producer.sensor.Altitude()
		utils.LogErrorandExit("Cannot get altitude", err)

		tempStr := strconv.FormatFloat(temperature, 'f', 2, 64)
		pressStr := strconv.Itoa(pressure)
		altStr := strconv.FormatFloat(altitude, 'f', 2, 64)

		finalMessage := fmt.Sprintf("%v T:%v P:%v A:%v", time.Now().Format(time.Kitchen), tempStr, pressStr, altStr)
		finalEvent := events.Event{finalMessage, "bmp", producer}

		producer.outputChan <- finalEvent
		time.Sleep(every)

	}
}
