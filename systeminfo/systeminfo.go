package systeminfo

import (
	"log"
	"time"
		linuxproc "github.com/c9s/goprocinfo/linux"
)

func InitSystemMonitoring() chan string {

	output:= make(chan string)
	go pollSystemInfo(output, 10 * time.Second)
	return output
}

func pollSystemInfo(output chan string, every time.Duration) {
	for {
		uptime, err := linuxproc.ReadUptime("/proc/uptime")
		if err != nil {
			log.Fatal("Failed while getting uptime")
		}
		uptimeStr := "Uptime: " + time.Duration(time.Duration(int64(uptime.Total)) * time.Second).String()

		output<- uptimeStr
		time.Sleep(every)
	}

}