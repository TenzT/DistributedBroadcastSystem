package native

import (
	"BroadcastService/eventbus"
	"log"
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {
	eventBus := New()
	go func() {
		event := eventbus.Event{
			EventType: eventbus.EVENT_TYPE_NEW_DATA,
			Payload:   "hello",
		}

		time.Sleep(2 * time.Second)
		eventBus.Publish(event)
	}()

	event := <-eventBus.Subscribe()
	log.Println(event)
	log.Println("Test succss")
}
