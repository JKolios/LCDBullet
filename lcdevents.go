<<<<<<< Updated upstream
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/consumers/httplog"
	"github.com/JKolios/goLcdEvents/consumers/lcd"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/producers/bmp"
	"github.com/JKolios/goLcdEvents/producers/pushbullet"
	"github.com/JKolios/goLcdEvents/producers/systeminfo"
	"github.com/JKolios/goLcdEvents/utils"
	_ "github.com/kidoman/embd/host/rpi"
)

func eventHub(producerChan chan events.Event, consumerChans []chan events.Event, control chan os.Signal) {
	var incoming events.Event
	for {
		log.Println("Hub iteration")

		select {
		case incoming = <-producerChan:
			log.Println(incoming)
			for _, consumerChan := range consumerChans {
				consumerChan <- incoming
			}
		case <-control:

			shutdownEvent := events.Event{false, "shutdown", nil}
			for _, consumerChan := range consumerChans {
				consumerChan <- shutdownEvent
			}
			producerChan <- shutdownEvent
			control <- syscall.SIGINT

		}
	}
}

func main() {
	config, err := conf.ParseJSONConf("conf.json")
	if err != nil {
		log.Fatalln("Error while parsing config: " + err.Error())
	}
	log.Printf("Config:%+v", config)

	// Consumer Init

	var Consumers []events.Consumer

	if utils.SliceContains(config.Consumers, "lcd") {
		Consumers = append(Consumers, &lcd.LCDConsumer{})
	}

	if utils.SliceContains(config.Consumers, "httplog") {
		Consumers = append(Consumers, &httplog.HttpConsumer{})
	}

	var consumerChannels []chan events.Event
	var newChan chan events.Event

	for _, consumer := range Consumers {
		consumer.Initialize(config)
		newChan = make(chan events.Event, 100)
		consumer.Register(newChan)
		consumerChannels = append(consumerChannels, newChan)
	}

	// Producer Init

	var Producers []events.Producer

	if utils.SliceContains(config.Producers, "pushbullet") {
		Producers = append(Producers, &pushbullet.PushbulletProducer{})
	}

	if utils.SliceContains(config.Producers, "bmp") {
		Producers = append(Producers, &bmp.BMPProducer{})
	}

	if utils.SliceContains(config.Producers, "systeminfo") {
		Producers = append(Producers, &systeminfo.SystemInfoProducer{})
	}

	producerChan := make(chan events.Event, 100)

	for _, producer := range Producers {
		producer.Initialize(config)
		producer.Subscribe(producerChan)

	}

	controlChan := make(chan os.Signal, 1)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	go eventHub(producerChan, consumerChannels, controlChan)

	<-controlChan
}
=======
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/consumers/httplog"
	"github.com/JKolios/goLcdEvents/consumers/lcd"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/producers/bmp"
	"github.com/JKolios/goLcdEvents/producers/pushbullet"
	"github.com/JKolios/goLcdEvents/producers/systeminfo"
	"github.com/JKolios/goLcdEvents/producers/wunderground"
	"github.com/JKolios/goLcdEvents/utils"
	_ "github.com/kidoman/embd/host/rpi"
)

func eventHub(producerChan chan events.Event, consumerChans []chan events.Event, control chan os.Signal) {
	var incoming events.Event
	for {
		log.Println("Hub iteration")

		select {
		case incoming = <-producerChan:
			log.Println(incoming)
			for _, consumerChan := range consumerChans {
				consumerChan <- incoming
			}
		case <-control:

			shutdownEvent := events.Event{false, "shutdown", nil}
			for _, consumerChan := range consumerChans {
				consumerChan <- shutdownEvent
			}
			producerChan <- shutdownEvent
			control <- syscall.SIGINT

		}
	}
}

func main() {
	config, err := conf.ParseJSONConf("conf.json")
	if err != nil {
		log.Fatalln("Error while parsing config: " + err.Error())
	}
	log.Printf("Config:%+v", config)

	// Consumer Init

	var Consumers []events.Consumer

	if utils.SliceContains(config.Consumers, "lcd") {
		Consumers = append(Consumers, &lcd.LCDConsumer{})
	}

	if utils.SliceContains(config.Consumers, "httplog") {
		Consumers = append(Consumers, &httplog.HttpConsumer{})
	}

	var consumerChannels []chan events.Event
	var newChan chan events.Event

	for _, consumer := range Consumers {
		consumer.Initialize(config)
		newChan = make(chan events.Event, 100)
		consumer.Register(newChan)
		consumerChannels = append(consumerChannels, newChan)
	}

	// Producer Init

	var Producers []events.Producer

	if utils.SliceContains(config.Producers, "pushbullet") {
		Producers = append(Producers, &pushbullet.PushbulletProducer{})
	}

	if utils.SliceContains(config.Producers, "bmp") {
		Producers = append(Producers, &bmp.BMPProducer{})
	}

	if utils.SliceContains(config.Producers, "systeminfo") {
		Producers = append(Producers, &systeminfo.SystemInfoProducer{})
	}

	if utils.SliceContains(config.Producers, "wunderground") {
		Producers = append(Producers, &wunderground.WundergroundProducer{})
	}

	producerChan := make(chan events.Event, 100)

	for _, producer := range Producers {
		producer.Initialize(config)
		producer.Subscribe(producerChan)

	}

	controlChan := make(chan os.Signal, 1)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	go eventHub(producerChan, consumerChannels, controlChan)

	<-controlChan
}
>>>>>>> Stashed changes
