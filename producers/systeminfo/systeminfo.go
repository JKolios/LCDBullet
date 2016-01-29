// +build linux

package systeminfo

import (
	"fmt"
	"log"
	"time"

	"github.com/JKolios/EventsToGo/events"
	"github.com/JKolios/EventsToGo/producers"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

func ProducerWaitFunction(producer *producers.GenericProducer) {
	time.Sleep(30 * time.Second)
}

func ProducerRunFuction(producer *producers.GenericProducer) events.Event {

	uptime, err := linuxproc.ReadUptime("/proc/uptime")
	if err != nil {
		log.Fatal("Failed while getting uptime")
	}

	load, err := linuxproc.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Fatal("Failed while getting uptime")
	}

	uptimeStr := fmt.Sprintf("Up:%v ", time.Duration(time.Duration(int64(uptime.Total))*time.Second).String())
	loadStr := fmt.Sprintf("Load:%v %v %v", load.Last1Min, load.Last5Min, load.Last15Min)
	return events.Event{uptimeStr + loadStr, producer.Name, time.Now(), events.PRIORITY_LOW}

}
