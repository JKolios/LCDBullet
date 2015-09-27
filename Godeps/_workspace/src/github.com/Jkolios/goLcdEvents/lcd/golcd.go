package lcd

import (
	"github.com/JKolios/goLcdEvents/Godeps/_workspace/src/github.com/kidoman/embd/controller/hd44780"
	_ "github.com/JKolios/goLcdEvents/Godeps/_workspace/src/github.com/kidoman/embd/host/rpi"
	"log"
	"math"
	"sync"
	"time"
)

const (
	NO_FLASH         = -1
	BEFORE           = 0
	AFTER            = 1
	BEFORE_AND_AFTER = 2
)

//SharedDisplay represents instance of an HD44780 LCD shareable between many goroutines
type SharedDisplay struct {
	driver *hd44780.HD44780
	mutex  sync.Mutex
	Input  chan *LcdEvent
}

type LcdEvent struct {
	message                 string
	duration                time.Duration
	flash, flashRepetitions int
	clearAfter              bool
}

func NewLcdEvent(message string, duration time.Duration, flash int, flashRepetitions int, clearAfter bool) *LcdEvent {
	return &LcdEvent{message, duration, flash, flashRepetitions, clearAfter}
}

//NewDisplay Generates a pointer to a new SharedDisplay instance
func NewDisplay(pinout []int, blPolarity bool) *SharedDisplay {
	driver, err := hd44780.NewGPIO(pinout[0], pinout[1], pinout[2], pinout[3], pinout[4], pinout[5], pinout[6], hd44780.BacklightPolarity(blPolarity), hd44780.RowAddress16Col, hd44780.TwoLine, hd44780.DisplayOn)
	logErrorandExit("Cannot init LCD:", err)
	err = driver.Clear()
	logErrorandExit("Cannot clear LCD:", err)
	input := make(chan *LcdEvent, 100)
	display := SharedDisplay{driver, sync.Mutex{}, input}
	go monitorInputChannel(&display, input)
	return &display
}

func (display *SharedDisplay) displaySingleFrame(bytes []byte, duration time.Duration) {

	//Display Line 0
	rightBound := 16
	if len(bytes) < 16 {
		rightBound = len(bytes)
	}
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
	display.mutex.Lock()
	err := display.driver.Clear()
	logErrorandExit("Cannot clear LCD:", err)

	if event.flash == BEFORE || event.flash == BEFORE_AND_AFTER {
		display.flashDisplay(event.flashRepetitions, 1*time.Second)
	}

	bytes := []byte(event.message)

	frames := int(math.Ceil(float64(len(bytes)) / 32.0))
	frametime := int64(math.Ceil(float64(event.duration) / float64(frames)))

	for i := 0; i < frames; i++ {
		log.Printf("Displaying frame %v\n", i)
		rightBound := (i + 1) * 32

		if rightBound > len(bytes) {
			rightBound = len(bytes)
		}
		display.displaySingleFrame(bytes[i*32:rightBound], time.Duration(frametime))

		if i != (frames-1) || event.clearAfter {
			display.driver.Clear()
		}

	}

	if event.flash == AFTER || event.flash == BEFORE_AND_AFTER {
		display.flashDisplay(event.flashRepetitions, 1*time.Second)
	}
	display.mutex.Unlock()
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
		display.DisplayEvent(<-input)
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
