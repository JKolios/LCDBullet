package utils

import (
	"log"
	"reflect"
)

func LogErrorandExit(message string, err error) {
	if err != nil {
		log.Fatalln(message + err.Error())
	}
}

func SliceContainsString(container []interface{}, element string) bool {
	for _, a := range container {
		if reflect.DeepEqual(a.(string), element) {
			return true
		}
	}
	return false
}
