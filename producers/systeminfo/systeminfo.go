package systeminfo

import (
	"fmt"
	"log"
	"time"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
	linuxproc "github.com/c9s/goprocinfo/linux"
)

type SystemInfoProducer struct {
	outputChan chan events.Event
}

func (producer *SystemInfoProducer) Initialize(config conf.Configuration) {
	//Dummy
	return
}

func (producer *SystemInfoProducer) Subscribe(producerChan chan events.Event) {
	producer.outputChan = producerChan
	go pollSystemInfo(producer, time.Second*10)
}

func pollSystemInfo(producer *SystemInfoProducer, every time.Duration) {
	tick := time.Tick(every)
	for {
		log.Println("Starting systeminfo polling")
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
		uptimeEvent := events.Event{uptimeStr + loadStr, "systeminfo", producer, time.Now(), events.PRIORITY_LOW}

		producer.outputChan <- uptimeEvent
		log.Println("Systeminfo polling done")
		<- tick
	}

}
