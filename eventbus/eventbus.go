package eventbus

const (
	EVENT_TYPE_UNKNOWN = iota
	EVENT_TYPE_NEW_DATA
)

// EventBus transmiss events
type EventBus interface {
	// Publish sends event to event bus
	Publish(event Event)

	// Subscribe rechives a read-only channel of events
	Subscribe() <-chan Event
}

type Event struct {
	EventType int
	Payload   interface{}
}
