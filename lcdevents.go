package main

import (
	"container/list"
	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/consumers"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/inspector"
	"github.com/JKolios/goLcdEvents/producers"
	_ "github.com/kidoman/embd/host/rpi"
	"log"
	"os"
	"os/signal"
	"time"
)

var highPriorityeventList = list.New()
var lowPriorityeventList = list.New()

func producerHub(done chan struct{}, producerChan chan events.Event) {

	var incomingEvent events.Event

	for {
		select {
		case <-done:
			log.Println("producerHub halting")
			return

		case incomingEvent = <-producerChan:
			log.Printf("producerHub got event: %+v\n", incomingEvent)
			// Inspect the event and handle according to type and priority
			if incomingEvent.Priority == events.PRIORITY_HIGH {
				highPriorityeventList.PushBack(incomingEvent)
			} else {
				lowPriorityeventList.PushBack(incomingEvent)
			}

		}
	}
}

func consumerHub(done chan struct{}, consumerChans []chan<- events.Event) {

	var selectedEvent events.Event

	for {
		select {

		case <-done:
			log.Println("consumerHub halting")
			return

		default:

			if highPriorityeventList.Len() > 0 {
				selectedEvent = highPriorityeventList.Remove(highPriorityeventList.Front()).(events.Event)
			} else if lowPriorityeventList.Len() > 0 {
				selectedEvent = lowPriorityeventList.Remove(lowPriorityeventList.Front()).(events.Event)
			} else {
				continue
			}

			log.Printf("consumerHub selected event: %+v \n", selectedEvent)

			for _, consumerChan := range consumerChans {
				consumerChan <- selectedEvent
			}

		}
	}
}

func main() {
	config, err := conf.ParseJSONFile("conf.json")
	if err != nil {
		log.Fatalln("Error while parsing config: " + err.Error())
	}

	// Consumer Init
	var consumerChannels []chan<- events.Event

	done := make(chan struct{})

	for _, consumerName := range config.Consumers {
		consumer := consumers.NewConsumer(consumerName)
		consumer.Initialize(config)
		consumerChannels = append(consumerChannels, consumer.Start(done))
	}

	// Producer Init
	producerChan := make(chan events.Event)

	for _, producerName := range config.Producers {
		producer := producers.NewProducer(producerName)
		producer.Initialize(config)
		producer.Start(done, producerChan)
	}

	go producerHub(done, producerChan)
	go consumerHub(done, consumerChannels)

	go inspector.ListReport(lowPriorityeventList, "Low Priority", time.Minute*10)
	go inspector.ListReport(highPriorityeventList, "High Priority", time.Minute*10)

	go inspector.ListCleaner(lowPriorityeventList, "Low Priority", time.Minute*10)
	go inspector.ListCleaner(highPriorityeventList, "High Priority", time.Minute*10)

	controlChan := make(chan os.Signal)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	<-controlChan
	log.Println("Main got an OS signal, starting shutdown...")
	close(done)
	time.Sleep(20 * time.Second)

}
