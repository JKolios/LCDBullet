package lcd

import (
	"time"

	"github.com/JKolios/goLcdEvents/events"
	"log"
)

const (
	NO_FLASH         = -1
	BEFORE           = 0
	AFTER            = 1
	BEFORE_AND_AFTER = 2

	EVENT_DISPLAY  = 0
	EVENT_SHUTDOWN = 1
)

type LcdEvent struct {
	eventType               int
	message                 string
	duration                time.Duration
	flash, flashRepetitions int
	clearAfter              bool
}

func newLcdEvent(eventType int, message string, duration time.Duration, flash int, flashRepetitions int, clearAfter bool) *LcdEvent {
	return &LcdEvent{eventType, message, duration, flash, flashRepetitions, clearAfter}
}

func monitorlcdEventInputChannel(display *LCDConsumer, lcdEventInput chan events.Event) {
	var incomingEvent events.Event
	var incomingLcdEvent *LcdEvent

	EventLoop:
	for {
		incomingEvent = <-lcdEventInput

		switch incomingEvent.Type {
		case "pushbullet":
			incomingLcdEvent = newLcdEvent(EVENT_DISPLAY,incomingEvent.Payload.(string), 8*time.Second, BEFORE, 1, true)
		case "bmp":
			incomingLcdEvent = newLcdEvent(EVENT_DISPLAY, incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
		case "systeminfo":
			incomingLcdEvent = newLcdEvent(EVENT_DISPLAY, incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
		case "shutdown":
			incomingLcdEvent = newLcdEvent(EVENT_SHUTDOWN, "Shutting down...", 5 * time.Second, BEFORE, 1, true)
			display.displayEvent(incomingLcdEvent)
			break EventLoop

		default:
			incomingLcdEvent = newLcdEvent(EVENT_DISPLAY, incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
		}
		display.displayEvent(incomingLcdEvent)

	}
	log.Println("LCD monitor goroutine exiting...")
	lcdEventInput <- incomingEvent

}
