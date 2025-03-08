package core

import (
	"fmt"
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

// WithData returns a new event with the given data. It returns a new event.
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

// eventListener wraps a listener callback with an ID.
type eventListener struct {
	id       string
	callback func(*Event)
}

// EventEmitter is responsible for emitting events.
type EventEmitter struct {
	listeners map[EventType][]eventListener
	mu        sync.RWMutex // Mutex for thread safety when emitting events.
	counter   int          // Used to generate unique IDs for listeners.
}

// NewEventEmitter creates a new EventEmitter.
//
// Returns:
//   - *EventEmitter: A new EventEmitter.
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[EventType][]eventListener),
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
	// Generate a unique ID for the listener.
	e.mu.Lock()
	defer e.mu.Unlock()
	e.counter++
	id := fmt.Sprintf("%s-%d", eventType, e.counter)

	// Add the listener to the list.
	e.listeners[eventType] = append(e.listeners[eventType], eventListener{
		id:       id,
		callback: listener,
	})
	return e
}

// RemoveListener removes a listener for a specific event type.
//
// Parameters:
//   - eventType: The type of the event.
//   - listener: The listener function.
//
// Returns:
//   - *EventEmitter: The EventEmitter.
func (e *EventEmitter) RemoveListener(eventType EventType, id string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if list, found := e.listeners[eventType]; found {
		for i, l := range list {
			if l.id == id {
				// Remove the listener with the matching ID.
				e.listeners[eventType] = append(list[:i], list[i+1:]...)
				break
			}
		}
	}
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
		for _, l := range listeners {
			l.callback(event)
		}
	}
}
