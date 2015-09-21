package bmp

import (
		"strconv"
		"fmt"
		"log"
		"time"
		"github.com/kidoman/embd"
		"github.com/kidoman/embd/sensor/bmp085"
)

func InitBMP085(i2cBus byte) chan string {
	bus := embd.NewI2CBus(i2cBus)
	bmp := bmp085.New(bus)
	output:= make(chan string)
	go pollBMP085(bmp, output, 3 * time.Second)
	return output
}

func pollBMP085(sensor *bmp085.BMP085, output chan string, every time.Duration) {
	for {
		temperature, err := sensor.Temperature()
	if err != nil {
		log.Println("Cannot get temperature:" + err.Error())
		return
	}
	pressure, err := sensor.Pressure()
	if err != nil {
		log.Println("Cannot get pressure:" + err.Error())
		return
	}
	altitude, err := sensor.Altitude()
	if err != nil {
		log.Println("Cannot get altitude:" + err.Error())
		return
	}
	tempStr := strconv.FormatFloat(temperature, 'f', 2, 64)
	pressStr := strconv.Itoa(pressure)
	altStr := strconv.FormatFloat(altitude, 'f', 2, 64)

	finalMessage := fmt.Sprintf("%v T:%v P:%v A:%v", time.Now().Format(time.Kitchen), tempStr, pressStr, altStr)
	output<- finalMessage
		time.Sleep(every)

	}
}