package main

import (
	"bufio"
	"github.com/Jkolios/goLcdEvents/lcd"
	"github.com/Jkolios/goLcdEvents/pushbullet"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bmp085"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"time"
	"strconv"
)

func parseYAMLConf(filename string) map[string]interface{} {

	confFile, err := os.Open("conf.yml")
	if err != nil {
		log.Println("Cannot open conf file:" + err.Error())
		return nil
	}
	defer confFile.Close()

	scanner := bufio.NewScanner(confFile)
	confObject := make(map[string]interface{})
	for scanner.Scan() {
		log.Println(scanner.Text())
		yaml.Unmarshal(scanner.Bytes(), &confObject)
	}
	return confObject
}

func displayOnLCD(input chan string, lcd *lcd.SharedDisplay) {
	for {
		lcd.DisplayMessage(<-input, 5*time.Second, true)
	}

}

func main() {
	config := parseYAMLConf("conf.yml")
	log.Printf("Config:%v", config)

	var pinout []int
	for _, val := range config["pinout"].([]interface{}) {
		pinout = append(pinout, val.(int))
	}
	display := lcd.NewDisplay(pinout, true)
	defer display.Close()

	display.DisplayMessage("Display initialized", 3*time.Second, true)

	client := pushbullet.NewClient(config["apiToken"].(string))
	client.StartMonitoring()
	go displayOnLCD(client.Output, display)
	bus := embd.NewI2CBus(1)
	bmp := bmp085.New(bus)
	temperature, err := bmp.Temperature()
	if err != nil {
		log.Println("Cannot get temperature:" + err.Error())
		return
	}
	pressure, err := bmp.Pressure()
	if err != nil {
		log.Println("Cannot get pressure:" + err.Error())
		return
	}
	altitude, err := bmp.Altitude()
	if err != nil {
		log.Println("Cannot get altitude:" + err.Error())
		return
	}
	client.Output<- strconv.FormatFloat(temperature, 'f', 2, 64)
	client.Output<- strconv.Itoa(pressure)
	client.Output<- strconv.FormatFloat(altitude, 'f', 2, 64)
	for {}
}
