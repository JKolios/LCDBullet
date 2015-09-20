package pushbullet

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	wsURI  = "wss://stream.pushbullet.com:443/websocket/"
	apiURL = "https://api.pushbullet.com/v2/"
)

type APIClient struct {
	token   string
	control chan int
	Output  chan string
}

func NewClient(apiToken string) *APIClient {
	channel := make(chan int)
	output := make(chan string)
	return &APIClient{token: apiToken, control: channel, Output: output}
}

func (client *APIClient) StartMonitoring() {

	dialer := websocket.Dialer{}
	wsHeaders := http.Header{
		"Origin":                   {"http://jsone.system-ns.net"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	wsConnection, response, err := dialer.Dial(wsURI+client.token, wsHeaders)
	if err != nil {
		log.Fatal("Error opening websocket connection:" + err.Error())
	}
	log.Println(response)
	log.Println("Websocket connection appears to be up, monitoring")
	go wsMonitor(wsConnection, client.token, client.control, client.Output)
}

func (client *APIClient) StopMonitoring() {
	log.Println("Stopping wsMonitor...")
	client.control <- 0
	client.control <- 0

}

func wsMonitor(connection *websocket.Conn, token string, control chan int, output chan string) {
	//set up the message pump
	rawMessageChannel := make(chan map[string]interface{}, 10)
	lastPushCheck := time.Now()
	go messagePump(connection, rawMessageChannel, control)
	for {

		select {
		case <-control:
			{
				log.Println("wsMonitor Terminated")
				return
			}
		case message := <-rawMessageChannel:
			{
				log.Printf("Received JSON Message:Content %v\n", message)
				if message["type"] == "push" {
					log.Println("Got an ephemeral push, data:")
					log.Println(message)
				} else if message["type"] == "tickle" {
					log.Println("Got a tickle push, fetching body/bodies...")
					go getPushesSince(lastPushCheck, token, output)
					lastPushCheck = time.Now()
				} else {
					log.Println("Got a nop message, ignoring")
				}

			}
		}
	}

}

func messagePump(conn *websocket.Conn, messageChannel chan map[string]interface{}, control chan int) {
	messageMap := make(map[string]interface{})
	for {
		select {
		case <-control:
			{
				log.Println("Message Pump Terminated")
				conn.Close()
				return
			}
		default:
			err := conn.ReadJSON(&messageMap)
			log.Println("Parsing message")
			if err != nil {
				log.Println("Error while parsing message" + err.Error())
			}
			messageChannel <- messageMap
		}
	}
}

func getPushesSince(since time.Time, token string, output chan string) {

	httpClient := &http.Client{}
	requestUrl := apiURL + "pushes?active=true&modified_after=" + strconv.Itoa(int(since.Unix()))
	log.Println(requestUrl)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Println("Error constructing API request:" + err.Error())
		return

	}
	req.Header.Add("Access-Token", token)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending API request:" + err.Error())
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading JSON response:" + err.Error())
		return
	}

	responseStruct := make(map[string][]map[string]interface{})

	err = json.Unmarshal(response, &responseStruct)
	if err != nil {
		log.Println("Error unmarshalling JSON response:" + err.Error())
		return
	}
	log.Println(responseStruct)
	for _, message := range responseStruct["pushes"] {
		body, ok := message["body"].(string)
		if !ok  {
			messageId, ok :=message["iden"].(string)
			if !ok {
				log.Println("Got a malformed push, ignoring")
			}
			log.Println("Message " + messageId + " contains no body, ignoring")
		}else {
			output <- body
		}
	}
}
