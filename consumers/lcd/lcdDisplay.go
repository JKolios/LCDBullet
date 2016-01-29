package lcd

import (
	"log"
	"math"
	"time"

	"github.com/JKolios/goLcdEvents/utils"
	"github.com/kidoman/embd/controller/hd44780"
)

func displaySingleFrame(display *hd44780.HD44780, bytes []byte, duration time.Duration) {

	//Display Line 0
	rightBound := 16
	if len(bytes) < 16 {
		rightBound = len(bytes)
	}
	// log.Println("Line 0: " + string(bytes[0:rightBound]))
	for _, char := range bytes[0:rightBound] {
		err := display.WriteChar(char)
		utils.LogErrorandExit("Cannot write char to LCD:", err)
	}

	//Display Line 1
	if len(bytes) > 16 {
		display.SetCursor(0, 1)
		rightBound = 32
		if len(bytes) < 32 {
			rightBound = len(bytes)
		}
		// log.Println("Line 1: " + string(bytes[16:rightBound]))

		for _, char := range bytes[16:rightBound] {
			err := display.WriteChar(char)
			utils.LogErrorandExit("Cannot write char to LCD:", err)
		}
	}
	//Wait for the given duration
	time.Sleep(duration)
}

//DisplayMessage shows the given message on the display. The message is split in pages if needed (no scrolling is used)
//In general, only strings that can be mapped onto ASCII can be displayed correctly.
func displayEvent(display *hd44780.HD44780, event *LcdEvent) {
	log.Println("LCD: Displaying message: " + event.message)
	err := display.Clear()
	utils.LogErrorandExit("Cannot clear LCD:", err)

	if event.flash == BEFORE || event.flash == BEFORE_AND_AFTER {
		flashDisplay(display, event.flashRepetitions, 1*time.Second)
	}

	bytes := []byte(event.message)

	frames := int(math.Ceil(float64(len(bytes)) / 32.0))
	// log.Println("Frames: " + strconv.Itoa(frames))
	frametime := int64(math.Ceil(float64(event.duration) / float64(frames)))
	// log.Println("Frame time: " + strconv.Itoa(int(frametime)))

	for i := 0; i < frames; i++ {
		// log.Printf("Displaying frame %v\n", i)
		rightBound := (i + 1) * 32

		if rightBound > len(bytes) {
			rightBound = len(bytes)
		}
		// log.Println("Frame Content: " + string(bytes[i*32:rightBound]))
		displaySingleFrame(display, bytes[i*32:rightBound], time.Duration(frametime))

		if i != (frames-1) || event.clearAfter {
			display.Clear()
		}

	}

	if event.flash == AFTER || event.flash == BEFORE_AND_AFTER {
		flashDisplay(display, event.flashRepetitions, 1*time.Second)
	}

	if event.eventType == EVENT_SHUTDOWN {
		log.Println("LCD: Driver shutting down...")
		display.Close()
	}
}

//FlashDisplay will trigger the LCD's display on and off
func flashDisplay(display *hd44780.HD44780, repetitions int, duration time.Duration) {
	for i := 0; i < repetitions; i++ {
		err := display.BacklightOn()
		utils.LogErrorandExit("Failed while flashing display", err)
		time.Sleep(duration / 2)
		err = display.BacklightOff()
		utils.LogErrorandExit("Failed while flashing display", err)
		time.Sleep(duration / 2)
	}
}
