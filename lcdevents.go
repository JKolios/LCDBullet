package main

import (
	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/eventqueue"
	_ "github.com/kidoman/embd/host/rpi"
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

	eventQueue := eventqueue.NewQueue(config.Consumers, config.Producers, config)
	eventQueue.Start()

	controlChan := make(chan os.Signal)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	<-controlChan
	log.Println("Main got an OS signal, starting shutdown...")
	eventQueue.Stop()

	time.Sleep(20 * time.Second)

}
