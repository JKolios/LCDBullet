// +build !linux

package bmp

import (
	"github.com/JKolios/EventsToGo/events"
	"github.com/JKolios/EventsToGo/producers"
)

func ProducerSetupFuction(producer *producers.GenericProducer, config map[string]interface{}) {
	panic("BMP producer: not available on windows")
}

func ProducerWaitFunction(producer *producers.GenericProducer) {
}

func ProducerRunFuction(producer *producers.GenericProducer) events.Event {
	return events.Event{}
}
