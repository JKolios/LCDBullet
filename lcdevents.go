package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"container/list"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/consumers/httplog"
	"github.com/JKolios/goLcdEvents/consumers/lcd"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/producers/bmp"
	"github.com/JKolios/goLcdEvents/producers/pushbullet"
	"github.com/JKolios/goLcdEvents/producers/systeminfo"
	"github.com/JKolios/goLcdEvents/producers/wunderground"
	"github.com/JKolios/goLcdEvents/producers/bitcoinaverage"
	"github.com/JKolios/goLcdEvents/utils"
	_ "github.com/kidoman/embd/host/rpi"
)

var highPriorityeventList = list.New()
var lowPriorityeventList = list.New()

func producerHub(producerChan chan events.Event, control chan os.Signal) {

	var incomingEvent events.Event

	ReceiveLoop:
		for {
			select {
			case incomingEvent = <-producerChan:
				log.Printf("producerHub got event: %+v\n", incomingEvent)
				// Inspect the event and handle according to type and priority
				if incomingEvent.Priority == events.PRIORITY_HIGH {
					highPriorityeventList.PushBack(incomingEvent)
				}else{
					lowPriorityeventList.PushBack(incomingEvent)
				}
			case <-control:
				log.Println("producerHub got an OS signal, enqueueing shutdown event")
				shutdownEvent := events.Event{false, "shutdown", nil, time.Now(), events.PRIORITY_IMMEDIATE}
				highPriorityeventList.PushFront(shutdownEvent)
				log.Println("producerHub halting")
				control <- syscall.SIGINT
				break ReceiveLoop

			}
		}
}


func consumerHub(consumerChans []chan events.Event, control chan os.Signal) {

	var selectedEvent events.Event

	SendLoop:
		for {
			if highPriorityeventList.Len() > 0 {
				selectedEvent = highPriorityeventList.Remove(highPriorityeventList.Front()).(events.Event)
			}else if lowPriorityeventList.Len() > 0 {
				selectedEvent = lowPriorityeventList.Remove(lowPriorityeventList.Front()).(events.Event)
			}else {
				continue
			}


			log.Printf("consumerHub selected event: %+v \n", selectedEvent)

			for _, consumerChan := range (consumerChans) {
				consumerChan <- selectedEvent
			}

			if selectedEvent.Type == "shutdown" {
				log.Println("consumerHub halting")
				control <- syscall.SIGINT
				break SendLoop
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
		newChan = make(chan events.Event)
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

	if utils.SliceContains(config.Producers, "bitcoinaverage") {
		Producers = append(Producers, &bitcoinaverage.BitcoinAverageProducer{})
	}

	producerChan := make(chan events.Event)

	for _, producer := range Producers {
		producer.Initialize(config)
		producer.Subscribe(producerChan)

	}

	controlChan := make(chan os.Signal)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	go producerHub(producerChan, controlChan)
	go consumerHub(consumerChannels, controlChan)

	controlChan <- <- controlChan
	log.Println("Main got an OS signal, bouncing...")
	<- controlChan
	<- controlChan
	time.Sleep(10 * time.Second)
}
