package events

import (
	"sync"
)

// EventType represents the type of event.
type EventType string

// Event represents an emitted event.
type Event struct {
	Type    EventType
	Message string
}

func NewEvent(eventType EventType, message string) Event {
	return Event{
		Type:    eventType,
		Message: message,
	}
}

// EventEmitter defines the interface for an event emitter.
type EventEmitter interface {
	RegisterListener(eventType EventType, listener func(Event))
	Emit(event Event)
}

// DefaultEventEmitter is responsible for emitting events.
type DefaultEventEmitter struct {
	listeners map[EventType][]func(Event)
	mu        sync.RWMutex
}

// NewDefaultEventEmitter creates a new DefaultEventEmitter.
func NewDefaultEventEmitter() *DefaultEventEmitter {
	return &DefaultEventEmitter{
		listeners: make(map[EventType][]func(Event)),
	}
}

// RegisterListener registers a listener for a specific event type.
func (emitter *DefaultEventEmitter) RegisterListener(
	eventType EventType,
	listener func(Event),
) {
	emitter.mu.Lock()
	defer emitter.mu.Unlock()
	_, ok := emitter.listeners[eventType]
	if !ok {
		emitter.listeners[eventType] = []func(Event){}
	}
	emitter.listeners[eventType] = append(emitter.listeners[eventType], listener)
}

// Emit emits an event to all registered listeners.
func (emitter *DefaultEventEmitter) Emit(event Event) {
	emitter.mu.RLock()
	defer emitter.mu.RUnlock()
	if listeners, found := emitter.listeners[event.Type]; found {
		for _, listener := range listeners {
			listener(event)
		}
	}
}
