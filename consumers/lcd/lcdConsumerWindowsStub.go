// +build !linux

package lcd

import (
	"github.com/JKolios/EventsToGo/consumers"
	"github.com/JKolios/EventsToGo/events"
)

func RunFunction(consumer *consumers.GenericConsumer, incomingEvent events.Event) {}

func StopFunction(consumer *consumers.GenericConsumer) {}

func SetupFunction(consumer *consumers.GenericConsumer, config map[string]interface{}) {
	panic("LCD consumer: not available on windows")

}
