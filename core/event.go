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
//
// Parameters:
//   - eventType: The type of the event.
//   - message: The message of the event.
//
// Returns:
//   - *Event: A new event.
func NewEvent(eventType EventType, message string) *Event {
	return &Event{
		Type:    eventType,
		Message: message,
	}
}

// WithData returns a new event with the given data.
//
// Parameters:
//   - data: The data to include in the event.
//
// Returns:
//   - *Event: A new event.
func (e *Event) WithData(data any) *Event {
	return &Event{
		Type:    e.Type,
		Message: e.Message,
		Data:    data,
	}
}

// EventEmitter is responsible for emitting events.
type EventEmitter struct {
	listeners map[EventType][]func(*Event)
	mu        sync.RWMutex
}

// NewEventEmitter creates a new EventEmitter.
//
// Returns:
//   - *EventEmitter: A new EventEmitter.
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[EventType][]func(*Event)),
	}
}

// RegisterListener registers a listener for a specific event type.
//
// Parameters:
//   - eventType: The type of the event.
//   - listener: The listener function.
//
// Returns:
//   - *EventEmitter: The EventEmitter.
func (e *EventEmitter) RegisterListener(
	eventType EventType, listener func(*Event),
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
//
// Parameters:
//   - event: The event to emit.
//
// Returns:
//   - *EventEmitter: The EventEmitter.
func (e *EventEmitter) Emit(event *Event) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if listeners, found := e.listeners[event.Type]; found {
		for _, listener := range listeners {
			listener(event)
		}
	}
}
