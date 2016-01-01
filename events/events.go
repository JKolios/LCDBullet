package events

import (
	"time"

	"github.com/JKolios/goLcdEvents/conf"
)

type Producer interface {
	Initialize(config conf.Configuration)
	Subscribe(chan Event)
}

type Consumer interface {
	Initialize(config conf.Configuration)
	Register(chan Event)
}

type Event struct {
	Payload   interface{}
	Type      string
	From      Producer
	CreatedOn time.Time
}
