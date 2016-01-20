package utils

import (
	"log"
)

func LogErrorandExit(message string, err error) {
	if err != nil {
		log.Fatalln(message + err.Error())
	}
}
