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
	output     chan<- events.Event
	done       <-chan struct{}
}

type wsMessage struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
}

func (producer *PushbulletProducer) Initialize(config conf.Configuration) {
	producer.token = config.PushbulletApiToken

}

func (producer *PushbulletProducer) Start(done <-chan struct{}, EventOutput chan<- events.Event) {

	dialer := websocket.Dialer{}
	wsHeaders := http.Header{
		"Origin":                   {"http://jsone.system-ns.net"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsConnection, _, err := dialer.Dial(wsURI+producer.token, wsHeaders)
	if err != nil {
		log.Fatal("PushbulletProducer :Error opening websocket connection:" + err.Error())
	}
	log.Println("PushbulletProducer: Websocket connection appears to be up, monitoring")
	producer.connection = wsConnection
	producer.output = EventOutput
	producer.done = done
	go pushbulletMonitor(producer)
}

func pushbulletMonitor(producer *PushbulletProducer) {
	//set up the message pump
	wsMessageChannel := make(chan wsMessage)
	lastcheckTimestamp := float64(time.Now().Unix())
	go wsMessagePump(producer.done, producer.connection, wsMessageChannel)
	for {

		select {
		case <-producer.done:
			{
				log.Println("pushbulletMonitor Terminated")
				return
			}
		case message := <-wsMessageChannel:
			{
				log.Printf("Received JSON Message:Content %v\n", message)
				if message.Type == "tickle" {

					log.Println("Got a tickle push, fetching body/bodies...")
					ListPushesResponse := getPushesSince(lastcheckTimestamp, producer.token)

					for _, push := range ListPushesResponse.Pushes {

						pushEvent := events.Event{push.Body, "pushbullet", producer, time.Now(), events.PRIORITY_HIGH}
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

func wsMessagePump(done <-chan struct{}, conn *websocket.Conn, messageChannel chan wsMessage) {
	message := wsMessage{}
	for {
		select {
		case <-done:
			{
				log.Println("Message Pembdump Terminated")
				conn.Close()
				return
			}
		default:
			err := conn.ReadJSON(&message)
			// log.Println("Parsing message")
			if err != nil {
				log.Println("Error while parsing message" + err.Error())
			}
			messageChannel <- message
		}
	}
}
