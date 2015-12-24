package main

import (
	"encoding/json"
	_ "github.com/JKolios/goLcdEvents/Godeps/_workspace/src/github.com/kidoman/embd/host/rpi"
	"github.com/JKolios/goLcdEvents/consumers/lcd"
	"github.com/JKolios/goLcdEvents/producers/bmp"
	"github.com/JKolios/goLcdEvents/producers/pushbullet"
	"github.com/JKolios/goLcdEvents/producers/systeminfo"
	"github.com/JKolios/goLcdEvents/utils"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Producer interface {
	Initialize(config utils.Configuration)
	Subscribe(chan string)
	Terminate()
}

type Consumer interface {
	Initialize(config utils.Configuration)
	Consume(mType string, message string)
	Terminate()
}

func parseJSONConf(filename string) (utils.Configuration, error) {

	var confObject utils.Configuration
	confFile, err := ioutil.ReadFile(filename)

	err = json.Unmarshal(confFile, &confObject)
	return confObject, err
}

func messageHub(producerChan chan string, Consumers []Consumer, control chan os.Signal) {
	for {
		select {
		case incoming := <-producerChan:
			for _, consumer := range Consumers {
				consumer.Consume("other", incoming)
			}
		case <-control:
			for _, consumer := range Consumers {
				consumer.Terminate()
			}
			control <- syscall.SIGINT
		}

	}
}

func main() {
	config, err := parseJSONConf("conf.json")
	if err != nil {
		log.Fatalln("Error while parsing config: " + err.Error())
	}
	log.Printf("Config:%+v", config)

	// Consumer Init

	var Consumers []Consumer

	if utils.SliceContains(config.Consumers, "lcd") {

		lcdConsumer := lcd.NewLCDConsumer()
		lcdConsumer.Initialize(config)
		lcdConsumer.Consume("other", "Display Initialized")
		Consumers = append(Consumers, lcdConsumer)
		defer lcdConsumer.Terminate()

	}

	// Producer Init

	var producers []Producer
	producerChan := make(chan string)

	if utils.SliceContains(config.Producers, "bmp") {
		bmpProducer := bmp.NewBMPProducer()
		bmpProducer.Initialize(config)
		bmpProducer.Subscribe(producerChan)
		producers = append(producers, bmpProducer)
	}

	if utils.SliceContains(config.Producers, "pushbullet") {
		pushbulletProducer := pushbullet.NewPushbulletProducer()
		pushbulletProducer.Initialize(config)
		pushbulletProducer.Subscribe(producerChan)
		producers = append(producers, pushbulletProducer)
	}

	if utils.SliceContains(config.Producers, "systeminfo") {
		sysinfoProducer := systeminfo.NewSystemInfoProducer()
		sysinfoProducer.Initialize(config)
		sysinfoProducer.Subscribe(producerChan)
		producers = append(producers, sysinfoProducer)
	}

	controlChan := make(chan os.Signal, 1)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	go messageHub(producerChan, Consumers, controlChan)

	<-controlChan
}
