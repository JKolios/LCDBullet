// +build !linux

package systeminfo

import (
	"github.com/JKolios/EventsToGo/events"
	"github.com/JKolios/EventsToGo/producers"
)

func ProducerSetupFuction(producer *producers.GenericProducer, config map[string]interface{}) {
	panic("Systeminfo producer: not available in windows")
}

func ProducerWaitFunction(producer *producers.GenericProducer) {
}

func ProducerRunFuction(producer *producers.GenericProducer) events.Event {
	return events.Event{}
}
