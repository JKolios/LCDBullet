package bmp

import (
	"fmt"
	"github.com/JKolios/goLcdEvents/utils"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bmp085"
	"strconv"
	"time"
)

type BMPProducer struct {
	sensor     *bmp085.BMP085
	outputChan chan string
}

func NewBMPProducer() *BMPProducer {
	return &BMPProducer{}
}

func (producer *BMPProducer) Initialize(config utils.Configuration) {

	bus := embd.NewI2CBus(config.Bmpi2c)
	producer.sensor = bmp085.New(bus)

}

func (producer *BMPProducer) Subscribe(producerChan chan string) {
	producer.outputChan = producerChan
	go pollBMP085(producer.sensor, producer.outputChan, 10*time.Second)
}
func (producer *BMPProducer) Terminate() {
	//Dummy, no termination needed(?)
}

func pollBMP085(sensor *bmp085.BMP085, output chan string, every time.Duration) {
	for {
		temperature, err := sensor.Temperature()
		utils.LogErrorandExit("Cannot get temperature", err)

		pressure, err := sensor.Pressure()
		utils.LogErrorandExit("Cannot get pressure", err)

		altitude, err := sensor.Altitude()
		utils.LogErrorandExit("Cannot get altitude", err)

		tempStr := strconv.FormatFloat(temperature, 'f', 2, 64)
		pressStr := strconv.Itoa(pressure)
		altStr := strconv.FormatFloat(altitude, 'f', 2, 64)

		finalMessage := fmt.Sprintf("%v T:%v P:%v A:%v", time.Now().Format(time.Kitchen), tempStr, pressStr, altStr)
		output <- finalMessage
		time.Sleep(every)

	}
}
