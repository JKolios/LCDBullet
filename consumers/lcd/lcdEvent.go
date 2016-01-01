package lcd

import (
	"time"

	"github.com/JKolios/goLcdEvents/events"
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

func newShutdownEvent() *LcdEvent {
	return &LcdEvent{EVENT_SHUTDOWN, "", 0 * time.Second, 0, 0, true}
}

func newDisplayEvent(message string, duration time.Duration, flash int, flashRepetitions int, clearAfter bool) *LcdEvent {
	return &LcdEvent{EVENT_DISPLAY, message, duration, flash, flashRepetitions, clearAfter}
}

func monitorlcdEventInputChannel(display *LCDConsumer, lcdEventInput chan events.Event) {
	var incomingEvent events.Event
	var incomingLcdEvent *LcdEvent
	for {
		incomingEvent = <-lcdEventInput

		switch incomingEvent.Type {
		case "pushbullet":
			incomingLcdEvent = newDisplayEvent(incomingEvent.Payload.(string), 8*time.Second, BEFORE, 1, true)
		case "bmp":
			incomingLcdEvent = newDisplayEvent(incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
		case "systeminfo":
			incomingLcdEvent = newDisplayEvent(incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
		case "shutdown":
			incomingLcdEvent = newShutdownEvent()
		default:
			incomingLcdEvent = newDisplayEvent(incomingEvent.Payload.(string), 8*time.Second, NO_FLASH, 1, false)
		}

		display.displayEvent(incomingLcdEvent)
	}

}
