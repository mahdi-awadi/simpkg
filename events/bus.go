package events

import (
	"sync"
)

// IEventbus interface
type IEventbus interface {
	Subscribe(topic string, handler func(v ...any))
	SubscribersCount(topic string) int
	Unsubscribe(topic string)
	Dispatch(topic string, v ...any)
	DispatchAsync(topic string, v ...any)
	DispatchFunc(topic string, fn func(v ...any), v ...any)
	DispatchFuncAsync(topic string, fn func(v ...any), v ...any)
}

// EventBus struct
type EventBus struct {
	handlers map[string][]func(...any)
	locker   sync.Mutex
}

// New returns a new EventBus object.
func New() *EventBus {
	return &EventBus{
		handlers: make(map[string][]func(...any)),
		locker:   sync.Mutex{},
	}
}

// Subscribe to a topic
func (bus *EventBus) Subscribe(topic string, handler func(v ...any)) {
	bus.locker.Lock()
	bus.handlers[topic] = append(bus.handlers[topic], handler)
	bus.locker.Unlock()
}

// SubscribersCount returns subscribers count
func (bus *EventBus) SubscribersCount(topic string) int {
	bus.locker.Lock()
	count := len(bus.handlers[topic])
	bus.locker.Unlock()

	return count
}

// Unsubscribe from a topic
func (bus *EventBus) Unsubscribe(topic string) {
	bus.locker.Lock()
	delete(bus.handlers, topic)
	bus.locker.Unlock()
}

// Dispatch a topic
func (bus *EventBus) Dispatch(topic string, v ...any) {
	bus.dispatch(topic, false, v...)
}

// DispatchAsync dispatch a topic asynchronously
func (bus *EventBus) DispatchAsync(topic string, v ...any) {
	bus.dispatch(topic, true, v...)
}

// DispatchFunc call a function inside subscribe handler
func (bus *EventBus) DispatchFunc(topic string, fn func(v ...any), v ...any) {
	bus.dispatchFunc(topic, false, fn, v...)
}

// DispatchFuncAsync call a function inside subscribe handler asynchronously
func (bus *EventBus) DispatchFuncAsync(topic string, fn func(v ...any), v ...any) {
	bus.dispatchFunc(topic, true, fn, v...)
}

// dispatch a topic
func (bus *EventBus) dispatch(topic string, async bool, v ...any) {
	bus.locker.Lock()
	handlers, ok := bus.handlers[topic]
	bus.locker.Unlock()
	if ok {
		for _, handler := range handlers {
			if async {
				go handler(v...)
			} else {
				handler(v...)
			}
		}
	}
}

// dispatchFunc dispatch a topic
func (bus *EventBus) dispatchFunc(topic string, async bool, fn func(v ...any), v ...any) {
	bus.locker.Lock()
	handlers, ok := bus.handlers[topic]
	bus.locker.Unlock()
	if ok {
		for _, h := range handlers {
			if async {
				go func(handler func(...any), callback func(...any), data ...any) {
					callback(data...)
					handler(data...)
				}(h, fn, v...)
			} else {
				fn(v...)
				h(v...)
			}
		}
	}
}
