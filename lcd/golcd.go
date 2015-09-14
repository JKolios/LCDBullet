package lcd

import (
	"github.com/kidoman/embd/controller/hd44780"
	_ "github.com/kidoman/embd/host/rpi"
	"log"
	"math"
	"sync"
	"time"
)

//SharedDisplay represents instance of a HD44780 LCD shareable between many goroutines
type SharedDisplay struct {
	driver *hd44780.HD44780
	mutex  sync.Mutex
}

//NewDisplay Generates a pointer to a new SharedDisplay instance
func NewDisplay(rs, en, d4, d5, d6, d7, backlight int, blPolarity bool) *SharedDisplay {
	driver, err := hd44780.NewGPIO(rs, en, d4, d5, d6, d7, backlight, hd44780.BacklightPolarity(blPolarity), hd44780.RowAddress16Col, hd44780.TwoLine, hd44780.DisplayOn)
	logErrorandExit("Cannot init LCD:", err)
	err = driver.Clear()
	logErrorandExit("Cannot clear LCD:", err)
	return &SharedDisplay{driver, sync.Mutex{}}
}

func logErrorandExit(message string, err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
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
	err := display.driver.Clear()
	logErrorandExit("Cannot clear LCD:", err)

}

//DisplayMessage shows the given message on the display. The message is split in pages if needed (no scrolling is used)
//In general, only strings that can be mapped onto ASCII can be displayed correctly.
func (display *SharedDisplay) DisplayMessage(message string, duration time.Duration) {
	display.mutex.Lock()
	err := display.driver.Clear()
	logErrorandExit("Cannot clear LCD:", err)

	bytes := []byte(message)

	frames := int(math.Ceil(float64(len(bytes)) / 32.0))
	frametime := int64(math.Ceil(float64(duration) / float64(frames)))

	for i := 0; i < frames; i++ {
		log.Printf("Displaying frame %v\n", i)
		rightBound := (i + 1) * 32

		if rightBound > len(bytes) {
			rightBound = len(bytes)
		}
		display.displaySingleFrame(bytes[i*32:rightBound], time.Duration(frametime))
	}
	display.mutex.Unlock()
}

//FlashDisplay will trigger the LCD's display on and off
func (display *SharedDisplay) FlashDisplay(repetitions int, duration time.Duration) {
	display.mutex.Lock()
	for i := 0; i < repetitions; i++ {
		err := display.driver.BacklightOff()
		logErrorandExit("Failed while flashing display", err)
		time.Sleep(duration / 2)
		err = display.driver.BacklightOn()
		logErrorandExit("Failed while flashing display", err)
		time.Sleep(duration / 2)
	}
	display.mutex.Unlock()
}

//Close closes the connection to the display and frees the GPIO pins for other uses
func (display *SharedDisplay) Close() {
	err := display.driver.Close()
	logErrorandExit("failed while closing display", err)
}
