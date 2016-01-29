package pushbullet

import (
	"log"
	"time"

	"github.com/JKolios/EventsToGo/events"
	"github.com/JKolios/EventsToGo/producers"
	"github.com/gorilla/websocket"
)

const (
	wsURI  = "wss://stream.pushbullet.com:443/websocket/"
	apiURL = "https://api.pushbullet.com/v2/"
)

type PushbulletProducer struct {
	connection *websocket.Conn
	token      string
}

type wsMessage struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype"`
}

func wsMessagePump(producer *producers.GenericProducer) {
	message := wsMessage{}
	for {
		err := producer.RuntimeObjects["connection"].(*websocket.Conn).ReadJSON(&message)
		if err != nil {
			log.Println("Error while parsing message" + err.Error())
		}

		producer.RuntimeObjects["wsMessageChannel"].(chan wsMessage) <- message
	}
}

func ProducerSetupFuction(producer *producers.GenericProducer, config map[string]interface{}) {

	dialer := websocket.Dialer{}

	wsConnection, _, err := dialer.Dial(wsURI+config["PushbulletApiToken"].(string), nil)
	if err != nil {
		log.Fatal("PushbulletProducer: Error opening websocket connection:" + err.Error())
	}

	producer.RuntimeObjects["PushbulletApiToken"] = config["PushbulletApiToken"].(string)
	producer.RuntimeObjects["connection"] = wsConnection
	producer.RuntimeObjects["wsMessageChannel"] = make(chan wsMessage)
	producer.RuntimeObjects["lastcheckTimestamp"] = float64(time.Now().Unix())

	go wsMessagePump(producer)

}

func ProducerWaitFunction(producer *producers.GenericProducer) {
	for {
		message := <-producer.RuntimeObjects["wsMessageChannel"].(chan wsMessage)

		if message.Type == "tickle" {
			return

		}
	}
}

func ProducerRunFuction(producer *producers.GenericProducer) events.Event {

	ListPushesResponse := getPushesSince(producer.RuntimeObjects["lastcheckTimestamp"].(float64), producer.RuntimeObjects["PushbulletApiToken"].(string))
	producer.RuntimeObjects["lastcheckTimestamp"] = ListPushesResponse.Pushes[0].Modified

	pushbulletUpdates := "Pushbullet updates: "
	for _, push := range ListPushesResponse.Pushes {
		pushbulletUpdates += push.Body + " "
	}

	pushEvent := events.Event{pushbulletUpdates, "pushbullet", time.Now(), events.PRIORITY_HIGH}
	return pushEvent
}

func ProducerStopFunction(producer *producers.GenericProducer) {
	producer.RuntimeObjects["connection"].(*websocket.Conn).Close()
}
