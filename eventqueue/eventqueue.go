package eventqueue

import (
	"container/list"
	"github.com/JKolios/goLcdEvents/conf"

	"github.com/JKolios/goLcdEvents/consumers"
	"github.com/JKolios/goLcdEvents/events"
	"github.com/JKolios/goLcdEvents/inspector"
	"github.com/JKolios/goLcdEvents/producers"
	"log"
	"time"
)

type EventQueue struct {
	highPriorityeventList, lowPriorityeventList list.List
	producers                                   []producers.Producer
	producerChan                                chan events.Event
	consumerChannels                            []chan events.Event
	consumers                                   []consumers.Consumer
	done                                        chan struct{}
}

func (queue *EventQueue) producerHub() {

	var incomingEvent events.Event

	for {
		select {
		case <-queue.done:
			log.Println("producerHub halting")
			return

		case incomingEvent = <-queue.producerChan:
			log.Printf("producerHub got event: %+v\n", incomingEvent)
			// Inspect the event and handle according to type and priority
			if incomingEvent.Priority == events.PRIORITY_HIGH {
				queue.highPriorityeventList.PushBack(incomingEvent)
			} else {
				queue.lowPriorityeventList.PushBack(incomingEvent)
			}

		}
	}
}

func (queue *EventQueue) consumerHub() {

	var selectedEvent events.Event

	for {
		select {

		case <-queue.done:
			log.Println("consumerHub halting")
			return

		default:

			if queue.highPriorityeventList.Len() > 0 {
				selectedEvent = queue.highPriorityeventList.Remove(queue.highPriorityeventList.Front()).(events.Event)
			} else if queue.lowPriorityeventList.Len() > 0 {
				selectedEvent = queue.lowPriorityeventList.Remove(queue.lowPriorityeventList.Front()).(events.Event)
			} else {
				continue
			}

			log.Printf("consumerHub selected event: %+v \n", selectedEvent)

			for _, consumerChan := range queue.consumerChannels {
				consumerChan <- selectedEvent
			}

		}
	}
}

func NewQueue(consumerNames, producerNames []string, config conf.Configuration) *EventQueue {

	queue := &EventQueue{}

	queue.done = make(chan struct{})

	// Consumer Init
	for _, consumerName := range consumerNames {
		consumer := consumers.NewConsumer(consumerName)
		consumer.Initialize(config)
		queue.consumers = append(queue.consumers, consumer)

	}

	// Producer Init
	for _, producerName := range producerNames {
		producer := producers.NewProducer(producerName)
		producer.Initialize(config)
		queue.producers = append(queue.producers, producer)

	}

	return queue
}

func (queue *EventQueue) Start() {

	queue.producerChan = make(chan events.Event)

	for _, producer := range queue.producers {
		producer.Start(queue.done, queue.producerChan)
	}

	for _, consumer := range queue.consumers {
		queue.consumerChannels = append(queue.consumerChannels, consumer.Start(queue.done))
	}

	go queue.producerHub()
	go queue.consumerHub()

	go inspector.ListReport(&queue.lowPriorityeventList, "Low Priority", time.Minute*10)
	go inspector.ListReport(&queue.highPriorityeventList, "High Priority", time.Minute*10)

	go inspector.ListCleaner(&queue.lowPriorityeventList, "Low Priority", time.Minute*10)
	go inspector.ListCleaner(&queue.highPriorityeventList, "High Priority", time.Minute*10)
}

func (queue *EventQueue) Stop() {
	close(queue.done)
}
