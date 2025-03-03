package core

import "net/http"

// Middleware represents a function that wraps an http.Handler with additional behavior.
// A Middleware typically performs actions before and/or after calling the next handler.
type Middleware func(http.Handler) http.Handler

// ApplyMiddlewares applies a sequence of middlewares to an http.Handler.
// The middlewares are applied in the order given, such that the first middleware in the list
// will be the outermost wrapper and the last middleware will wrap the target handler directly.
//
// Example: ApplyMiddlewares(finalHandler, m1, m2) yields m1(m2(finalHandler)).
func ApplyMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	wrapped := h
	// Apply each middleware in reverse order so that the first in the list is outermost
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}
