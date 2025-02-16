package core

import (
	"sync"
)

// EventType represents the type of event.
type EventType string

// Event represents an emitted event.
type Event struct {
	Type    EventType
	Message string
	Data    any
}

// NewEvent creates a new event.
func NewEvent(eventType EventType, message string) *Event {
	return &Event{
		Type:    eventType,
		Message: message,
	}
}

func (e *Event) WithData(data any) *Event {
	e.Data = data
	return e
}

// EventEmitter is responsible for emitting events.
type EventEmitter struct {
	listeners map[EventType][]func(*Event)
	mu        sync.RWMutex
}

// NewEventEmitter creates a new EventEmitter.
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[EventType][]func(*Event)),
	}
}

// RegisterListener registers a listener for a specific event type.
func (e *EventEmitter) RegisterListener(
	eventType EventType,
	listener func(*Event),
) *EventEmitter {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.listeners[eventType]
	if !ok {
		e.listeners[eventType] = []func(*Event){}
	}
	e.listeners[eventType] = append(e.listeners[eventType], listener)
	return e
}

// Emit emits an event to all registered listeners.
func (e *EventEmitter) Emit(event *Event) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if listeners, found := e.listeners[event.Type]; found {
		for _, listener := range listeners {
			listener(event)
		}
	}
}
