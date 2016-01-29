// +build linux

package bmp

import (
	"fmt"
	"strconv"

	"github.com/JKolios/EventsToGo/events"
	"github.com/JKolios/EventsToGo/producers"
	"github.com/JKolios/goLcdEvents/utils"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/sensor/bmp085"
	"time"
)

func ProducerSetupFuction(producer *producers.GenericProducer, config map[string]interface{}) {
	producer.RuntimeObjects["sensor"] = bmp085.New(embd.NewI2CBus(config["Bmpi2c"].(byte)))
}

func ProducerWaitFunction(producer *producers.GenericProducer) {
	time.Sleep(30 * time.Second)
}

func ProducerRunFuction(producer *producers.GenericProducer) events.Event {

	sensor := producer.RuntimeObjects["sensor"].(bmp085.BMP085)

	temperature, err := sensor.Temperature()
	utils.LogErrorandExit("Cannot get temperature", err)

	pressure, err := sensor.Pressure()
	utils.LogErrorandExit("Cannot get pressure", err)

	altitude, err := sensor.Altitude()
	utils.LogErrorandExit("Cannot get altitude", err)

	tempStr := strconv.FormatFloat(temperature, 'f', 2, 64)
	pressStr := strconv.Itoa(pressure)
	altStr := strconv.FormatFloat(altitude, 'f', 2, 64)

	finalMessage := fmt.Sprintf("BMP Reports: Temperature:%v Pressure:%v Altitude:%v", tempStr, pressStr, altStr)
	return events.Event{finalMessage, producer.Name, time.Now(), events.PRIORITY_LOW}

}
