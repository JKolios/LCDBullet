package utils

import (
	"log"
	"reflect"
)

func SliceContains(container []string, element string) bool {
	for _, a := range container {
		if reflect.DeepEqual(a, element) {
			return true
		}
	}
	return false
}

func LogErrorandExit(message string, err error) {
	if err != nil {
		log.Fatalln(message + err.Error())
	}
}

type Configuration struct {
	Producers  []string `json:"Producers"`
	Consumers  []string `json:"Consumers"`
	Pinout     []int    `json:"Pinout"`
	Bmpi2c     byte     `json:"Bmpi2c"`
	ApiToken   string   `json:"ApiToken"`
	BlPolarity bool     `json:"BlPolarity"`
}
