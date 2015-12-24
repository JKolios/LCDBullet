package lcd

import (
	"github.com/JKolios/goLcdEvents/utils"
	"github.com/kidoman/embd/controller/hd44780"
	_ "github.com/kidoman/embd/host/rpi"
	"time"
)

//SharedDisplay represents instance of an HD44780 LCD shareable between many goroutines
type LCDConsumer struct {
	driver        *hd44780.HD44780
	lcdEventInput chan *LcdEvent
}

//NewDisplay Creates and initialises a new LCD consumer
func NewLCDConsumer() *LCDConsumer {
	return &LCDConsumer{}
}

func (consumer *LCDConsumer) Initialize(config utils.Configuration) {
	// Config Parsing
	pinout := config.Pinout
	blPolarity := config.BlPolarity

	// Driver Init
	driver, err := hd44780.NewGPIO(pinout[0], pinout[1], pinout[2], pinout[3], pinout[4], pinout[5], pinout[6], hd44780.BacklightPolarity(blPolarity), hd44780.RowAddress16Col, hd44780.TwoLine, hd44780.DisplayOn)
	utils.LogErrorandExit("Cannot init LCD:", err)
	err = driver.Clear()
	utils.LogErrorandExit("Cannot clear LCD:", err)
	consumer.driver = driver

	// Input Channel Init
	lcdEventInput := make(chan *LcdEvent, 100)
	consumer.lcdEventInput = lcdEventInput

	// Input Monitor Goroutine Startup
	go monitorlcdEventInputChannel(consumer, lcdEventInput)
}

func (consumer *LCDConsumer) Consume(mType, message string) {
	switch mType {
	case "pushbullet":
		consumer.lcdEventInput <- newDisplayEvent(message, 5*time.Second, BEFORE, 1, true)
	case "bmp":
		consumer.lcdEventInput <- newDisplayEvent(message, 3*time.Second, NO_FLASH, 1, false)
	case "systeminfo":
		consumer.lcdEventInput <- newDisplayEvent(message, 3*time.Second, NO_FLASH, 1, false)
	default:
		consumer.lcdEventInput <- newDisplayEvent(message, 5*time.Second, NO_FLASH, 1, false)
	}
}

//Terminate closes the connection to the display and frees the GPIO pins for other uses
func (consumer *LCDConsumer) Terminate() {
	consumer.lcdEventInput <- newShutdownEvent()
}
