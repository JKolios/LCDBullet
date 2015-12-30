package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	Producers            []string `json:"Producers"`
	Consumers            []string `json:"Consumers"`
	Pinout               []int    `json:"Pinout"`
	Bmpi2c               byte     `json:"Bmpi2c"`
	PushbulletApiToken   string   `json:"pushbulletApiToken"`
	WundergroundApiToken string   `json:"wundergroundApiToken"`
	WundergroundLocation string   `json:"wundergroundLocation"`
	BlPolarity           bool     `json:"BlPolarity"`
	ListenAddress        string   `json:"ListenAddress"`
	Endpoint             string   `json:"Endpoint"`
}

func ParseJSONConf(filename string) (Configuration, error) {

	var confObject Configuration
	confFile, err := ioutil.ReadFile(filename)

	err = json.Unmarshal(confFile, &confObject)
	return confObject, err
}
