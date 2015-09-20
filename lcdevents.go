package main

import (
	"github.com/Jkolios/goLcdEvents/lcd"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/Jkolios/goLcdEvents/pushbullet"
	"gopkg.in/yaml.v2"
	"bufio"
	"os"
	"log"
	"time"
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
	for scanner.Scan(){
		log.Println(scanner.Text())
		yaml.Unmarshal(scanner.Bytes(), &confObject)
	}
	return confObject
}

func displayOnLCD(input chan string, lcd *lcd.SharedDisplay) {
	for {
		lcd.DisplayMessage(<-input, 3 *time.Second)
	}

}

func main() {
	config := parseYAMLConf("conf.yml")
	log.Printf("Config:%v", config)

	var pinout []int
	for _,val := range(config["pinout"].([]interface{})) {
		pinout = append(pinout, val.(int))
	}
	display := lcd.NewDisplay(pinout, true)
	defer display.Close()

	display.DisplayMessage("Display initialized", 3*time.Second)

	client := pushbullet.NewClient(config["apiToken"].(string))
	client.StartMonitoring()
	go displayOnLCD(client.Output, display)
	for {
	}
}
