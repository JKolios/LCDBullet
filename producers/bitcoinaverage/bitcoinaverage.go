package bitcoinaverage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
)

const (
	apiURL = "https://api.bitcoinaverage.com/ticker/global/"
)

type BitcoinAverageProducer struct {
	currencySymbol string
	output   chan events.Event
}

type apiResponse struct {

	Ask float32 `json:"ask"`
	Bid float32 `json:"bid"`
	Last float32 `json:"last"`
	Avg24h float32 `json:"24h_avg"`

}

var httpClient *http.Client = &http.Client{}

func (producer *BitcoinAverageProducer) Initialize(config conf.Configuration) {
	producer.currencySymbol = config.BitcoinAverageCurrency

}

func (producer *BitcoinAverageProducer) Subscribe(producerChan chan events.Event) {

	producer.output = producerChan
	go pollBitcoinAverage(producer, 30*time.Second)
}

func pollBitcoinAverage(producer *BitcoinAverageProducer, every time.Duration) {
	tick := time.Tick(every)
	for {
		log.Println("Starting bitcoinaverage polling")
		averages := getCurrentBTCAverages(producer.currencySymbol)

		finalMessage := fmt.Sprintf("Bitcoin Global Average: Ask: %v Bid:%v Last:%v 24H Average: %v", averages.Ask, averages.Bid, averages.Last, averages.Avg24h)
		finalEvent := events.Event{finalMessage, "bitcoinaverage", producer, time.Now(), events.PRIORITY_LOW}

		producer.output <- finalEvent
		log.Println("Wunderground polling done")
		<- tick
	}
}

func getCurrentBTCAverages(currencySymbol string) apiResponse {

	requestUrl := apiURL + currencySymbol
	var responseStruct apiResponse

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Println("Error constructing API request:" + err.Error())
		return apiResponse{}

	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending API request:" + err.Error())
		return apiResponse{}
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading JSON response:" + err.Error())
		return apiResponse{}
	}

	err = json.Unmarshal(response, &responseStruct)
	if err != nil {
		log.Println("Error unmarshalling JSON response:" + err.Error())
		return apiResponse{}
	}

	return responseStruct
}
