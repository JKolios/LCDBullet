// +build linux

package lcd

import (
	"github.com/JKolios/EventsToGo/consumers"
	"github.com/JKolios/EventsToGo/events"
	"github.com/kidoman/embd/controller/hd44780"
	_ "github.com/kidoman/embd/host/rpi"
	"log"

	"time"
)

func RunFunction(consumer *consumers.GenericConsumer, incomingEvent events.Event) {
	var incomingLCDEvent *LcdEvent

	switch incomingEvent.Type {
	case "pushbullet":
		incomingLCDEvent = newLcdEvent(EVENT_DISPLAY, incomingEvent.Payload.(string), 8*time.Second, BEFORE, 1, true)
	case "bmp":
		incomingLCDEvent = newLcdEvent(EVENT_DISPLAY, incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
	case "systeminfo":
		incomingLCDEvent = newLcdEvent(EVENT_DISPLAY, incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
	default:
		incomingLCDEvent = newLcdEvent(EVENT_DISPLAY, incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
	}
	displayEvent(consumer.RuntimeObjects["driver"].(*hd44780.HD44780), incomingLCDEvent)
}

func StopFunction(consumer *consumers.GenericConsumer) {
	err := consumer.RuntimeObjects["driver"].(*hd44780.HD44780).Close()
	if err != nil {
		log.Fatal("Cannot close LCD:", err.Error())
	}
}

func SetupFunction(consumer *consumers.GenericConsumer, config map[string]interface{}) {

	pinout := config["pinout"].([]interface{})

	driver, err := hd44780.NewGPIO(pinout[0].(int), pinout[1].(int), pinout[2].(int), pinout[3].(int), pinout[4].(int), pinout[5].(int), pinout[6].(int), hd44780.BacklightPolarity(config["BlPolarity"].(bool)), hd44780.RowAddress16Col, hd44780.TwoLine, hd44780.DisplayOn, hd44780.Dots5x10)
	consumer.RuntimeObjects["driver"] = driver

	if err != nil {
		log.Fatal("Cannot init LCD:", err.Error())
	}

	err = driver.Clear()
	if err != nil {
		log.Fatal("Cannot clear LCD:", err.Error())
	}

}
