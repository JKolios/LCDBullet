package lcd

import (
	"github.com/kidoman/embd/controller/hd44780"
	_ "github.com/kidoman/embd/host/rpi"
	"log"
	"math"
	"strconv"
	"time"
)

const (
	NO_FLASH         = -1
	BEFORE           = 0
	AFTER            = 1
	BEFORE_AND_AFTER = 2

	EVENT_DISPLAY  = 0
	EVENT_SHUTDOWN = 1
)

//SharedDisplay represents instance of an HD44780 LCD shareable between many goroutines
type SharedDisplay struct {
	driver *hd44780.HD44780
	Input  chan *LcdEvent
}

type LcdEvent struct {
	eventType               int
	message                 string
	duration                time.Duration
	flash, flashRepetitions int
	clearAfter              bool
}

func NewLcdEvent(eventType int, message string, duration time.Duration, flash int, flashRepetitions int, clearAfter bool) *LcdEvent {
	return &LcdEvent{eventType, message, duration, flash, flashRepetitions, clearAfter}
}

func NewShutdownEvent() *LcdEvent {
	return &LcdEvent{EVENT_SHUTDOWN, "", 0 * time.Second, 0, 0, true}
}

func NewDisplayEvent(message string, duration time.Duration, flash int, flashRepetitions int, clearAfter bool) *LcdEvent {
	return &LcdEvent{EVENT_DISPLAY, message, duration, flash, flashRepetitions, clearAfter}
}

//NewDisplay Generates a pointer to a new SharedDisplay instance
func NewDisplay(pinout []int, blPolarity bool) *SharedDisplay {
	driver, err := hd44780.NewGPIO(pinout[0], pinout[1], pinout[2], pinout[3], pinout[4], pinout[5], pinout[6], hd44780.BacklightPolarity(blPolarity), hd44780.RowAddress16Col, hd44780.TwoLine, hd44780.DisplayOn)
	logErrorandExit("Cannot init LCD:", err)
	err = driver.Clear()
	logErrorandExit("Cannot clear LCD:", err)
	input := make(chan *LcdEvent, 100)
	display := SharedDisplay{driver, input}
	go monitorInputChannel(&display, input)
	return &display
}

func (display *SharedDisplay) displaySingleFrame(bytes []byte, duration time.Duration) {

	//Display Line 0
	rightBound := 16
	if len(bytes) < 16 {
		rightBound = len(bytes)
	}
	log.Println("Line 0: " + string(bytes[0:rightBound]))
	for _, char := range bytes[0:rightBound] {
		err := display.driver.WriteChar(char)
		logErrorandExit("Cannot write char to LCD:", err)
	}

	//Display Line 1
	if len(bytes) > 16 {
		display.driver.SetCursor(0, 1)
		rightBound = 32
		if len(bytes) < 32 {
			rightBound = len(bytes)
		}
		log.Println("Line 1: " + string(bytes[16:rightBound]))

		for _, char := range bytes[16:rightBound] {
			err := display.driver.WriteChar(char)
			logErrorandExit("Cannot write char to LCD:", err)
		}
	}
	//Wait for the given duration
	time.Sleep(duration)
}

//DisplayMessage shows the given message on the display. The message is split in pages if needed (no scrolling is used)
//In general, only strings that can be mapped onto ASCII can be displayed correctly.
func (display *SharedDisplay) DisplayEvent(event *LcdEvent) {
	log.Println("Displaying message: " + event.message)
	err := display.driver.Clear()
	logErrorandExit("Cannot clear LCD:", err)

	if event.flash == BEFORE || event.flash == BEFORE_AND_AFTER {
		display.flashDisplay(event.flashRepetitions, 1*time.Second)
	}

	bytes := []byte(event.message)

	frames := int(math.Ceil(float64(len(bytes)) / 32.0))
	log.Println("Frames: " + strconv.Itoa(frames))
	frametime := int64(math.Ceil(float64(event.duration) / float64(frames)))
	log.Println("Frame time: " + strconv.Itoa(int(frametime)))

	for i := 0; i < frames; i++ {
		log.Printf("Displaying frame %v\n", i)
		rightBound := (i + 1) * 32

		if rightBound > len(bytes) {
			rightBound = len(bytes)
		}
		log.Println("Frame Content: " + string(bytes[i*32:rightBound]))
		display.displaySingleFrame(bytes[i*32:rightBound], time.Duration(frametime))

		if i != (frames-1) || event.clearAfter {
			display.driver.Clear()
		}

	}

	if event.flash == AFTER || event.flash == BEFORE_AND_AFTER {
		display.flashDisplay(event.flashRepetitions, 1*time.Second)
	}
}

//FlashDisplay will trigger the LCD's display on and off
func (display *SharedDisplay) flashDisplay(repetitions int, duration time.Duration) {
	for i := 0; i < repetitions; i++ {
		err := display.driver.BacklightOn()
		logErrorandExit("Failed while flashing display", err)
		time.Sleep(duration / 2)
		err = display.driver.BacklightOff()
		logErrorandExit("Failed while flashing display", err)
		time.Sleep(duration / 2)
	}
}

func monitorInputChannel(display *SharedDisplay, input chan *LcdEvent) {
	for {
		incomingEvent := <-input
		switch incomingEvent.eventType {
		case EVENT_DISPLAY:
			display.DisplayEvent(incomingEvent)
		case EVENT_SHUTDOWN:
			display.Close()
			return
		}

	}
}

//Close closes the connection to the display and frees the GPIO pins for other uses
func (display *SharedDisplay) Close() {
	err := display.driver.Close()
	logErrorandExit("Failed while closing display", err)

}

func logErrorandExit(message string, err error) {
	if err != nil {
		log.Fatal(message + err.Error())
	}
}
