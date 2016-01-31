package main

import (
	"github.com/JKolios/EventsToGo"
	"github.com/JKolios/EventsToGo/consumers"
	"github.com/JKolios/EventsToGo/producers"
	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/consumers/lcd"
	"github.com/JKolios/goLcdEvents/consumers/wsclient"
	"github.com/JKolios/goLcdEvents/producers/bitcoinaverage"
	"github.com/JKolios/goLcdEvents/producers/bmp"
	"github.com/JKolios/goLcdEvents/producers/pushbullet"
	"github.com/JKolios/goLcdEvents/producers/systeminfo"
	"github.com/JKolios/goLcdEvents/producers/wunderground"
	"github.com/JKolios/goLcdEvents/utils"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	config, err := conf.ParseJSONFile("conf.json")
	if err != nil {
		log.Fatalln("Error while parsing config: " + err.Error())
	}

	eventTTL := time.Minute * 5
	eventQueue := EventsToGo.NewQueue(&eventTTL)

	if utils.SliceContainsString(config["consumers"].([]interface{}), "lcd") {
		lcdClientConsumer := consumers.NewGenericConsumer("lcd", config)
		lcdClientConsumer.RegisterFunctions(lcd.SetupFunction, lcd.RunFunction, nil)

		eventQueue.AddConsumer(lcdClientConsumer)
	}

	if utils.SliceContainsString(config["consumers"].([]interface{}), "wsclient") {
		wsClientConsumer := consumers.NewGenericConsumer("wsclient", config)
		wsClientConsumer.RegisterFunctions(wsclient.SetupFunction, wsclient.RunFunction, nil)

		eventQueue.AddConsumer(wsClientConsumer)
	}

	if utils.SliceContainsString(config["producers"].([]interface{}), "wunderground") {
		wundergroundProducer := producers.NewGenericProducer("wunderground", config)
		wundergroundProducer.RegisterFunctions(wunderground.ProducerSetupFuction, wunderground.ProducerRunFuction, wunderground.ProducerWaitFunction, nil)

		eventQueue.AddProducer(wundergroundProducer)
	}

	if utils.SliceContainsString(config["producers"].([]interface{}), "bitcoinaverage") {
		bitcoinaverageProducer := producers.NewGenericProducer("bitcoinaverage", config)
		bitcoinaverageProducer.RegisterFunctions(bitcoinaverage.ProducerSetupFuction, bitcoinaverage.ProducerRunFuction, bitcoinaverage.ProducerWaitFunction, nil)

		eventQueue.AddProducer(bitcoinaverageProducer)
	}

	if utils.SliceContainsString(config["producers"].([]interface{}), "bitcoinaverage") {
		pushbulletProducer := producers.NewGenericProducer("pushbullet", config)
		pushbulletProducer.RegisterFunctions(pushbullet.ProducerSetupFuction, pushbullet.ProducerRunFuction, pushbullet.ProducerWaitFunction, pushbullet.ProducerStopFunction)

		eventQueue.AddProducer(pushbulletProducer)
	}

	if utils.SliceContainsString(config["producers"].([]interface{}), "systeminfo") {
		systeminfoProducer := producers.NewGenericProducer("systeminfo", config)
		systeminfoProducer.RegisterFunctions(systeminfo.ProducerSetupFuction, systeminfo.ProducerRunFuction, systeminfo.ProducerWaitFunction, nil)

		eventQueue.AddProducer(systeminfoProducer)
	}

	if utils.SliceContainsString(config["producers"].([]interface{}), "bmp") {
		bmpProducer := producers.NewGenericProducer("bmp", config)
		bmpProducer.RegisterFunctions(bmp.ProducerSetupFuction, bmp.ProducerRunFuction, bmp.ProducerWaitFunction, nil)

		eventQueue.AddProducer(bmpProducer)
	}

	eventQueue.Start()

	controlChan := make(chan os.Signal)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	<-controlChan
	log.Println("Main got an OS signal, starting shutdown...")
	eventQueue.Stop()

}
