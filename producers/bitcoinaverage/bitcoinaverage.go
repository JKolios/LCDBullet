package bitcoinaverage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/JKolios/EventsToGo/events"
	"github.com/JKolios/EventsToGo/producers"
	"time"
)

const (
	API_URL = "https://api.bitcoinaverage.com/ticker/global/"
)

type apiResponse struct {
	Ask    float32 `json:"ask"`
	Bid    float32 `json:"bid"`
	Last   float32 `json:"last"`
	Avg24h float32 `json:"24h_avg"`
}

type BitcoinAverageProducer struct {
	producers.GenericProducer
	requestEndpoint string
}

func ProducerSetupFuction(producer *producers.GenericProducer, config map[string]interface{}) {
	producer.RuntimeObjects["requestEndpoint"] = API_URL + config["BitcoinAverageCurrency"].(string)

}

func ProducerWaitFunction(producer *producers.GenericProducer) {
	time.Sleep(10 * time.Second)
}

func ProducerRunFuction(producer *producers.GenericProducer) events.Event {

	var responseStruct apiResponse

	req, err := http.NewRequest("GET", producer.RuntimeObjects["requestEndpoint"].(string), nil)
	if err != nil {
		log.Println("Error constructing API request:" + err.Error())
		return events.Event{}

	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending API request:" + err.Error())
		return events.Event{}
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading JSON response:" + err.Error())
		return events.Event{}
	}

	err = json.Unmarshal(response, &responseStruct)
	if err != nil {
		log.Println("Error unmarshalling JSON response:" + err.Error())
		return events.Event{}
	}

	finalMessage := fmt.Sprintf("Bitcoin Global Average: Ask: %v Bid:%v Last:%v 24H Average: %v", responseStruct.Ask, responseStruct.Bid, responseStruct.Last, responseStruct.Avg24h)

	return events.Event{finalMessage, producer.Name, time.Now(), events.PRIORITY_LOW}

}
