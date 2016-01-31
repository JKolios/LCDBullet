package wunderground

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
	apiURL = "http://api.wunderground.com/api/"
)

type apiResponse struct {
	Obs Observation `json:"current_observation"`
}

type Observation struct {
	Weather     string  `json:"weather"`
	Temp_c      float32 `json:"temp_c"`
	Feelslike_c string  `json:"feelslike_c"`
}

var httpClient *http.Client = &http.Client{}

func ProducerSetupFuction(producer *producers.GenericProducer, config map[string]interface{}) {
	producer.RuntimeObjects["token"] = config["wundergroundApiToken"].(string)
	producer.RuntimeObjects["location"] = config["wundergroundLocation"].(string)

}

func ProducerWaitFunction(producer *producers.GenericProducer) {
	time.Sleep(30 * time.Second)
}

func ProducerRunFuction(producer *producers.GenericProducer) events.Event {

	eventPayload, eventPriority := getCurrentConditions(producer.RuntimeObjects)
	return events.Event{eventPayload, producer.Name, time.Now(), eventPriority}
}

func getCurrentConditions(config map[string]interface{}) (string, int) {

	requestUrl := apiURL + config["token"].(string) + "/conditions/q/" + config["location"].(string) + ".json"
	var responseStruct apiResponse

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Println("Error constructing API request:" + err.Error())
		return "Wunderground Error", events.PRIORITY_LOW

	}
	req.Header.Add("Access-Token", config["token"].(string))
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending API request:" + err.Error())
		return "Wunderground Error", events.PRIORITY_LOW
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading JSON response:" + err.Error())
		return "Wunderground Error", events.PRIORITY_LOW
	}

	err = json.Unmarshal(response, &responseStruct)
	if err != nil {
		log.Println("Error unmarshalling JSON response:" + err.Error())
		return "Wunderground Error", events.PRIORITY_LOW
	}

	finalMessage := fmt.Sprintf("%v Temp:%v Feels like:%v", responseStruct.Obs.Weather, responseStruct.Obs.Temp_c, responseStruct.Obs.Feelslike_c)
	return finalMessage, events.PRIORITY_LOW
}
