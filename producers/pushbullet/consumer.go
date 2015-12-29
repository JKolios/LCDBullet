package pushbullet

import (
	"log"
	"net/http"
	"time"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/gorilla/websocket"
)

const (
	wsURI  = "wss://stream.pushbullet.com:443/websocket/"
	apiURL = "https://api.pushbullet.com/v2/"
)

type PushbulletProducer struct {
	connection *websocket.Conn
	token      string
	control    chan int
	output     chan events.Event
}

type wsMessage struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
}

func (producer *PushbulletProducer) Initialize(config conf.Configuration) {
	producer.control = make(chan int)
	producer.token = config.ApiToken

}

func (producer *PushbulletProducer) Subscribe(producerChan chan events.Event) {

	dialer := websocket.Dialer{}
	wsHeaders := http.Header{
		"Origin":                   {"http://jsone.system-ns.net"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsConnection, response, err := dialer.Dial(wsURI+producer.token, wsHeaders)
	if err != nil {
		log.Fatal("Error opening websocket connection:" + err.Error())
	}
	log.Println(response)
	log.Println("Websocket connection appears to be up, monitoring")
	producer.connection = wsConnection
	producer.output = producerChan
	go pushbulletMonitor(producer)
}

func (producer *PushbulletProducer) StopMonitoring() {
	log.Println("Stopping wsMonitor...")
	producer.control <- 0
	producer.control <- 0

}

func pushbulletMonitor(producer *PushbulletProducer) {
	//set up the message pump
	wsMessageChannel := make(chan wsMessage)
	lastcheckTimestamp := float64(time.Now().Unix())
	go wsMessagePump(producer.connection, wsMessageChannel, producer.control)
	for {

		select {
		case <-producer.control:
			{
				log.Println("wsMonitor Terminated")
				return
			}
		case message := <-wsMessageChannel:
			{
				log.Printf("Received JSON Message:Content %v\n", message)
				if message.Type == "tickle" {

					log.Println("Got a tickle push, fetching body/bodies...")
					ListPushesResponse := getPushesSince(lastcheckTimestamp, producer.token)

					for _, push := range ListPushesResponse.Pushes {

						pushEvent := events.Event{push.Body, "pushbullet", producer}
						producer.output <- pushEvent
						lastcheckTimestamp = push.Modified

					}

				} else if message.Type == "push" {
					log.Println("Got an ephemeral push, data:")
					log.Println(message)

				} else {
					log.Println("Got a nop message, ignoring")
				}
			}

		}
	}
}

func wsMessagePump(conn *websocket.Conn, messageChannel chan wsMessage, control chan int) {
	message := wsMessage{}
	for {
		select {
		case <-control:
			{
				log.Println("Message Pump Terminated")
				conn.Close()
				return
			}
		default:
			err := conn.ReadJSON(&message)
			log.Println("Parsing message")
			if err != nil {
				log.Println("Error while parsing message" + err.Error())
			}
			messageChannel <- message
		}
	}
}
