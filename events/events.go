package events

import (
	"time"

	"github.com/JKolios/goLcdEvents/conf"
)

const (
	PRIORITY_IMMEDIATE = 2
	PRIORITY_HIGH = 1
	PRIORITY_LOW = 0
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
	Priority int
}
