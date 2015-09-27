package main

import (
	"bufio"
	"github.com/Jkolios/goLcdEvents/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/Jkolios/goLcdEvents/bmp"
	"github.com/Jkolios/goLcdEvents/lcd"
	"github.com/Jkolios/goLcdEvents/pushbullet"
	"github.com/Jkolios/goLcdEvents/systeminfo"
	_ "github.com/kidoman/embd/host/rpi"
	"log"
	"os"
	"os/signal"
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
	for scanner.Scan() {
		log.Println(scanner.Text())
		yaml.Unmarshal(scanner.Bytes(), &confObject)
	}
	return confObject
}

func lcdHub(pushBullet, bmp, sysinfo chan string, lcdChan chan *lcd.LcdEvent, control chan os.Signal) {
	for {
		select {
		case <-control:
			lcdChan <- lcd.NewShutdownEvent()
			close(lcdChan)
			return
		case pushBulletMessage := <-pushBullet:
			lcdChan <- lcd.NewDisplayEvent(pushBulletMessage, 5*time.Second, lcd.BEFORE, 1, true)
		case bmpMessage := <-bmp:
			lcdChan <- lcd.NewDisplayEvent(bmpMessage, 3*time.Second, lcd.NO_FLASH, 1, false)
		case sysinfoMessage := <-sysinfo:
			lcdChan <- lcd.NewDisplayEvent(sysinfoMessage, 3*time.Second, lcd.NO_FLASH, 1, false)

		}
	}
}

func main() {
	config := parseYAMLConf("conf.yml")
	log.Printf("Config:%v", config)

	var pinout []int
	for _, val := range config["pinout"].([]interface{}) {
		pinout = append(pinout, val.(int))
	}
	i2cBus := uint8(config["bmpI2c"].(int))

	display := lcd.NewDisplay(pinout, false)
	defer display.Close()

	initEvent := lcd.NewLcdEvent(lcd.EVENT_DISPLAY, "Display initialized", 3*time.Second, lcd.BEFORE_AND_AFTER, 1, true)
	display.Input <- initEvent

	client := pushbullet.NewClient(config["apiToken"].(string))
	client.StartMonitoring()

	bmpChan := bmp.InitBMP085(i2cBus)

	sysInfoChan := systeminfo.InitSystemMonitoring()

	controlChan := make(chan os.Signal, 1)
	signal.Notify(controlChan, os.Interrupt, os.Kill)

	go lcdHub(client.Output, bmpChan, sysInfoChan, display.Input, controlChan)

	for {
	}
}
