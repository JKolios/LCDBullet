package producers

import (
	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/producers/bitcoinaverage"
	"github.com/JKolios/goLcdEvents/producers/bmp"
	"github.com/JKolios/goLcdEvents/producers/pushbullet"
	"github.com/JKolios/goLcdEvents/producers/systeminfo"
	"github.com/JKolios/goLcdEvents/producers/wunderground"
)

type Producer interface {
	Initialize(config conf.Configuration)
	Start(<-chan struct{}, chan<- events.Event)
}

var producerMap = map[string]func() Producer{
	"pushbullet":     func() Producer { return &pushbullet.PushbulletProducer{} },
	"bmp":            func() Producer { return &bmp.BMPProducer{} },
	"systeminfo":     func() Producer { return &systeminfo.SystemInfoProducer{} },
	"wunderground":   func() Producer { return &wunderground.WundergroundProducer{} },
	"bitcoinaverage": func() Producer { return &bitcoinaverage.BitcoinAverageProducer{} }}

func NewProducer(producerName string) Producer {
	return producerMap[producerName]()
}
