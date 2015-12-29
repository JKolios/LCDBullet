package lcd

import (
	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/utils"
	"github.com/kidoman/embd/controller/hd44780"
	_ "github.com/kidoman/embd/host/rpi"
)

type LCDConsumer struct {
	Driver        *hd44780.HD44780
	LcdEventInput chan events.Event
}

func (consumer *LCDConsumer) Initialize(config conf.Configuration) {
	// Config Parsing
	pinout := config.Pinout
	blPolarity := config.BlPolarity

	// Driver Init
	driver, err := hd44780.NewGPIO(pinout[0], pinout[1], pinout[2], pinout[3], pinout[4], pinout[5], pinout[6], hd44780.BacklightPolarity(blPolarity), hd44780.RowAddress16Col, hd44780.TwoLine, hd44780.DisplayOn)
	utils.LogErrorandExit("Cannot init LCD:", err)

	consumer.Driver = driver
	err = consumer.Driver.Clear()
	utils.LogErrorandExit("Cannot clear LCD:", err)

}

func (consumer *LCDConsumer) Register(LcdEventInput chan events.Event) {

	consumer.LcdEventInput = LcdEventInput
	// Input Monitor Goroutine Startup
	go monitorlcdEventInputChannel(consumer, LcdEventInput)
}
