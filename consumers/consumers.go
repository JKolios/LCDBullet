package consumers

import (
	"github.com/JKolios/goLcdEvents/conf"
	"github.com/JKolios/goLcdEvents/consumers/httplog"
	"github.com/JKolios/goLcdEvents/consumers/lcd"
	"github.com/JKolios/goLcdEvents/consumers/wsclient"
	"github.com/JKolios/goLcdEvents/events"
)

type Consumer interface {
	Initialize(config conf.Configuration)
	Start(<-chan struct{}) chan events.Event
}

var consumerMap = map[string]func() Consumer{
	"lcd":      func() Consumer { return &lcd.LCDConsumer{} },
	"httplog":  func() Consumer { return &httplog.HttpConsumer{} },
	"wsclient": func() Consumer { return &wsclient.WebsocketConsumer{} }}

func NewConsumer(consumerName string) Consumer {
	return consumerMap[consumerName]()
}
