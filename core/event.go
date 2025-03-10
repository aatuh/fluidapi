package core

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// EventType represents the type of event.
type EventType string

// Event represents an emitted event.
type Event struct {
	Type    EventType
	Message string
	Data    any
}

// EventCallback is a function that handles an event.
type EventCallback func(event *Event)

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
	mu        sync.RWMutex   // Mutex for thread safety when emitting events.
	counter   int            // Used to generate unique IDs for listeners.
	timeout   *time.Duration // Optional timeout for each callback.
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

// WithTimeout returns a function that sets the timeout for each callback.
//
// Parameters:
//   - timeout: The timeout duration.
//
// Returns:
//   - func(*EventEmitter): A function that sets the timeout.
func WithTimeout(timeout time.Duration) func(*EventEmitter) {
	return func(e *EventEmitter) {
		e.timeout = &timeout
	}
}

// RegisterListener registers a listener for a specific event type.
//
// Parameters:
//   - eventType: The type of the event.
//   - callback: The function to call when the event is emitted.
//
// Returns:
//   - *EventEmitter: The EventEmitter.
func (e *EventEmitter) RegisterListener(
	eventType EventType, callback EventCallback,
) *EventEmitter {
	// Generate a unique ID for the listener.
	e.mu.Lock()
	defer e.mu.Unlock()
	e.counter++
	id := fmt.Sprintf("%s-%d", eventType, e.counter)

	// Add the listener to the list.
	e.listeners[eventType] = append(e.listeners[eventType], eventListener{
		id:       id,
		callback: callback,
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

// Emit emits an event to all registered listeners. It runs each callback in a
// separate goroutine. It will use the timeout for each callback if it is set.
// If the timeout is not set, the callbacks will be run immediately.
//
// Parameters:
//   - event: The event to emit.
//
// Returns:
//   - *EventEmitter: The EventEmitter.
func (e *EventEmitter) Emit(event *Event) {
	e.mu.RLock()
	listeners := e.listeners[event.Type]
	e.mu.RUnlock()
	// Determine the timeout for each callback.
	var timeout *time.Duration
	if e.timeout != nil {
		timeout = new(time.Duration)
		*timeout = *e.timeout
	}
	// Run each callback in a separate goroutine.
	for _, l := range listeners {
		go func(cb EventCallback, timeout *time.Duration) {
			runCallback(event, cb, timeout)
		}(l.callback, timeout)
	}
}

// runCallback runs a callback with an optional timeout.
func runCallback(event *Event, cb EventCallback, timeout *time.Duration) {
	if timeout == nil {
		// Run the callback immediately.
		cb(event)
	} else {
		// Run the callback with a timeout.
		_, cancel := context.WithTimeout(context.Background(), *timeout)
		defer cancel()
		cb(event)
	}
}
