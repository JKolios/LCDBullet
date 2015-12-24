package systeminfo

import (
	"github.com/JKolios/goLcdEvents/utils"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"log"
	"time"
)

type SystemInfoProducer struct {
	outputChan chan string
}

func NewSystemInfoProducer() *SystemInfoProducer {
	return &SystemInfoProducer{}
}

func (producer *SystemInfoProducer) Initialize(config utils.Configuration) {
	//Dummy
	return
}

func (producer *SystemInfoProducer) Subscribe(producerChan chan string) {
	producer.outputChan = producerChan
	go pollSystemInfo(producer.outputChan, time.Second*10)
}

func pollSystemInfo(output chan string, every time.Duration) {
	for {
		uptime, err := linuxproc.ReadUptime("/proc/uptime")
		if err != nil {
			log.Fatal("Failed while getting uptime")
		}
		uptimeStr := "Uptime: " + time.Duration(time.Duration(int64(uptime.Total))*time.Second).String()

		output <- uptimeStr
		time.Sleep(every)
	}

}

func (consumer *SystemInfoProducer) Terminate() {
}
