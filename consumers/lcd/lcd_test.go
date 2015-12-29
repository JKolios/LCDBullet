package lcd

import (
	"testing"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
	_ "github.com/kidoman/embd/host/rpi"
)

var testConfig = conf.Configuration{
	Pinout:     []int{21, 16, 6, 13, 19, 26, 5},
	BlPolarity: true,
}

var testEvent = events.Event{
	Type:    "test",
	Payload: "Testing...",
	From:    nil,
}

func TestInitAndRegister(t *testing.T) {
	consumer := LCDConsumer{}
	testChan := make(chan events.Event)
	consumer.Initialize(testConfig)
	consumer.Register(testChan)
	testChan <- testEvent
}
