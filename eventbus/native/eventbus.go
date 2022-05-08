package native

import "BroadcastService/eventbus"

// NativeChannelEventBus implements EventBus with native channel
type NativeChannelEventBus struct {
	messageQueue chan eventbus.Event
}

func (eb *NativeChannelEventBus) Publish(event eventbus.Event) {
	eb.messageQueue <- event
}

func (eb *NativeChannelEventBus) Subscribe() <-chan eventbus.Event {
	return eb.messageQueue
}

func NewNativeChannelEventBus() *NativeChannelEventBus {
	return &NativeChannelEventBus{
		messageQueue: make(chan eventbus.Event, 10000),
	}
}
