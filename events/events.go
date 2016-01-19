package events

import (
	"time"

	"github.com/JKolios/goLcdEvents/conf"
)

const (
	PRIORITY_LOW = iota
	PRIORITY_HIGH
	PRIORITY_IMMEDIATE
)

type Producer interface {
	Initialize(config conf.Configuration)
	Start(<-chan struct{}, chan<- Event)
}

type Consumer interface {
	Initialize(config conf.Configuration)
	Start(<-chan struct{}, <-chan Event)
}

type Event struct {
	Payload   interface{}
	Type      string
	From      Producer
	CreatedOn time.Time
	Priority  int
}
