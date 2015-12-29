package pushbullet

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Push struct {
	Type     string  `json:"type"`
	Body     string  `json:"body"`
	Modified float64 `json:"modified"`
	Id       string  `json:"iden"`
}

type ListPushesResponse struct {
	Pushes []Push `json:"pushes"`
}

var httpClient *http.Client = &http.Client{}

func getPushesSince(since float64, token string) ListPushesResponse {

	requestUrl := apiURL + "pushes?active=true&modified_after=" + url.QueryEscape(strconv.FormatFloat(since, 'e', -1, 64))
	log.Println(requestUrl)
	var responseStruct ListPushesResponse

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.Println("Error constructing API request:" + err.Error())
		return responseStruct

	}
	req.Header.Add("Access-Token", token)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("Error sending API request:" + err.Error())
		return responseStruct
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading JSON response:" + err.Error())
		return responseStruct
	}

	log.Println(string(response))

	err = json.Unmarshal(response, &responseStruct)
	if err != nil {
		log.Println("Error unmarshalling JSON response:" + err.Error())
		return responseStruct
	}

	log.Println(responseStruct)
	return responseStruct
}
