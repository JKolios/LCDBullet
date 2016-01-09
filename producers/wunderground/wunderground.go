package wunderground

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
	apiURL = "http://api.wunderground.com/api/"
)

type WundergroundProducer struct {
	token    string
	location string
	output   chan events.Event
}

type apiResponse struct {
	Obs Observation `json:"current_observation"`
}

type Observation struct {
	Weather     string  `json:"weather"`
	Temp_c      float32 `json:"temp_c"`
	Feelslike_c string  `json:"feelslike_c"`
}

var httpClient *http.Client = &http.Client{}

func (producer *WundergroundProducer) Initialize(config conf.Configuration) {
	producer.token = config.WundergroundApiToken
	producer.location = config.WundergroundLocation

}

func (producer *WundergroundProducer) Subscribe(producerChan chan events.Event) {

	producer.output = producerChan
	go pollWunderground(producer, 30*time.Second)
}

func pollWunderground(producer *WundergroundProducer, every time.Duration) {
	tick := time.Tick(every)
	for {
		log.Println("Starting wunderground polling")
		conditions := getCurrentConditions(producer.token, producer.location)

		finalMessage := fmt.Sprintf("%v Temp:%v Feels like:%v", conditions.Weather, conditions.Temp_c, conditions.Feelslike_c)
		finalEvent := events.Event{finalMessage, "wunderground", producer, time.Now(), events.PRIORITY_HIGH}

		producer.output <- finalEvent
		log.Println("Wunderground polling done")
		<- tick
	}
}

func getCurrentConditions(token, location string) Observation {

	requestUrl := apiURL + token + "/conditions/q/" + location + ".json"
	var responseStruct apiResponse

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Println("Error constructing API request:" + err.Error())
		return Observation{}

	}
	req.Header.Add("Access-Token", token)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending API request:" + err.Error())
		return Observation{}
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading JSON response:" + err.Error())
		return Observation{}
	}

	err = json.Unmarshal(response, &responseStruct)
	if err != nil {
		log.Println("Error unmarshalling JSON response:" + err.Error())
		return Observation{}
	}

	return responseStruct.Obs
}
