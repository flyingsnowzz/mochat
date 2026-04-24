package event

import (
	"reflect"
	"sync"

	"mochat-api-server/internal/pkg/logger"
)

type Event interface {
	Name() string
}

type Listener func(event Event)

type EventBus struct {
	listeners map[string][]Listener
	mu        sync.RWMutex
}

var Bus = &EventBus{
	listeners: make(map[string][]Listener),
}

func (b *EventBus) Subscribe(eventName string, listener Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners[eventName] = append(b.listeners[eventName], listener)
	logger.Sugar.Debugf("registered listener for event: %s (%s)", eventName, reflect.TypeOf(listener).String())
}

func (b *EventBus) Publish(event Event) {
	b.mu.RLock()
	listeners, ok := b.listeners[event.Name()]
	b.mu.RUnlock()

	if !ok {
		return
	}

	logger.Sugar.Debugf("publishing event: %s", event.Name())
	for _, listener := range listeners {
		go func(l Listener) {
			defer func() {
				if r := recover(); r != nil {
					logger.Sugar.Errorf("event listener panic: %v", r)
				}
			}()
			l(event)
		}(listener)
	}
}

func (b *EventBus) PublishSync(event Event) {
	b.mu.RLock()
	listeners, ok := b.listeners[event.Name()]
	b.mu.RUnlock()

	if !ok {
		return
	}

	for _, listener := range listeners {
		listener(event)
	}
}

func Subscribe(eventName string, listener Listener) {
	Bus.Subscribe(eventName, listener)
}

func Publish(event Event) {
	Bus.Publish(event)
}

func PublishSync(event Event) {
	Bus.PublishSync(event)
}
