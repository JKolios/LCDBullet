package main

import (
	"log"
	"os"
	"os/signal"
	"container/list"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/consumers/httplog"
	"github.com/JKolios/goLcdEvents/consumers/lcd"
	"github.com/JKolios/goLcdEvents/consumers/wsclient"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/producers/bmp"
	"github.com/JKolios/goLcdEvents/producers/pushbullet"
	"github.com/JKolios/goLcdEvents/producers/systeminfo"
	"github.com/JKolios/goLcdEvents/producers/wunderground"
	"github.com/JKolios/goLcdEvents/producers/bitcoinaverage"
	_ "github.com/kidoman/embd/host/rpi"
	"time"
)

var highPriorityeventList = list.New()
var lowPriorityeventList = list.New()

var consumerMap =  map[string]events.Consumer{
	"lcd": &lcd.LCDConsumer{},
	"httplog": &httplog.HttpConsumer{},
	"wsclient": &wsclient.WebsocketConsumer{}}

var producerMap =  map[string]events.Producer{
	"pushbullet": &pushbullet.PushbulletProducer{},
	"bmp": &bmp.BMPProducer{},
	"systeminfo": &systeminfo.SystemInfoProducer{},
	"wunderground":&wunderground.WundergroundProducer{},
	"bitcoinaverage":&bitcoinaverage.BitcoinAverageProducer{}}

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
				}else{
					lowPriorityeventList.PushBack(incomingEvent)
				}


			}
		}
}


func consumerHub(done chan struct{}, consumerChans []chan events.Event) {

	var selectedEvent events.Event

		for {
			select {

			case <-done:
				log.Println("consumerHub halting")
				return

			default:

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

			}
		}
		}



func main() {
	config, err := conf.ParseJSONConf("conf.json")
	if err != nil {
		log.Fatalln("Error while parsing config: " + err.Error())
	}
	// Consumer Init

	var consumerChannels []chan events.Event

	done := make(chan struct{})

	for _, consumerName := range config.Consumers {
		consumerMap[consumerName].Initialize(config)
		newChan := make(chan events.Event, 100)
		consumerMap[consumerName].Start(done, newChan)
		consumerChannels = append(consumerChannels, newChan)
	}

	// Producer Init
	producerChan := make(chan events.Event)

	for _, producerName := range config.Producers {
		producerMap[producerName].Initialize(config)
		producerMap[producerName].Start(done, producerChan)
	}

	go producerHub(done, producerChan)
	go consumerHub(done, consumerChannels)

	controlChan := make(chan os.Signal)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	<- controlChan
	log.Println("Main got an OS signal, starting shutdown...")
	close(done)
	time.Sleep(30 * time.Second)

}