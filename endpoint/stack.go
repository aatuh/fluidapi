package endpoint

import (
	"net/http"
	"sync"

	"github.com/pakkasys/fluidapi/core"
)

// Wrapper encapsulates a middleware with an identifier and optional metadata.
// ID can be used to identify the middleware type (e.g. for reordering or
// documentation). Data can carry any type of additional information.
type Wrapper struct {
	Middleware core.Middleware
	ID         string
	Data       any
}

// Stack manages an ordered list of middleware wrappers with concurrency safety.
type Stack struct {
	mu       sync.RWMutex
	wrappers []Wrapper
}

// NewStack creates and returns an initialized middleware Stack.
func NewStack(wrappers ...Wrapper) *Stack {
	return &Stack{
		wrappers: wrappers,
	}
}

// Clone creates a deep copy of the Stack.
//
// Returns:
//   - *Stack: The cloned middleware stack.
func (s *Stack) Clone() *Stack {
	s.mu.RLock()
	defer s.mu.RUnlock()
	newStack := NewStack()
	newStack.wrappers = make([]Wrapper, len(s.wrappers))
	copy(newStack.wrappers, s.wrappers)
	return newStack
}

// Add appends a new middleware Wrapper to the stack and returns the stack for
// chaining.
//
// Parameters:
//   - w: The wrapper to add.
//
// Returns:
//   - *Stack: The updated middleware stack.
func (s *Stack) Add(w Wrapper) *Stack {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wrappers = append(s.wrappers, w)
	return s
}

// InsertBefore inserts a middleware Wrapper before the one with the specified
// ID. Returns true if a matching wrapper was found and insertion happened
// before it; if no match is found, the new wrapper is appended and false is
// returned.
//
// Parameters:
//   - id: The ID of the wrapper to insert before.
//   - w: The wrapper to insert.
//
// Returns:
//   - *Stack: The updated middleware stack.
//   - bool: True if a matching wrapper was found and insertion succeeded.
func (s *Stack) InsertBefore(id string, w Wrapper) (*Stack, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, wrapper := range s.wrappers {
		if wrapper.ID == id {
			s.wrappers = append(s.wrappers[:i],
				append([]Wrapper{w}, s.wrappers[i:]...)...)
			return s, true
		}
	}
	s.wrappers = append(s.wrappers, w)
	return s, false
}

// InsertAfter inserts a middleware Wrapper after the one with the specified ID.
// Returns true if a matching wrapper was found and insertion happened after it.
// If no match is found, the new wrapper is appended and false is returned.
//
// Parameters:
//   - id: The ID of the wrapper to insert after.
//   - w: The wrapper to insert.
//
// Returns:
//   - *Stack: The updated middleware stack.
//   - bool: True if a matching wrapper was found and insertion succeeded.
func (s *Stack) InsertAfter(id string, w Wrapper) (*Stack, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, wrapper := range s.wrappers {
		if wrapper.ID == id {
			pos := i + 1
			s.wrappers = append(s.wrappers[:pos],
				append([]Wrapper{w}, s.wrappers[pos:]...)...)
			return s, true
		}
	}
	s.wrappers = append(s.wrappers, w)
	return s, false
}

// Remove deletes the middleware Wrapper with the specified ID from the stack.
// Returns true if the middleware was found and removed; false otherwise.
//
// Parameters:
//   - id: The ID of the wrapper to remove.
//
// Returns:
//   - *Stack: The updated middleware stack.
//   - bool: True if the middleware was found and removed; false otherwise.
func (s *Stack) Remove(id string) (*Stack, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, wrapper := range s.wrappers {
		if wrapper.ID == id {
			s.wrappers = append(s.wrappers[:i], s.wrappers[i+1:]...)
			return s, true
		}
	}
	return s, false
}

// Apply applies the middleware wrappers to the given http.Handler.
// Wrappers are applied in reverse order so that the first wrapper is the
// outermost.
//
// Parameters:
//   - h: The http.Handler to wrap.
//
// Returns:
//   - http.Handler: The wrapped http.Handler.
func (s *Stack) Apply(h http.Handler) http.Handler {
	s.mu.RLock()
	defer s.mu.RUnlock()
	middlewares := make([]core.Middleware, len(s.wrappers))
	for i, wrapper := range s.wrappers {
		middlewares[i] = wrapper.Middleware
	}
	core.ApplyMiddlewares(h, middlewares...)
	return h
}
