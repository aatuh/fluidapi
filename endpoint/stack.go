package endpoint

import "github.com/pakkasys/fluidapi/core"

// Wrapper encapsulates a middleware with an identifier and optional metadata.
// ID can be used to identify the middleware type (for reordering or documentation).
// Data can carry additional information (e.g., input schema or config for the middleware).
type Wrapper struct {
	Middleware core.Middleware
	ID         string
	Data       any
}

// Stack is a list of Wrapper, representing an ordered middleware chain.
type Stack []Wrapper

// Middlewares extracts the core.Middleware functions from the stack in order.
func (s Stack) Middlewares() []core.Middleware {
	mws := make([]core.Middleware, len(s))
	for i, wrapper := range s {
		mws[i] = wrapper.Middleware
	}
	return mws
}

// InsertAfterID inserts a middleware Wrapper immediately after the middleware with the given ID in the stack.
// Returns true if inserted, or false if no middleware with that ID was found.
func (s *Stack) InsertAfterID(id string, wrapper Wrapper) bool {
	for i, mw := range *s {
		if mw.ID == id {
			// Insert after position i
			if i == len(*s)-1 {
				*s = append(*s, wrapper)
			} else {
				// Slice insert operation
				*s = append((*s)[:i+1], append([]Wrapper{wrapper}, (*s)[i+1:]...)...)
			}
			return true
		}
	}
	return false
}

// InsertBeforeID inserts a middleware Wrapper immediately before the middleware with the given ID.
// Returns true if inserted, false if the ID was not found.
func (s *Stack) InsertBeforeID(id string, wrapper Wrapper) bool {
	for i, mw := range *s {
		if mw.ID == id {
			*s = append((*s)[:i], append([]Wrapper{wrapper}, (*s)[i:]...)...)
			return true
		}
	}
	return false
}
