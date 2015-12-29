package httplog

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
)

var httpContent = make(chan string)

type HttpConsumer struct {
	listenAddress string
	endpoint      string
}

func (consumer *HttpConsumer) Initialize(config conf.Configuration) {
	// Config Parsing
	consumer.listenAddress = config.ListenAddress
	consumer.endpoint = config.Endpoint
}

func (consumer *HttpConsumer) Register(EventInput chan events.Event) {

	// Input Monitor Goroutine Startup
	go monitorInput(EventInput)

	//HTTP Handler Startup
	http.HandleFunc(consumer.endpoint, httpPollHandler)
	go http.ListenAndServe(consumer.listenAddress, nil)
	log.Println("HTTP log: started, listening at " + consumer.listenAddress + consumer.endpoint)
}

func monitorInput(inputChan chan events.Event) {
	for {
		incomingEvent := <-inputChan
		eventAsText := fmt.Sprintf("%s:%s\n", incomingEvent.Type, incomingEvent.Payload.(string))
		httpContent <- eventAsText
	}

}

func httpPollHandler(w http.ResponseWriter, req *http.Request) {

	writeFlusher, ok := w.(http.Flusher)
	if !ok {
		panic("HTTP log: expected http.ResponseWriter to be an http.Flusher")
	}

	for {
		incomingText := <-httpContent
		log.Print("HTTP log: incoming text:" + incomingText)
		_, err := io.WriteString(w, incomingText)
		if err != nil {
			log.Println("HTTP log: Failed to append to an HTTP Response:" + err.Error())
			break
		}
		writeFlusher.Flush()
		log.Println("HTTP log: Flushing response")
	}
}
