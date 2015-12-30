package lcd

import (
	"log"
	"math"
	"strconv"
	"time"

	"github.com/JKolios/goLcdEvents/utils"
)

func (display *LCDConsumer) displaySingleFrame(bytes []byte, duration time.Duration) {

	//Display Line 0
	rightBound := 16
	if len(bytes) < 16 {
		rightBound = len(bytes)
	}
	log.Println("Line 0: " + string(bytes[0:rightBound]))
	for _, char := range bytes[0:rightBound] {
		err := display.Driver.WriteChar(char)
		utils.LogErrorandExit("Cannot write char to LCD:", err)
	}

	//Display Line 1
	if len(bytes) > 16 {
		display.Driver.SetCursor(0, 1)
		rightBound = 32
		if len(bytes) < 32 {
			rightBound = len(bytes)
		}
		log.Println("Line 1: " + string(bytes[16:rightBound]))

		for _, char := range bytes[16:rightBound] {
			err := display.Driver.WriteChar(char)
			utils.LogErrorandExit("Cannot write char to LCD:", err)
		}
	}
	//Wait for the given duration
	time.Sleep(duration)
}

//DisplayMessage shows the given message on the display. The message is split in pages if needed (no scrolling is used)
//In general, only strings that can be mapped onto ASCII can be displayed correctly.
func (display *LCDConsumer) displayEvent(event *LcdEvent) {
	log.Println("Displaying message: " + event.message)
	err := display.Driver.Clear()
	utils.LogErrorandExit("Cannot clear LCD:", err)

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
			display.Driver.Clear()
		}

	}

	if event.flash == AFTER || event.flash == BEFORE_AND_AFTER {
		display.flashDisplay(event.flashRepetitions, 1*time.Second)
	}
}

//FlashDisplay will trigger the LCD's display on and off
func (display *LCDConsumer) flashDisplay(repetitions int, duration time.Duration) {
	for i := 0; i < repetitions; i++ {
		err := display.Driver.BacklightOn()
		utils.LogErrorandExit("Failed while flashing display", err)
		time.Sleep(duration / 2)
		err = display.Driver.BacklightOff()
		utils.LogErrorandExit("Failed while flashing display", err)
		time.Sleep(duration / 2)
	}
}
